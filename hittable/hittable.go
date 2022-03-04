package hittable

import (
	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/material"
)

// Hittable is the common interface for all objects in the ray tracing scene
// that can be hit by rays
type Hittable interface {
	PdfUsingHittable
	Hit(r geo.Ray, rayLength util.Interval) (bool, *material.HitRecord)
	BoundingBox() aabb
}

type PdfUsingHittable interface {
	PdfValue(origin, direction geo.Vec3) float64
	RandomDirection(origin geo.Vec3) geo.Vec3
}

type NonPdfUsingHittable struct{}

func (h NonPdfUsingHittable) PdfValue(o, v geo.Vec3) float64 {
	panic("Should not be used")
}

func (h NonPdfUsingHittable) RandomDirection(o geo.Vec3) geo.Vec3 {
	panic("Should not be used")
}
