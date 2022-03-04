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
	IsLight() bool
}

// PdfUsingHittable has methods used by hittables that can use pdf for scattering
type PdfUsingHittable interface {
	PdfValue(origin, direction geo.Vec3) float64
	RandomDirection(origin geo.Vec3) geo.Vec3
}

// NonPdfUsingHittable is used by hittables that never uses pdfs directly
type NonPdfUsingHittable struct{}

// PdfValue panics if invoked
func (h NonPdfUsingHittable) PdfValue(origin, direction geo.Vec3) float64 {
	panic("Should not be used")
}

// RandomDirection panics if invoked
func (h NonPdfUsingHittable) RandomDirection(origin geo.Vec3) geo.Vec3 {
	panic("Should not be used")
}
