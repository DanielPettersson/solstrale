package hittable

import (
	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/material"
	"github.com/udhos/gwob"
)

type objModel struct {
}

func NewObjModel(path string) (Hittable, error) {

	options := &gwob.ObjParserOptions{}
	object, err := gwob.NewObjFromFile(path, options)
	if err != nil {
		return nil, err
	}

	triangles := make([]Hittable, 0, object.NumberOfElements())

	whiteMat := material.Lambertian{
		Tex: material.SolidColor{
			ColorValue: geo.NewVec3(.8, .8, .8),
		},
	}

	numIndices := len(object.Indices)
	for i := 0; i < numIndices; i += 3 {

		x1, y1, z1 := object.VertexCoordinates(object.Indices[i])
		x2, y2, z2 := object.VertexCoordinates(object.Indices[i+1])
		x3, y3, z3 := object.VertexCoordinates(object.Indices[i+2])
		v0 := geo.NewVec3(float64(x1), float64(y1), float64(z1))
		v1 := geo.NewVec3(float64(x2), float64(y2), float64(z2))
		v2 := geo.NewVec3(float64(x3), float64(y3), float64(z3))
		t := NewTriangle(v0, v1, v2, whiteMat)
		triangles = append(triangles, t)
	}

	return NewBoundingVolumeHierarchy(triangles), nil
}
