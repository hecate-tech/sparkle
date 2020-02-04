package main

import (
	"encoding/binary"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hecate-tech/sparkle"
	"golang.org/x/mobile/gl"
)

type Texture struct {
	gl.Texture
	TexType string
}

type Mesh struct {
	Vertices []mgl64.Vec3
	Indices  []uint
	Textures []Texture

	vao gl.VertexArray
	vbo gl.Buffer
	ebo gl.Buffer
}

func (m *Mesh) setupMesh(ctx *sparkle.Context)  {
	m.vao = ctx.GL.CreateVertexArray()
	m.vbo = ctx.GL.CreateBuffer()
	m.ebo = ctx.GL.CreateBuffer()

	ctx.GL.BindVertexArray(m.vao)
	ctx.GL.BindBuffer(gl.ARRAY_BUFFER, m.vbo)

	ctx.GL.BufferData(gl.ARRAY_BUFFER,
		RobustBytes(binary.LittleEndian, m.Vertices...), gl.STATIC_DRAW)

	ctx.GL.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.ebo)

}

func NewMesh(verts []mgl64.Vec3, ind []uint, texs []Texture, ctx *sparkle.Context) *Mesh {
	m := &Mesh{
		Vertices: verts,
		Indices: ind,
		Textures: texs,
	}
	
	m.setupMesh(ctx)
	
	return m
}
