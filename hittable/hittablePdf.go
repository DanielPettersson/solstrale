package hittable

import (
	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/pdf"
)

// HittablePdf is a wrapper for generating pdfs for a list of hittables
type HittablePdf struct {
	objects *HittableList
	origin  geo.Vec3
}

// NewHittablePdf creates a new instance of HittablePdf
func NewHittablePdf(objects *HittableList, origin geo.Vec3) pdf.Pdf {
	return HittablePdf{
		objects: objects,
		origin:  origin,
	}
}

// Value implements pdf.Pdf
func (p HittablePdf) Value(direction geo.Vec3) float64 {
	return p.objects.PdfValue(p.origin, direction)
}

// Generate implements pdf.Pdf
func (p HittablePdf) Generate() geo.Vec3 {
	return p.objects.RandomDirection(p.origin)
}
