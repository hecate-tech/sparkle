package main

import (
	"encoding/binary"
	"time"

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

	position sparkle.Attrib
	color    sparkle.Uniform
	offset   sparkle.Uniform
	green    float32 // Color of the triangle.
	// triangleData = f32.Bytes(binary.LittleEndian,
	// 	0.0, 0.4, 0.0, // top left
	// 	0.0, 0.0, 0.0, // bottom left
	// 	0.4, 0.0, 0.0, // bottom right
	// )
	triangleData = RobustBytes(binary.LittleEndian,
		mgl64.Vec3{0.0, 0.4, 0.0},
		mgl64.Vec3{0.0, 0.0, 0.0},
		mgl64.Vec3{0.4, 0.0, 0.0},
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

// func Bytes(byteOrder binary.ByteOrder, values ...mgl64.Vec3) []byte {
// 	le := false
// 	switch byteOrder {
// 	case binary.BigEndian:
// 	case binary.LittleEndian:
// 		le = true
// 	default:
// 		panic(fmt.Sprintf("invalid byte order %v", byteOrder))
// 	}

// 	b := make([]byte, 4*(len(values)*3))
// 	for i, v := range values {
// 		for j, vv := range v {
// 			u := math.Float32bits(float32(vv))
// 			// u := math.Float64bits(vv)
// 			if le {
// 				b[4*(i+j+0)] = byte(u >> 0)
// 				b[4*(i+j+1)] = byte(u >> 8)
// 				b[4*(i+j+2)] = byte(u >> 16)
// 				b[4*(i+j+3)] = byte(u >> 24)
// 			} else {
// 				b[4*(i+j+0)] = byte(u >> 24)
// 				b[4*(i+j+1)] = byte(u >> 16)
// 				b[4*(i+j+2)] = byte(u >> 8)
// 				b[4*(i+j+3)] = byte(u >> 0)
// 			}
// 		}
// 	}
// 	return b
// }

func start(ctx *sparkle.Context) {
	var aspect float64 = 1600.0 / 900.0
	cam = view.NewCamera(mgl64.Vec3{1, 0.75, 1}, mgl64.DegToRad(90), aspect)
	model, _ = view.LoadObj("dwarf.obj", "dwarf_diffuse.png", 1600, 900, cam)

	ctx.AddDrawer(draw)
	vertShader, _ = sparkle.LoadShader("basic.vert")
	fragShader, _ = sparkle.LoadShader("basic.frag")

	program, _ = glutil.CreateProgram(ctx.GL, vertShader, fragShader)

	buffer = ctx.GL.CreateBuffer()
	ctx.GL.BindBuffer(gl.ARRAY_BUFFER, buffer)
	ctx.GL.BufferData(gl.ARRAY_BUFFER, triangleData, gl.STATIC_DRAW)

	position = ctx.GL.GetAttribLocation(program, "position")
	color = ctx.GL.GetUniformLocation(program, "color")
	offset = ctx.GL.GetUniformLocation(program, "offset")

}

func draw(ctx *sparkle.Context) {
	green += 0.01
	if green > 1 {
		green = 0
	}

	ctx.GL.UseProgram(program)
	ctx.GL.Uniform4f(
		color,
		0, green, 0, 1,
	)

	ctx.GL.Uniform2f(
		offset,
		ctx.GetTouchX()/float32(ctx.Size.WidthPx),
		ctx.GetTouchY()/float32(ctx.Size.HeightPx),
	)

	ctx.GL.BindBuffer(gl.ARRAY_BUFFER, buffer)

	ctx.GL.EnableVertexAttribArray(position)
	{
		ctx.GL.VertexAttribPointer(
			position, 3,
			gl.FLOAT, false, 0, 0,
		)
		ctx.GL.DrawArrays(gl.TRIANGLES, 0, 3)
	}
	ctx.GL.DisableVertexAttribArray(position)

	model.Draw(ctx.Buffer)

	// ctx.GL.BufferData(
	// 	gl.ARRAY_BUFFER,
	// 	RobustBytes(binary.LittleEndian, model.GetOutVertices()...),
	// 	gl.STATIC_DRAW,
	// )
	// ctx.GL.DrawArrays(gl.TRIANGLES, 0, len(model.GetOutVertices()))
}

var last = time.Now()

func update(ctx *sparkle.Context) {
	// dt := float64(time.Since(last).Milliseconds()) / 1000.0
	// last = time.Now()
	// model.AddRotation(mgl64.RadToDeg(dt*50000), mgl64.Vec3{0, 1, 0})
}

func stop(ctx *sparkle.Context) {

}
