package hittable

import (
	"errors"
	"fmt"
	"image"
	"os"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/material"
	"github.com/udhos/gwob"
)

// NewObjModel reads a Wavefront .obj file and creates a bvh containing
// all triangles. It also read materials from the referred .mat file.
// Support for colored and textured lambertian materials.
func NewObjModel(path string) (Hittable, error) {
	return NewObjModelWithDefaultMaterial(
		path,
		material.Lambertian{Tex: material.SolidColor{ColorValue: geo.NewVec3(1, 1, 1)}},
	)
}

// NewObjModelWithDefaultMaterial reads a Wavefront .obj file and creates a bvh containing
// all triangles. It also read materials from the referred .mat file.
// Support for colored and textured lambertian materials.
// Applies supplied default material if none in model
func NewObjModelWithDefaultMaterial(path string, defaultMaterial material.Material) (Hittable, error) {

	options := &gwob.ObjParserOptions{IgnoreNormals: true}
	object, err := gwob.NewObjFromFile(path, options)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to parse obj file: %v", err.Error()))
	}

	mats := map[string]material.Material{
		"_": defaultMaterial,
	}

	// Read all materials if a library is defined

	if object.Mtllib != "" {
		materialLib, err := gwob.ReadMaterialLibFromFile(object.Mtllib, options)
		if err != nil {
			return nil, err
		}

		for name, m := range materialLib.Lib {

			// If a texture
			if m.MapKd != "" {

				f, _ := os.Open(m.MapKd)
				defer f.Close()
				image, _, err := image.Decode(f)
				if err != nil {
					return nil, err
				}
				mats[name] = material.Lambertian{Tex: material.ImageTexture{Image: image}}

				// Otherwise use the diffuse color for a lambertian
			} else {

				mats[name] = material.Lambertian{Tex: material.SolidColor{ColorValue: geo.NewVec3(
					float64(m.Kd[0]),
					float64(m.Kd[1]),
					float64(m.Kd[2]),
				)}}
			}
		}
	}

	triangles := make([]Hittable, 0, object.NumberOfElements())

	for _, group := range object.Groups {

		// For each group in object, read all triangles and set material

		for i := group.IndexBegin; i < group.IndexBegin+group.IndexCount; i += 3 {

			mat, found := mats[group.Usemtl]
			if !found {
				mat = mats["_"]
			}

			x, y, z := object.VertexCoordinates(object.Indices[i])
			v0 := geo.NewVec3(float64(x), float64(y), float64(z))
			x, y, z = object.VertexCoordinates(object.Indices[i+1])
			v1 := geo.NewVec3(float64(x), float64(y), float64(z))
			x, y, z = object.VertexCoordinates(object.Indices[i+2])
			v2 := geo.NewVec3(float64(x), float64(y), float64(z))

			var tu0, tv0, tu1, tv1, tu2, tv2 float64

			// Read texture coordinates if any

			if object.TextCoordFound {
				tu0, tv0 = textureCoordinates(*object, object.Indices[i])
				tu1, tv1 = textureCoordinates(*object, object.Indices[i+1])
				tu2, tv2 = textureCoordinates(*object, object.Indices[i+2])
			}

			t := NewTriangleWithTexCoords(v0, v1, v2, tu0, tv0, tu1, tv1, tu2, tv2, mat)
			triangles = append(triangles, t)
		}
	}

	return NewBoundingVolumeHierarchy(triangles), nil
}

func textureCoordinates(o gwob.Obj, stride int) (float64, float64) {
	offset := o.StrideOffsetTexture / 4
	floatsPerStride := o.StrideSize / 4
	f := offset + stride*floatsPerStride
	return float64(o.Coord[f]), float64(o.Coord[f+1])
}
