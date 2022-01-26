package trace

import (
	"math"
	"math/rand"
)

type material interface {
	scatter(rayIn ray, rec hitRecord) (bool, vec3, ray)
	emitted(rec hitRecord) vec3
}

type lambertian struct {
	tex texture
}

func (m lambertian) scatter(rayIn ray, rec hitRecord) (bool, vec3, ray) {
	scatterDirection := rec.normal.add(randomUnitVector())
	if scatterDirection.nearZero() {
		scatterDirection = rec.normal
	}

	scatterRay := ray{rec.hitPoint, scatterDirection, rayIn.time}
	return true, m.tex.color(rec), scatterRay
}

func (m lambertian) emitted(rec hitRecord) vec3 {
	return black
}

type metal struct {
	tex  texture
	fuzz float64
}

func (m metal) scatter(rayIn ray, rec hitRecord) (bool, vec3, ray) {
	reflected := rayIn.direction.unit().reflect(rec.normal)
	scatterRay := ray{rec.hitPoint, reflected.add(randomInUnitSphere().mulS(m.fuzz)), rayIn.time}
	scatter := scatterRay.direction.dot(rec.normal) > 0

	return scatter, m.tex.color(rec), scatterRay
}

func (m metal) emitted(rec hitRecord) vec3 {
	return black
}

type dielectric struct {
	tex               texture
	indexOfRefraction float64
}

func (m dielectric) scatter(rayIn ray, rec hitRecord) (bool, vec3, ray) {
	var refractionRatio float64
	if rec.frontFace {
		refractionRatio = 1 / m.indexOfRefraction
	} else {
		refractionRatio = m.indexOfRefraction
	}

	unitDirection := rayIn.direction.unit()
	cosTheta := math.Min(unitDirection.neg().dot(rec.normal), 1)
	sinTheta := math.Sqrt(1 - cosTheta*cosTheta)
	cannotRefract := refractionRatio*sinTheta > 1

	var direction vec3
	if cannotRefract || reflectance(cosTheta, refractionRatio) > rand.Float64() {
		direction = unitDirection.reflect(rec.normal)
	} else {
		direction = unitDirection.refract(rec.normal, refractionRatio)
	}

	scatter := ray{rec.hitPoint, direction, rayIn.time}

	return true, m.tex.color(rec), scatter
}

func (m dielectric) emitted(rec hitRecord) vec3 {
	return black
}

// Calculate reflectance using Schlick's approximation
func reflectance(cosine, indexOfRefraction float64) float64 {
	r0 := (1 - indexOfRefraction) / (1 + indexOfRefraction)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}

type diffuseLight struct {
	emit texture
}

func (m diffuseLight) scatter(rayIn ray, rec hitRecord) (bool, vec3, ray) {
	return false, vec3{}, ray{}
}

func (m diffuseLight) emitted(rec hitRecord) vec3 {
	return m.emit.color(rec)
}

type isotropic struct {
	albedo texture
}

func (m isotropic) scatter(rayIn ray, rec hitRecord) (bool, vec3, ray) {
	attenuation := m.albedo.color(rec)
	scattered := ray{rec.hitPoint, randomUnitVector(), rayIn.time}
	return true, attenuation, scattered
}

func (m isotropic) emitted(rec hitRecord) vec3 {
	return black
}
