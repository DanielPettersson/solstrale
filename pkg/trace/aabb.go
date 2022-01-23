package trace

import "math"

// Axis Aligned Bounding Box
type aabb struct {
	x, y, z interval
}

func createAabbFromPoints(a vec3, b vec3) aabb {
	return aabb{
		interval{math.Min(a.x, b.x), math.Max(a.x, b.x)},
		interval{math.Min(a.y, b.y), math.Max(a.y, b.y)},
		interval{math.Min(a.z, b.z), math.Max(a.z, b.z)},
	}
}

func combineAabbs(a aabb, b aabb) aabb {
	return aabb{
		combineIntervals(a.x, b.x),
		combineIntervals(a.y, b.y),
		combineIntervals(a.z, b.z),
	}
}

func (a aabb) add(offset vec3) aabb {
	return aabb{
		a.x.add(offset.x),
		a.y.add(offset.y),
		a.z.add(offset.z),
	}
}

func (a aabb) padIfNeeded() aabb {
	delta := 0.0001
	var newX interval
	if a.x.size() >= delta {
		newX = a.x
	} else {
		newX = a.x.expand(delta)
	}
	var newY interval
	if a.y.size() >= delta {
		newY = a.y
	} else {
		newY = a.y.expand(delta)
	}
	var newZ interval
	if a.z.size() >= delta {
		newZ = a.z
	} else {
		newZ = a.z.expand(delta)
	}

	return aabb{newX, newY, newZ}
}

func (aabb aabb) hit(r ray, rayLength interval) bool {

	tMin := (aabb.x.min - r.origin.x) / r.direction.x
	tMax := (aabb.x.max - r.origin.x) / r.direction.x

	t0 := math.Min(tMin, tMax)
	t1 := math.Max(tMin, tMax)
	rayLengthMin := math.Max(t0, rayLength.min)
	rayLengthMax := math.Min(t1, rayLength.max)

	if rayLengthMax <= rayLengthMin {
		return false
	}

	tMin = (aabb.y.min - r.origin.y) / r.direction.y
	tMax = (aabb.y.max - r.origin.y) / r.direction.y

	t0 = math.Min(tMin, tMax)
	t1 = math.Max(tMin, tMax)
	rayLengthMin = math.Max(t0, rayLength.min)
	rayLengthMax = math.Min(t1, rayLength.max)

	if rayLengthMax <= rayLengthMin {
		return false
	}

	tMin = (aabb.z.min - r.origin.z) / r.direction.z
	tMax = (aabb.z.max - r.origin.z) / r.direction.z

	t0 = math.Min(tMin, tMax)
	t1 = math.Max(tMin, tMax)
	rayLengthMin = math.Max(t0, rayLength.min)
	rayLengthMax = math.Min(t1, rayLength.max)

	return rayLengthMax > rayLengthMin
}
