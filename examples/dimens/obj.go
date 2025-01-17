package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"image"
	"os"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/oakmound/oak/render"
	"golang.org/x/mobile/exp/f32"
)

// Model represents a 3D model in world space.
// This model is a render.Renderable.
type Model struct {
	// a render.Sprite has a position and a buffer of image data which
	// it uses to draw to the screen at that position.
	*render.Sprite
	// the textureData is the local texture file (.bmp in the original, .png in this version)
	// that is referred to to color each triangle face
	textureData *image.RGBA

	outVertices []mgl64.Vec3
	outUVs      []mgl64.Vec3
	outNormals  []mgl64.Vec3

	// quat represents the model's rotation.
	// the quaternion isn't directly applied to the transform, but instead
	// applied to the camera's viewing position and angles to create
	// the illusion of full 3D rotation.
	quat mgl64.Quat

	scale    mgl64.Mat4
	position mgl64.Mat4
	// transform represents the model's position and scale.
	// transform mgl64.Mat4
	angle float64
}

func RobustBytes(byteOrder binary.ByteOrder, values ...mgl64.Vec3) []byte {
	newVals := make([]float32, len(values)*3)

	for i := range values {
		newVals[3*i+0] = float32(values[i].X())
		newVals[3*i+1] = float32(values[i].Y())
		newVals[3*i+2] = float32(values[i].Z())
	}

	return f32.Bytes(byteOrder, newVals...)
}

// SetRotation resets the rotation of the model to what is provided.
// If you wish to rotate based on the current rotation then please refer to
// AddRotation instead.
func (m *Model) SetRotation(angle float64, axis mgl64.Vec3) {
	m.angle = angle
	m.quat = mgl64.QuatRotate(angle, axis)
}

// AddRotation adds to the existing rotation axis.
func (m *Model) AddRotation(angle float64, axis mgl64.Vec3) {
	m.angle += angle
	m.quat = mgl64.QuatRotate(m.angle, axis)
}

// GetTransform combines the scale and position to give the transform matrix of
// this model.
func (m *Model) GetTransform() mgl64.Mat4 {
	return m.scale.Mul4(m.position)
}

// GetScale gets the model's scale on its x, y, and z axis.
func (m *Model) GetScale() mgl64.Vec3 {
	return m.scale.Diag().Vec3()
}

// SetScale sets the model's relative scale on the x, y, and z axis
// If you wish to scale based on thge current rotation then please refer to
// AddScale instead.
func (m *Model) SetScale(x, y, z float64) {
	m.scale = mgl64.Scale3D(x, y, z)
}

// AddScale scales the object relative to its current scale.
func (m *Model) AddScale(x, y, z float64) {
	m.scale = m.scale.Add(mgl64.Scale3D(x, y, z))
}

// SetPosition sets the position of the object from 0,0,0.
// If you wish to set the scale based on its current position then please refer
// to AddPosition instead.
func (m *Model) SetPosition(x, y, z float64) {
	m.position = mgl64.Translate3D(x, y, z)
}

// AddPosition sets the model's position relative to its current position.
func (m *Model) AddPosition(x, y, z float64) {
	m.position = m.position.Add(mgl64.Translate3D(x, y, z))
}

