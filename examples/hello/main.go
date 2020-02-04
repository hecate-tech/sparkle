package main

import (
	"encoding/binary"

	"github.com/hecate-tech/sparkle"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

func main() {
	sparkle.Main(
		&screen.NewWindowOptions{},
		start, stop, update,
	)
}

var (
	tData = f32.Bytes(binary.LittleEndian,
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.5, 0.0,
	)
	tBuffer   gl.Buffer
	program   gl.Program
	aPosition gl.Attrib
	aOff      gl.Uniform
	VAO       gl.VertexArray
)

func start(ctx *sparkle.Context) {
	// VAO = ctx.GL.CreateVertexArray()

	sVert, err := sparkle.LoadShader("basic.vert")
	if err != nil {
		panic(err)
	}

	sFrag, err := sparkle.LoadShader("basic.frag")
	if err != nil {
		panic(err)
	}

	program, err = glutil.CreateProgram(ctx.GL, sVert, sFrag)
	if err != nil {
		panic(err)
	}

	tBuffer = ctx.GL.CreateBuffer()
	ctx.GL.BindBuffer(gl.ARRAY_BUFFER, tBuffer)
	ctx.GL.BufferData(gl.ARRAY_BUFFER, tData, gl.STATIC_DRAW)

	// ctx.GL.LinkProgram(program)

	aPosition = ctx.GL.GetAttribLocation(program, "position")
	aOff = ctx.GL.GetUniformLocation(program, "offset")

	// ctx.GL.BindVertexArray(VAO)

	ctx.AddDrawer(draw)
}

func update(ctx *sparkle.Context) {
}

func draw(ctx *sparkle.Context) {
	g := ctx.GL

	g.UseProgram(program)

	g.BindBuffer(gl.ARRAY_BUFFER, tBuffer)

	g.Uniform2f(
		aOff,
		ctx.GetTouchX()/float32(ctx.Size.WidthPx),
		ctx.GetTouchY()/float32(ctx.Size.HeightPx),
	)

	g.EnableVertexAttribArray(aPosition)
	{
		g.VertexAttribPointer(
			aPosition, 3,
			gl.FLOAT, false, 0, 0,
		)
		g.DrawArrays(gl.TRIANGLES, 0, 3)
	}
	g.DisableVertexAttribArray(aPosition)
}

func stop(ctx *sparkle.Context) {

}
