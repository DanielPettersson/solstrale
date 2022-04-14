package material

import (
	"math"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/pdf"
	"github.com/DanielPettersson/solstrale/random"
)

// ScatterRecord is a collection of attributes from the scattering of a ray with a material
type ScatterRecord struct {
	Attenuation geo.Vec3
	PdfPtr      *pdf.Pdf
	SkipPdf     bool
	SkipPdfRay  geo.Ray
}

// Material is the interface for types that describe how
// a ray behaves when hitting an object.
type Material interface {
	PdfGeneratingMaterial
	LightEmittingMaterial
	Scatter(rayIn geo.Ray, rec *HitRecord) (bool, ScatterRecord)
}

// PdfGeneratingMaterial is a material that can use pdfs for scattering of rays
type PdfGeneratingMaterial interface {
	ScatteringPdf(rayIn geo.Ray, rec *HitRecord, scattered geo.Ray) float64
}

// NonPdfGeneratingMaterial is to be used by materials that do not use pdfs
type NonPdfGeneratingMaterial struct{}

// ScatteringPdf as NonPdfGeneratingMaterial is used for materials that do not generate a pdf
// Just return 0
func (m NonPdfGeneratingMaterial) ScatteringPdf(rayIn geo.Ray, rec *HitRecord, scattered geo.Ray) float64 {
	return 0
}

// LightEmittingMaterial is a material that can emit light
type LightEmittingMaterial interface {
	Emitted(rec *HitRecord) geo.Vec3
}

// NonLightEmittingMaterial is to be used by materials that do not emit light
type NonLightEmittingMaterial struct{}

// Emitted a non emitting material emits zero light
func (m NonLightEmittingMaterial) Emitted(rec *HitRecord) geo.Vec3 {
	return geo.ZeroVector
}

// Lambertian is a typical matte material
type Lambertian struct {
	NonLightEmittingMaterial
	Tex Texture
}

// Scatter returns a randomish scatter of the ray for the matte material
func (m Lambertian) Scatter(rayIn geo.Ray, rec *HitRecord) (bool, ScatterRecord) {
	attenuation := m.Tex.Color(rec)
	pdfPtr := pdf.NewCosinePdf(rec.Normal)

	return true, ScatterRecord{
		Attenuation: attenuation,
		PdfPtr:      &pdfPtr,
		SkipPdf:     false,
	}
}

// ScatteringPdf returns the pdf value for a given rays for the lambertian material
func (m Lambertian) ScatteringPdf(rayIn geo.Ray, rec *HitRecord, scattered geo.Ray) float64 {
	cosTheta := rec.Normal.Dot(scattered.Direction.Unit())
	if cosTheta < 0 {
		return 0
	}
	return cosTheta / math.Pi
}

// Metal is a material that is reflective
type Metal struct {
	NonLightEmittingMaterial
	NonPdfGeneratingMaterial
	Tex  Texture
	Fuzz float64
}

// Scatter returns a reflected scattered ray for the metal material
// The Fuzz property of the metal defines the randomness applied to the reflection
func (m Metal) Scatter(rayIn geo.Ray, rec *HitRecord) (bool, ScatterRecord) {
	reflected := rayIn.Direction.Unit().Reflect(rec.Normal)
	scatterRay := geo.NewRay(
		rec.HitPoint,
		reflected.Add(geo.RandomInUnitSphere().MulS(m.Fuzz)),
		rayIn.Time,
	)

	return true, ScatterRecord{
		Attenuation: m.Tex.Color(rec),
		PdfPtr:      nil,
		SkipPdf:     true,
		SkipPdfRay:  scatterRay,
	}
}

// Dielectric is a glass type material with an index of refraction
type Dielectric struct {
	NonLightEmittingMaterial
	NonPdfGeneratingMaterial
	Tex               Texture
	IndexOfRefraction float64
}

// Scatter returns a refracted ray for the dielectric material
func (m Dielectric) Scatter(rayIn geo.Ray, rec *HitRecord) (bool, ScatterRecord) {
	var refractionRatio float64
	if rec.FrontFace {
		refractionRatio = 1 / m.IndexOfRefraction
	} else {
		refractionRatio = m.IndexOfRefraction
	}

	unitDirection := rayIn.Direction.Unit()
	cosTheta := math.Min(unitDirection.Neg().Dot(rec.Normal), 1)
	sinTheta := math.Sqrt(1 - cosTheta*cosTheta)
	cannotRefract := refractionRatio*sinTheta > 1

	var direction geo.Vec3
	if cannotRefract || reflectance(cosTheta, refractionRatio) > random.RandomNormalFloat() {
		direction = unitDirection.Reflect(rec.Normal)
	} else {
		direction = unitDirection.Refract(rec.Normal, refractionRatio)
	}

	scatterRay := geo.NewRay(
		rec.HitPoint,
		direction,
		rayIn.Time,
	)

	return true, ScatterRecord{
		Attenuation: m.Tex.Color(rec),
		PdfPtr:      nil,
		SkipPdf:     true,
		SkipPdfRay:  scatterRay,
	}
}

// Calculate reflectance using Schlick's approximation
func reflectance(cosine, indexOfRefraction float64) float64 {
	r0 := (1 - indexOfRefraction) / (1 + indexOfRefraction)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}

// DiffuseLight is a material used for emitting light
type DiffuseLight struct {
	NonPdfGeneratingMaterial
	Emit Texture
}

// Scatter a light never scatters a ray
func (m DiffuseLight) Scatter(rayIn geo.Ray, rec *HitRecord) (bool, ScatterRecord) {
	return false, ScatterRecord{}
}

// Emitted a light emits it's given color
func (m DiffuseLight) Emitted(rec *HitRecord) geo.Vec3 {
	if !rec.FrontFace {
		return geo.ZeroVector
	}
	return m.Emit.Color(rec)
}

// Isotropic is a fog type material
type Isotropic struct {
	NonLightEmittingMaterial
	Albedo Texture
}

// Scatter returns a randomly scattered ray in any direction
func (m Isotropic) Scatter(rayIn geo.Ray, rec *HitRecord) (bool, ScatterRecord) {
	attenuation := m.Albedo.Color(rec)
	pdfPtr := pdf.NewSpherePdf()

	return true, ScatterRecord{
		Attenuation: attenuation,
		PdfPtr:      &pdfPtr,
		SkipPdf:     false,
	}
}

// ScatteringPdf returns the pdf value for a given rays for the isotropic material
func (m Isotropic) ScatteringPdf(rayIn geo.Ray, rec *HitRecord, scattered geo.Ray) float64 {
	return 1 / (4 * math.Pi)
}
