package main

import (
	"encoding/binary"

	"github.com/damienfamed75/pine/view"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hecate-tech/sparkle"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

func main() {
	sparkle.Main(&screen.NewWindowOptions{
		Width:  1600,
		Height: 900,
		Title:  "Hello, World!",
	}, start, stop, update)
}

var (
	model *view.Model
	cam   *view.Camera

	vertShader string
	fragShader string
	buffer     sparkle.Buffer
	program    sparkle.Program

	texture      gl.Texture
	position     sparkle.Attrib
	color        sparkle.Uniform
	offset       sparkle.Uniform
	green        float32 // Color of the triangle.
	triangleData = RobustBytes(binary.LittleEndian,
		mgl64.Vec3{0.0, 0.4, 0.0},
		mgl64.Vec3{0.0, 0.0, 0.0},
		mgl64.Vec3{0.4, 0.0, 0.0},
	)

	cubeData = f32.Bytes(binary.LittleEndian,
		// Front
		-1.0, -1.0, -1.0, // triangle 1 : begin
		-1.0, -1.0, 1.0,
		-1.0, 1.0, 1.0, // triangle 1 : end
		1.0, 1.0, -1.0, // triangle 2 : begin
		-1.0, -1.0, -1.0,
		-1.0, 1.0, -1.0, // triangle 2 : end
		// Right
		1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,
		-1.0, -1.0, -1.0,
		// Back
		-1.0, -1.0, -1.0,
		-1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
		1.0, -1.0, 1.0,
		-1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0,
		// Left
		-1.0, 1.0, 1.0,
		-1.0, -1.0, 1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		// Bottom
		1.0, -1.0, -1.0,
		1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, -1.0,
		-1.0, 1.0, -1.0,
		// Top
		1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
		-1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
	)

	cubeColors = f32.Bytes(binary.LittleEndian,
		// Front
		0.583, 0.771, 0.014,
		0.609, 0.115, 0.436,
		0.327, 0.483, 0.844,
		0.822, 0.569, 0.201,
		0.435, 0.602, 0.223,
		0.310, 0.747, 0.185,
		// Right
		0.597, 0.770, 0.761,
		0.559, 0.436, 0.730,
		0.359, 0.583, 0.152,
		0.483, 0.596, 0.789,
		0.559, 0.861, 0.639,
		0.195, 0.548, 0.859,
		// Back
		0.014, 0.184, 0.576,
		0.771, 0.328, 0.970,
		0.406, 0.615, 0.116,
		0.676, 0.977, 0.133,
		0.971, 0.572, 0.833,
		0.140, 0.616, 0.489,
		// Left
		0.997, 0.513, 0.064,
		0.945, 0.719, 0.592,
		0.543, 0.021, 0.978,
		0.279, 0.317, 0.505,
		0.167, 0.620, 0.077,
		0.347, 0.857, 0.137,
		// Bottom
		0.055, 0.953, 0.042,
		0.714, 0.505, 0.345,
		0.783, 0.290, 0.734,
		0.722, 0.645, 0.174,
		0.302, 0.455, 0.848,
		0.225, 0.587, 0.040,
		// Top
		0.517, 0.713, 0.338,
		0.053, 0.959, 0.120,
		0.393, 0.621, 0.362,
		0.673, 0.211, 0.457,
		0.820, 0.883, 0.371,
		0.982, 0.099, 0.879,
	)
)

func RobustBytes(byteOrder binary.ByteOrder, values ...mgl64.Vec3) []byte {
	newVals := make([]float32, len(values)*3)

	for i := range values {
		newVals[3*i+0] = float32(values[i].X())
		newVals[3*i+1] = float32(values[i].Y())
		newVals[3*i+2] = float32(values[i].Z())
	}

	return f32.Bytes(byteOrder, newVals...)
}

var (
	cubeColorBuf gl.Buffer
)

func start(ctx *sparkle.Context) {
	var aspect float64 = 1600.0 / 900.0
	cam = view.NewCamera(mgl64.Vec3{1, 0.75, 1}, mgl64.DegToRad(90), aspect)
	model, _ = view.LoadObj("dwarf.obj", "dwarf_diffuse.png", 1600, 900, cam)

	ctx.AddDrawer(draw)
	vertShader, _ = sparkle.LoadShader("cube.vert")
	fragShader, _ = sparkle.LoadShader("cube.frag")

	program, _ = glutil.CreateProgram(ctx.GL, vertShader, fragShader)

	buffer = ctx.GL.CreateBuffer()
	ctx.GL.BindBuffer(gl.ARRAY_BUFFER, buffer)
	ctx.GL.BufferData(gl.ARRAY_BUFFER, cubeData, gl.STATIC_DRAW)
	ctx.GL.Enable(gl.DEPTH_TEST)
	// ctx.GL.BufferData(gl.ARRAY_BUFFER, triangleData, gl.STATIC_DRAW)

	// position = ctx.GL.GetAttribLocation(program, "position")
	// color = ctx.GL.GetUniformLocation(program, "color")
	// offset = ctx.GL.GetUniformLocation(program, "offset")
	// texture = ctx.GL.CreateTexture()

	cubeColorBuf = ctx.GL.CreateBuffer()
	ctx.GL.BindBuffer(gl.ARRAY_BUFFER, cubeColorBuf)
	ctx.GL.BufferData(gl.ARRAY_BUFFER, cubeColors, gl.STATIC_DRAW)
}

func draw(ctx *sparkle.Context) {
	// green += 0.01
	// if green > 1 {
	// 	green = 0
	// }

	ctx.GL.UseProgram(program)
	// ctx.GL.Uniform4f(
	// 	color,
	// 	0, green, 0, 1,
	// )

	// ctx.GL.Uniform2f(
	// 	offset,
	// 	ctx.GetTouchX()/float32(ctx.Size.WidthPx),
	// 	ctx.GetTouchY()/float32(ctx.Size.HeightPx),
	// )

	// ctx.GL.BindBuffer(gl.ARRAY_BUFFER, buffer)
	ctx.GL.EnableVertexAttribArray(gl.Attrib{Value: 0})
	{
		ctx.GL.BindBuffer(gl.ARRAY_BUFFER, buffer)
		ctx.GL.VertexAttribPointer(
			gl.Attrib{Value: 0}, 3,
			gl.FLOAT, false, 0, 0,
		)
		// ctx.GL.BindBuffer(gl.ARRAY_BUFFER, cubeColorBuf)
		// ctx.GL.VertexAttribPointer(
		// 	gl.Attrib{Value: 1}, 3,
		// 	gl.FLOAT, false, 0, 0,
		// )
		ctx.GL.DrawArrays(gl.TRIANGLES, 0, 12*3)
		// ctx.GL.DrawArrays(gl.TRIANGLES, 0, 12*3)
	}
	ctx.GL.DisableVertexAttribArray(gl.Attrib{Value: 0})

	// model.Draw(ctx.Buffer)

	// ctx.GL.BufferData(
	// 	gl.ARRAY_BUFFER,
	// 	RobustBytes(binary.LittleEndian, model.GetOutVertices()...),
	// 	gl.STATIC_DRAW,
	// )
	// ctx.GL.DrawArrays(gl.TRIANGLES, 0, len(model.GetOutVertices()))
}

// var last = time.Now()

func update(ctx *sparkle.Context) {
	// dt := float64(time.Since(last).Milliseconds()) / 1000.0
	// last = time.Now()
	// model.AddRotation(mgl64.RadToDeg(dt*50000), mgl64.Vec3{0, 1, 0})
}

func stop(ctx *sparkle.Context) {

}
