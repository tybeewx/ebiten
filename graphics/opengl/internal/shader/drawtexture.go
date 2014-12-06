package shader

import (
	"github.com/go-gl/gl"
	"github.com/hajimehoshi/ebiten/graphics"
	"github.com/hajimehoshi/ebiten/graphics/matrix"
	"sync"
)

var once sync.Once

func DrawTexture(native gl.Texture, projectionMatrix [16]float32, quads []graphics.TextureQuad, geo matrix.Geometry, color matrix.Color) {
	once.Do(func() {
		initialize()
	})

	if len(quads) == 0 {
		return
	}
	// TODO: Check performance
	shaderProgram := use(projectionMatrix, geo, color)

	native.Bind(gl.TEXTURE_2D)
	defer gl.Texture(0).Bind(gl.TEXTURE_2D)

	vertexAttrLocation := getAttributeLocation(shaderProgram, "vertex")
	texCoordAttrLocation := getAttributeLocation(shaderProgram, "tex_coord")

	gl.EnableClientState(gl.VERTEX_ARRAY)
	gl.EnableClientState(gl.TEXTURE_COORD_ARRAY)
	vertexAttrLocation.EnableArray()
	texCoordAttrLocation.EnableArray()
	defer func() {
		texCoordAttrLocation.DisableArray()
		vertexAttrLocation.DisableArray()
		gl.DisableClientState(gl.TEXTURE_COORD_ARRAY)
		gl.DisableClientState(gl.VERTEX_ARRAY)
	}()

	vertices := []float32{}
	texCoords := []float32{}
	indicies := []uint32{}
	// TODO: Check len(parts) and GL_MAX_ELEMENTS_INDICES
	for i, quad := range quads {
		x1 := quad.VertexX1
		x2 := quad.VertexX2
		y1 := quad.VertexY1
		y2 := quad.VertexY2
		vertices = append(vertices,
			x1, y1,
			x2, y1,
			x1, y2,
			x2, y2,
		)
		u1 := quad.TextureCoordU1
		u2 := quad.TextureCoordU2
		v1 := quad.TextureCoordV1
		v2 := quad.TextureCoordV2
		texCoords = append(texCoords,
			u1, v1,
			u2, v1,
			u1, v2,
			u2, v2,
		)
		base := uint32(i * 4)
		indicies = append(indicies,
			base, base+1, base+2,
			base+1, base+2, base+3,
		)
	}
	vertexAttrLocation.AttribPointer(2, gl.FLOAT, false, 0, vertices)
	texCoordAttrLocation.AttribPointer(2, gl.FLOAT, false, 0, texCoords)
	gl.DrawElements(gl.TRIANGLES, len(indicies), gl.UNSIGNED_INT, indicies)
}
