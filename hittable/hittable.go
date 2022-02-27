package hittable

import (
	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/material"
)

// Hittable is the common interface for all objects in the ray tracing scene
// that can be hit by rays
type Hittable interface {
	Hit(r geo.Ray, rayLength util.Interval) (bool, *material.HitRecord)
	BoundingBox() aabb
}
