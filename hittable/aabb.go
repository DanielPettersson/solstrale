package hittable

import (
	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
)

// Axis Aligned Bounding Box
type aabb struct {
	x, y, z util.Interval
}

func createAabbFromPoints(a geo.Vec3, b geo.Vec3) aabb {
	return aabb{
		util.Interval{Min: min(a.X, b.X), Max: max(a.X, b.X)},
		util.Interval{Min: min(a.Y, b.Y), Max: max(a.Y, b.Y)},
		util.Interval{Min: min(a.Z, b.Z), Max: max(a.Z, b.Z)},
	}
}

func combineAabbs(a aabb, b aabb) aabb {
	return aabb{
		util.CombineIntervals(a.x, b.x),
		util.CombineIntervals(a.y, b.y),
		util.CombineIntervals(a.z, b.z),
	}
}

func (a aabb) add(offset geo.Vec3) aabb {
	return aabb{
		a.x.Add(offset.X),
		a.y.Add(offset.Y),
		a.z.Add(offset.Z),
	}
}

func (a aabb) padIfNeeded() aabb {
	delta := 0.0001
	var newX util.Interval
	if a.x.Size() >= delta {
		newX = a.x
	} else {
		newX = a.x.Expand(delta)
	}
	var newY util.Interval
	if a.y.Size() >= delta {
		newY = a.y
	} else {
		newY = a.y.Expand(delta)
	}
	var newZ util.Interval
	if a.z.Size() >= delta {
		newZ = a.z
	} else {
		newZ = a.z.Expand(delta)
	}

	return aabb{newX, newY, newZ}
}

func (a aabb) hit(r geo.Ray, rayLength util.Interval) bool {

	tMin := (a.x.Min - r.Origin.X) / r.Direction.X
	tMax := (a.x.Max - r.Origin.X) / r.Direction.X

	t0 := min(tMin, tMax)
	t1 := max(tMin, tMax)
	rayLengthMin := max(t0, rayLength.Min)
	rayLengthMax := min(t1, rayLength.Max)

	if rayLengthMax <= rayLengthMin {
		return false
	}

	tMin = (a.y.Min - r.Origin.Y) / r.Direction.Y
	tMax = (a.y.Max - r.Origin.Y) / r.Direction.Y

	t0 = min(tMin, tMax)
	t1 = max(tMin, tMax)
	rayLengthMin = max(t0, rayLength.Min)
	rayLengthMax = min(t1, rayLength.Max)

	if rayLengthMax <= rayLengthMin {
		return false
	}

	tMin = (a.z.Min - r.Origin.Z) / r.Direction.Z
	tMax = (a.z.Max - r.Origin.Z) / r.Direction.Z

	t0 = min(tMin, tMax)
	t1 = max(tMin, tMax)
	rayLengthMin = max(t0, rayLength.Min)
	rayLengthMax = min(t1, rayLength.Max)

	return rayLengthMax > rayLengthMin
}

func min(v0, v1 float64) float64 {
	if v0 < v1 {
		return v0
	}
	return v1
}

func max(v0, v1 float64) float64 {
	if v0 > v1 {
		return v0
	}
	return v1
}
