package hittable

import (
	"math"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/material"
)

type constantMedium struct {
	Boundary               Hittable
	NegativeInverseDensity float64
	PhaseFunction          material.Material
}

func NewConstantMedium(boundary Hittable, density float64, color material.Texture) Hittable {
	return constantMedium{
		Boundary:               boundary,
		NegativeInverseDensity: -1 / density,
		PhaseFunction:          material.Isotropic{Albedo: color},
	}
}

func (cm constantMedium) Hit(r geo.Ray, rayLength util.Interval) (bool, *material.HitRecord) {

	hit1, rec1 := cm.Boundary.Hit(r, util.UniverseInterval)
	if !hit1 {
		return false, nil
	}

	hit2, rec2 := cm.Boundary.Hit(r, util.Interval{Min: rec1.RayLength + 0.0001, Max: util.Infinity})
	if !hit2 {
		return false, nil
	}

	rec1.RayLength = math.Max(rec1.RayLength, rayLength.Min)
	rec2.RayLength = math.Min(rec2.RayLength, rayLength.Max)

	if rec1.RayLength >= rec2.RayLength {
		return false, nil
	}

	rec1.RayLength = math.Max(rec1.RayLength, 0)

	rLength := r.Direction.Length()
	distanceInsideBoundary := (rec2.RayLength - rec1.RayLength) * rLength
	hitDistance := cm.NegativeInverseDensity * math.Log(util.RandomNormalFloat())

	if hitDistance > distanceInsideBoundary {
		return false, nil
	}

	t := rec1.RayLength + hitDistance/rLength

	hitRecord := material.HitRecord{
		HitPoint:  r.At(t),
		Material:  cm.PhaseFunction,
		RayLength: t,
	}

	return true, &hitRecord
}

func (cm constantMedium) BoundingBox() aabb {
	return cm.Boundary.BoundingBox()
}