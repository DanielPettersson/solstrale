package material

import (
	"github.com/DanielPettersson/solstrale/geo"
)

type HitRecord struct {
	HitPoint  geo.Vec3
	Normal    geo.Vec3
	Material  Material
	RayLength float64
	U         float64
	V         float64
	FrontFace bool
}
