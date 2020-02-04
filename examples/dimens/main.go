package main

import (
	"encoding/binary"

	"github.com/hecate-tech/sparkle"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

func main() {
	sparkle.Main(&screen.NewWindowOptions{
		Width:  800,
		Height: 600,
		Title:  "Hello",
	}, start, stop, update)
}

var (
	dwarf *Model

	// The shaders as strings.
	fragShader string
	vertShader string

	// The attributes we want to have or alter in the shaders.
	attribVCoord  gl.Attrib
	attribVNormal gl.Attrib

	// Buffers of the vertex information.
	// Vertex information is slices of points and then using little endian
	// pattern, they are turned into byte slices and the pointer value is stored
	// within the gl.Buffer.
	meshVertices gl.Buffer
	meshNormals  gl.Buffer
	meshElements gl.Buffer

	position gl.Attrib
	offset   gl.Uniform

	elemCount int

	// The shader program including the vertex and fragment shader.
	program gl.Program

	dwarfTriangleCount int

	VAO gl.VertexArray
)

func start(ctx *sparkle.Context) {
	dwarf, err := LoadObj(
		"dwarf.obj", "dwarf_diffuse.png", 1600, 900, ctx)
	if err != nil {
		panic(err)
	}

	VAO = ctx.GL.CreateVertexArray()

	vertShader, _ = sparkle.LoadShader("basic.vert")
	fragShader, _ = sparkle.LoadShader("basic.frag")

	program, _ = glutil.CreateProgram(ctx.GL, vertShader, fragShader)

	position = ctx.GL.GetAttribLocation(program, "position")
	offset = ctx.GL.GetUniformLocation(program, "offset")
	// attribVCoord = ctx.GL.GetAttribLocation(program, "v_coord")
	// attribVNormal = ctx.GL.GetAttribLocation(program, "v_normal")

	meshVertices = ctx.GL.CreateBuffer()
	vertDat := RobustBytes(binary.LittleEndian, dwarf.outVertices...)
	ctx.GL.BufferData(gl.ARRAY_BUFFER, vertDat, gl.STATIC_DRAW)

	meshNormals = ctx.GL.CreateBuffer()
	ctx.GL.BindBuffer(gl.ARRAY_BUFFER, meshNormals)
	normDat := RobustBytes(binary.LittleEndian, dwarf.outNormals...)
	ctx.GL.BufferData(gl.ARRAY_BUFFER, normDat, gl.STATIC_DRAW)

	dwarfTriangleCount = len(dwarf.outNormals)

	// meshElements = ctx.GL.CreateBuffer()
	// fe := make([]float32, len(dwarf.outElements))

	// elemCount = len(dwarf.outElements)

	// for i := range dwarf.outElements {
	// 	fe[i] = float32(dwarf.outElements[i])
	// }
	// elemDat := f32.Bytes(binary.LittleEndian, fe...)
	// ctx.GL.BufferData(gl.ELEMENT_ARRAY_BUFFER, elemDat, gl.STATIC_DRAW)

	// dwarfTri := RobustBytes(binary.LittleEndian, dwarf.outNormals...)
	// ctx.GL.BufferData(gl.ARRAY_BUFFER, dwarfTri, gl.STATIC_DRAW)

	// Add a drawing function to sparkle.
	ctx.AddDrawer(draw)
}

func stop(ctx *sparkle.Context) {

}

func update(*sparkle.Context) {

}

func draw(ctx *sparkle.Context) {
	g := ctx.GL

	// Use the shader program built by the start.
	// g.BindVertexArray(VAO)
	g.UseProgram(program)

	// Bind mesh elements.
	g.BindBuffer(gl.ARRAY_BUFFER, meshNormals)

	g.Uniform2f(
		offset,
		1,
		0.5,
	)

	g.EnableVertexAttribArray(position)
	{
		g.VertexAttribPointer(
			position, 3,
			gl.FLOAT, true, 0, 0,
		)

		g.DrawArrays(gl.TRIANGLES, 0, dwarfTriangleCount)
		// g.DrawArrays(gl.TRIANGLES, 0, dwarfTriangleCount)

		// g.DrawElements(gl.TRIANGLES, elemCount, gl.UNSIGNED_INT, 0)
		// size := g.GetBufferParameteri(gl.ELEMENT_ARRAY_BUFFER, gl.BUFFER_SIZE)

		// g.DrawElements(gl.ELEMENT_ARRAY_BUFFER, size, gl.TRIANGLES, 0)
	}
	g.DisableVertexAttribArray(position)

	// Enable the position of the model.
	// g.EnableVertexAttribArray(attribVCoord)
	// g.BindBuffer(gl.ARRAY_BUFFER, meshVertices)
	// g.VertexAttribPointer(
	// 	attribVCoord,
	// 	4,
	// 	gl.FLOAT,
	// 	false,
	// 	0,
	// 	0,
	// )

	// // Enables the position of the normalized vertices.
	// g.EnableVertexAttribArray(attribVNormal)
	// g.BindBuffer(gl.ARRAY_BUFFER, meshNormals)
	// g.VertexAttribPointer(
	// 	attribVNormal,
	// 	3,
	// 	gl.FLOAT,
	// 	false,
	// 	0,
	// 	0,
	// )

	// g.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, meshElements)
	// size := g.GetBufferParameteri(gl.ELEMENT_ARRAY_BUFFER, gl.BUFFER_SIZE)

	// g.DrawElements(gl.TRIANGLES, size, gl.UNSIGNED_INT, 0)

}
