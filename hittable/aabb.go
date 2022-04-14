package hittable

import (
	"math"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
)

// Axis Aligned Bounding Box
type aabb struct {
	x, y, z util.Interval
}

func createAabbFromPoints(a geo.Vec3, b geo.Vec3) aabb {
	return aabb{
		util.Interval{Min: math.Min(a.X, b.X), Max: math.Max(a.X, b.X)},
		util.Interval{Min: math.Min(a.Y, b.Y), Max: math.Max(a.Y, b.Y)},
		util.Interval{Min: math.Min(a.Z, b.Z), Max: math.Max(a.Z, b.Z)},
	}
}

func combineAabbs(a aabb, b aabb) aabb {
	return aabb{
		util.CombineIntervals(a.x, b.x),
		util.CombineIntervals(a.y, b.y),
		util.CombineIntervals(a.z, b.z),
	}
}

func (b aabb) add(offset geo.Vec3) aabb {
	return aabb{
		b.x.Add(offset.X),
		b.y.Add(offset.Y),
		b.z.Add(offset.Z),
	}
}

func (b aabb) padIfNeeded() aabb {
	delta := 0.0001
	var newX util.Interval
	if b.x.Size() >= delta {
		newX = b.x
	} else {
		newX = b.x.Expand(delta)
	}
	var newY util.Interval
	if b.y.Size() >= delta {
		newY = b.y
	} else {
		newY = b.y.Expand(delta)
	}
	var newZ util.Interval
	if b.z.Size() >= delta {
		newZ = b.z
	} else {
		newZ = b.z.Expand(delta)
	}

	return aabb{newX, newY, newZ}
}

func (b aabb) hit(r geo.Ray, rayLength util.Interval) bool {

	t1 := (b.x.Min - r.Origin.X) * r.DirectionInverted.X
	t2 := (b.x.Max - r.Origin.X) * r.DirectionInverted.X

	tmin := math.Min(t1, t2)
	tmax := math.Max(t1, t2)

	t1 = (b.y.Min - r.Origin.Y) * r.DirectionInverted.Y
	t2 = (b.y.Max - r.Origin.Y) * r.DirectionInverted.Y

	tmin = math.Max(tmin, math.Min(t1, t2))
	tmax = math.Min(tmax, math.Max(t1, t2))

	t1 = (b.z.Min - r.Origin.Z) * r.DirectionInverted.Z
	t2 = (b.z.Max - r.Origin.Z) * r.DirectionInverted.Z

	tmin = math.Max(tmin, math.Min(t1, t2))
	tmax = math.Min(tmax, math.Max(t1, t2))

	return tmax > math.Max(tmin, 0)
}
