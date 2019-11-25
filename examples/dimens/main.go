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
		Width:  1600,
		Height: 900,
		Title:  "Hello",
	}, start, stop, update)
}

var (
	dwarf *Model

	fragShader string
	vertShader string

	attribVCoord  gl.Attrib
	attribVNormal gl.Attrib

	meshVertices gl.Buffer
	meshNormals  gl.Buffer
	meshElements gl.Buffer

	program gl.Program
)

func start(ctx *sparkle.Context) {
	dwarf, _ = LoadObj(
		"dwarf.obj", "dwarf_diffuse.png", 1600, 900)

	fragShader, _ = sparkle.LoadShader("two_sided_shading.frag")
	vertShader, _ = sparkle.LoadShader("two_sided_shading.frag")

	program, _ = glutil.CreateProgram(ctx.GL, vertShader, fragShader)

	attribVCoord = ctx.GL.GetAttribLocation(program, "v_coord")
	attribVNormal = ctx.GL.GetAttribLocation(program, "v_normal")

	meshVertices = ctx.GL.CreateBuffer()
	vertDat := RobustBytes(binary.LittleEndian, dwarf.outVertices...)
	ctx.GL.BufferData(gl.ARRAY_BUFFER, vertDat, gl.STATIC_DRAW)

	meshNormals = ctx.GL.CreateBuffer()
	normDat := RobustBytes(binary.LittleEndian, dwarf.outNormals...)
	ctx.GL.BufferData(gl.ARRAY_BUFFER, normDat, gl.STATIC_DRAW)

	meshElements = ctx.GL.CreateBuffer()
	// TODO create elements based on OBJ faces.
}

func stop(ctx *sparkle.Context) {

}

func update(ctx *sparkle.Context) {
	g := ctx.GL

	g.EnableVertexAttribArray(attribVCoord)
	g.BindBuffer(gl.ARRAY_BUFFER, meshVertices)
	g.VertexAttribPointer(
		attribVCoord,
		4,
		gl.FLOAT,
		false,
		0,
		0,
	)

	g.EnableVertexAttribArray(attribVNormal)
	g.BindBuffer(gl.ARRAY_BUFFER, meshNormals)
	g.VertexAttribPointer(
		attribVNormal,
		3,
		gl.FLOAT,
		false,
		0,
		0,
	)

	g.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, meshElements)
	size := g.GetBufferParameteri(gl.ELEMENT_ARRAY_BUFFER, gl.BUFFER_SIZE)
	g.DrawElements(gl.TRIANGLES, size, gl.UNSIGNED_INT, 0)
}
