package hittable

import (
	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/material"
)

// Hittable is the common interface for all objects in the ray tracing scene
// that can be hit by rays
type Hittable interface {
	PdfLightHittable
	Hit(r geo.Ray, rayLength util.Interval) (bool, *material.HitRecord)
	BoundingBox() aabb
	Center() geo.Vec3
	IsLight() bool
}

// PdfLightHittable has methods used when other hittables have pdf scattering materials.
// This should be implemented by hittables that can have light materials.
type PdfLightHittable interface {
	PdfValue(origin, direction geo.Vec3) float64
	RandomDirection(origin geo.Vec3) geo.Vec3
}

// NonPdfLightHittable is used by hittables that never uses pdfs directly
// Will panic hittable is a light and other hittable has pdf generating material
type NonPdfLightHittable struct{}

// PdfValue panics if invoked
func (h NonPdfLightHittable) PdfValue(origin, direction geo.Vec3) float64 {
	panic("Should not be used")
}

// RandomDirection panics if invoked
func (h NonPdfLightHittable) RandomDirection(origin geo.Vec3) geo.Vec3 {
	panic("Should not be used")
}