// LoadObj loads a .obj file into memory, loading all its information
// including the texture for the .obj file.
//
// v - vertices
// vn - vertex normalized
// vt - vertex texture coordinate
// f - faces (triangles)
// mtl files are for another time for now.
// f 1/13/4 51/13/5 2/42/26
//				  3rd coord
//        2nd coord
// 1st coord
func LoadObj(objFile, texFile string, w, h int) (*Model, error) {
	fobj, err := os.Open(objFile)
	if err != nil {
		return nil, err
	}
	defer fobj.Close()

	tex, err := render.LoadSprite("model", texFile)
	if err != nil {
		return nil, err
	}

	mod := &Model{
		// Raw texture data from pixel to pixel.
		textureData: tex.GetRGBA(),
		// Empty sprite that has an assigned width and height.
		Sprite: render.NewEmptySprite(0, 0, w, h),
		// quat: mgl64.QuatRotate(mgl64.DegToRad(0), mgl64.Vec3{0, 1, 0}),
		scale:    mgl64.Scale3D(1, 1, 1),
		position: mgl64.Translate3D(0, 0, 0),
	}

	// quat := mgl64.QuatIdent().Rotate(mgl64.Vec3{
	// 	0, mgl64.DegToRad(45), 0,
	// })
	// mod.transform = mgl64.Translate3D(0, 0, 0).
	// 	Mul4(quat.Mat4()).
	// 	Mul4(mgl64.Scale3D(1, 1, 1))
	// mod.transform = mgl64.Translate3D(0, 0, -2).
	// 	Mul4(mgl64.HomogRotate3D(mgl64.DegToRad(45), mgl64.Vec3{0, 1, 0})).
	// 	Mul4(mgl64.Scale3D(1, 1, 1))

	var (
		uvIndices     []uint
		vertexIndices []uint
		normalIndices []uint

		tmpUVs      []mgl64.Vec3
		tmpVertices []mgl64.Vec3
		tmpNormals  []mgl64.Vec3
	)

	scanner := bufio.NewScanner(fobj)

	for scanner.Scan() {
		var (
			vertex struct{ x, y, z float64 }
		)

		line := scanner.Text()

		if len(line) < 2 {
			continue
		}
		if line[0] == 'v' && line[1] == 'n' {
			// vertex normals.
			fmt.Sscanf(line, "vn %f %f %f", &vertex.x, &vertex.y, &vertex.z)
			tmpNormals = append(tmpNormals, mgl64.Vec3{
				vertex.x, vertex.y, vertex.z,
			})
		} else if line[0] == 'v' && line[1] == 't' {
			// vertex texture coordinates.
			// Most of the time an obj file will not have a Z point, but we'll
			// include it anyway in the case that an obj file actually uses it.
			fmt.Sscanf(line, "vt %f %f %f", &vertex.x, &vertex.y, &vertex.z)
			tmpUVs = append(tmpUVs, mgl64.Vec3{
				vertex.x, vertex.y, vertex.z,
			})
		} else if line[0] == 'v' {
			// vertices
			fmt.Sscanf(line, "v %f %f %f", &vertex.x, &vertex.y, &vertex.z)
			tmpVertices = append(tmpVertices, mgl64.Vec3{
				vertex.x, vertex.y, vertex.z,
			})
		} else if line[0] == 'f' {
			var (
				uvIndex     [3]uint
				vertexIndex [3]uint
				normalIndex [3]uint
			)

			fmt.Sscanf(line, "f %d/%d/%d %d/%d/%d %d/%d/%d",
				&vertexIndex[0], &uvIndex[0], &normalIndex[0],
				&vertexIndex[1], &uvIndex[1], &normalIndex[1],
				&vertexIndex[2], &uvIndex[2], &normalIndex[2],
			)

			uvIndices = append(uvIndices,
				uvIndex[0], uvIndex[1], uvIndex[2])
			vertexIndices = append(vertexIndices,
				vertexIndex[0], vertexIndex[1], vertexIndex[2])
			normalIndices = append(normalIndices,
				normalIndex[0], normalIndex[1], normalIndex[2])
		}
	}

	// Looping through the faces and getting their according vertices.
	for i := range vertexIndices {
		vertIdx := vertexIndices[i]
		// The -1 is because OBJ files for arrays start at 1 not 0.
		// So to compensate for Golang we are subtracting the index by one.
		mod.outVertices = append(mod.outVertices, tmpVertices[vertIdx-1])
	}
	for i := range uvIndices {
		vertIdx := uvIndices[i]
		// The -1 is because OBJ files for arrays start at 1 not 0.
		// So to compensate for Golang we are subtracting the index by one.
		mod.outUVs = append(mod.outUVs, tmpUVs[vertIdx-1])
	}
	for i := range normalIndices {
		vertIdx := normalIndices[i]
		// The -1 is because OBJ files for arrays start at 1 not 0.
		// So to compensate for Golang we are subtracting the index by one.
		mod.outNormals = append(mod.outNormals, tmpNormals[vertIdx-1])
	}

	return mod, nil
}
