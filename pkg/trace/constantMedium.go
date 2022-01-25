package trace

import (
	"math"
	"math/rand"
)

type constantMedium struct {
	boundary               hittable
	negativeInverseDensity float64
	phaseFunction          material
}

func createConstantMedium(boundary hittable, density float64, color texture) constantMedium {
	return constantMedium{
		boundary:               boundary,
		negativeInverseDensity: -1 / density,
		phaseFunction:          isotropic{color},
	}
}

func (cm constantMedium) hit(r ray, rayLength interval) (bool, *hitRecord) {

	hit1, rec1 := cm.boundary.hit(r, universe_interval)
	if !hit1 {
		return false, nil
	}

	hit2, rec2 := cm.boundary.hit(r, interval{rec1.rayLength + 0.0001, infinity})
	if !hit2 {
		return false, nil
	}

	rec1.rayLength = math.Max(rec1.rayLength, rayLength.min)
	rec2.rayLength = math.Min(rec2.rayLength, rayLength.max)

	if rec1.rayLength >= rec2.rayLength {
		return false, nil
	}

	rec1.rayLength = math.Max(rec1.rayLength, 0)

	rLength := r.direction.length()
	distanceInsideBoundary := (rec2.rayLength - rec1.rayLength) * rLength
	hitDistance := cm.negativeInverseDensity * math.Log(rand.Float64())

	if hitDistance > distanceInsideBoundary {
		return false, nil
	}

	t := rec1.rayLength + hitDistance/rLength

	hitRecord := hitRecord{
		hitPoint:  r.at(t),
		material:  cm.phaseFunction,
		rayLength: t,
	}

	return true, &hitRecord
}

func (cm constantMedium) boundingBox() aabb {
	return cm.boundary.boundingBox()
}
