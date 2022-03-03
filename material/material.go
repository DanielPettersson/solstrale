package material

import (
	"math"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/pdf"
	"github.com/DanielPettersson/solstrale/random"
)

type ScatterRecord struct {
	Attenuation geo.Vec3
	PdfPtr      *pdf.Pdf
	SkipPdf     bool
	SkipPdfRay  geo.Ray
}

// Material is the interface for types that describe how
// a ray behaves when hitting an object.
type Material interface {
	Scatter(rayIn geo.Ray, rec *HitRecord) (bool, ScatterRecord)
	ScatteringPdf(rayIn geo.Ray, rec *HitRecord, scattered geo.Ray) float64
	Emitted(rec *HitRecord) geo.Vec3
}

// Lambertian is a typical matte material
type Lambertian struct {
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

func (m Lambertian) ScatteringPdf(rayIn geo.Ray, rec *HitRecord, scattered geo.Ray) float64 {
	cosTheta := rec.Normal.Dot(scattered.Direction.Unit())
	if cosTheta < 0 {
		return 0
	} else {
		return cosTheta / math.Pi
	}
}

// Emitted a lambertian material emits no light
func (m Lambertian) Emitted(rec *HitRecord) geo.Vec3 {
	return geo.ZeroVector
}

// Metal is a material that is reflective
type Metal struct {
	Tex  Texture
	Fuzz float64
}

// Scatter returns a reflected scattered ray for the metal material
// The Fuzz property of the metal defines the randomness applied to the reflection
func (m Metal) Scatter(rayIn geo.Ray, rec *HitRecord) (bool, ScatterRecord) {
	reflected := rayIn.Direction.Unit().Reflect(rec.Normal)
	scatterRay := geo.Ray{
		Origin:    rec.HitPoint,
		Direction: reflected.Add(geo.RandomInUnitSphere().MulS(m.Fuzz)),
		Time:      rayIn.Time,
	}

	return true, ScatterRecord{
		Attenuation: m.Tex.Color(rec),
		PdfPtr:      nil,
		SkipPdf:     true,
		SkipPdfRay:  scatterRay,
	}
}

func (m Metal) ScatteringPdf(rayIn geo.Ray, rec *HitRecord, scattered geo.Ray) float64 {
	return 0
}

// Emitted a metal material emits no light
func (m Metal) Emitted(rec *HitRecord) geo.Vec3 {
	return geo.ZeroVector
}

// Dielectric is a glass type material with an index of refraction
type Dielectric struct {
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

	scatterRay := geo.Ray{
		Origin:    rec.HitPoint,
		Direction: direction,
		Time:      rayIn.Time,
	}

	return true, ScatterRecord{
		Attenuation: m.Tex.Color(rec),
		PdfPtr:      nil,
		SkipPdf:     true,
		SkipPdfRay:  scatterRay,
	}
}

func (m Dielectric) ScatteringPdf(rayIn geo.Ray, rec *HitRecord, scattered geo.Ray) float64 {
	return 0
}

// Emitted a dielectric material emits no light
func (m Dielectric) Emitted(rec *HitRecord) geo.Vec3 {
	return geo.ZeroVector
}

// Calculate reflectance using Schlick's approximation
func reflectance(cosine, indexOfRefraction float64) float64 {
	r0 := (1 - indexOfRefraction) / (1 + indexOfRefraction)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}

// DiffuseLight is a material used for emitting light
type DiffuseLight struct {
	Emit Texture
}

// Scatter a light never scatters a ray
func (m DiffuseLight) Scatter(rayIn geo.Ray, rec *HitRecord) (bool, ScatterRecord) {
	return false, ScatterRecord{}
}

func (m DiffuseLight) ScatteringPdf(rayIn geo.Ray, rec *HitRecord, scattered geo.Ray) float64 {
	return 0
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

func (m Isotropic) ScatteringPdf(rayIn geo.Ray, rec *HitRecord, scattered geo.Ray) float64 {
	return 1 / (4 * math.Pi)
}

// Emitted a isotropic material emits no light
func (m Isotropic) Emitted(rec *HitRecord) geo.Vec3 {
	return geo.ZeroVector
}
