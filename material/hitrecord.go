// Package material provides material and texture related types and functions
package material

import (
	"github.com/DanielPettersson/solstrale/geo"
)

// HitRecord is a collection of all interesting properties from
// when a ray hits a hittable object
type HitRecord struct {
	HitPoint  geo.Vec3
	Normal    geo.Vec3
	Material  Material
	RayLength float64
	U         float64
	V         float64
	FrontFace bool
}
