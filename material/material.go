package material

import (
	"math"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
)

type Material interface {
	Scatter(rayIn geo.Ray, rec *HitRecord) (bool, geo.Vec3, geo.Ray)
	Emitted(rec *HitRecord) geo.Vec3
}

type Lambertian struct {
	Tex Texture
}

func (m Lambertian) Scatter(rayIn geo.Ray, rec *HitRecord) (bool, geo.Vec3, geo.Ray) {
	scatterDirection := rec.Normal.Add(geo.RandomUnitVector())
	if scatterDirection.NearZero() {
		scatterDirection = rec.Normal
	}

	scatterRay := geo.Ray{
		Origin:    rec.HitPoint,
		Direction: scatterDirection,
		Time:      rayIn.Time,
	}
	return true, m.Tex.Color(rec), scatterRay
}

func (m Lambertian) Emitted(rec *HitRecord) geo.Vec3 {
	return geo.ZeroVector
}

type Metal struct {
	Tex  Texture
	Fuzz float64
}

func (m Metal) Scatter(rayIn geo.Ray, rec *HitRecord) (bool, geo.Vec3, geo.Ray) {
	reflected := rayIn.Direction.Unit().Reflect(rec.Normal)
	scatterRay := geo.Ray{
		Origin:    rec.HitPoint,
		Direction: reflected.Add(geo.RandomInUnitSphere().MulS(m.Fuzz)),
		Time:      rayIn.Time,
	}
	scatter := scatterRay.Direction.Dot(rec.Normal) > 0

	return scatter, m.Tex.Color(rec), scatterRay
}

func (m Metal) Emitted(rec *HitRecord) geo.Vec3 {
	return geo.ZeroVector
}

type Dielectric struct {
	Tex               Texture
	IndexOfRefraction float64
}

func (m Dielectric) Scatter(rayIn geo.Ray, rec *HitRecord) (bool, geo.Vec3, geo.Ray) {
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
	if cannotRefract || reflectance(cosTheta, refractionRatio) > util.RandomNormalFloat() {
		direction = unitDirection.Reflect(rec.Normal)
	} else {
		direction = unitDirection.Refract(rec.Normal, refractionRatio)
	}

	scatter := geo.Ray{
		Origin:    rec.HitPoint,
		Direction: direction,
		Time:      rayIn.Time,
	}

	return true, m.Tex.Color(rec), scatter
}

func (m Dielectric) Emitted(rec *HitRecord) geo.Vec3 {
	return geo.ZeroVector
}

// Calculate reflectance using Schlick's approximation
func reflectance(cosine, indexOfRefraction float64) float64 {
	r0 := (1 - indexOfRefraction) / (1 + indexOfRefraction)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}

type DiffuseLight struct {
	Emit Texture
}

func (m DiffuseLight) Scatter(rayIn geo.Ray, rec *HitRecord) (bool, geo.Vec3, geo.Ray) {
	return false, geo.Vec3{}, geo.Ray{}
}

func (m DiffuseLight) Emitted(rec *HitRecord) geo.Vec3 {
	return m.Emit.Color(rec)
}

type Isotropic struct {
	Albedo Texture
}

func (m Isotropic) Scatter(rayIn geo.Ray, rec *HitRecord) (bool, geo.Vec3, geo.Ray) {
	attenuation := m.Albedo.Color(rec)
	scattered := geo.Ray{
		Origin:    rec.HitPoint,
		Direction: geo.RandomUnitVector(),
		Time:      rayIn.Time,
	}
	return true, attenuation, scattered
}

func (m Isotropic) Emitted(rec *HitRecord) geo.Vec3 {
	return geo.ZeroVector
}
