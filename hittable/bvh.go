package hittable

import (
	"math"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/material"
)

// Bounding Volume Hierarchy
type bvh struct {
	NonPdfLightHittable
	left  *Hittable
	right *Hittable
	bBox  aabb
}

// NewBoundingVolumeHierarchy creates a new hittable object from the given hittable list
// The bounding Volume Hierarchy sorts the hittables in a binary tree
// where each node has a bounding box.
// This is to optimize the ray intersection search when having many hittable objects.
func NewBoundingVolumeHierarchy(list []Hittable) Hittable {

	if len(list) == 0 {
		panic("Cannot create a Bvh with empty list of objects")
	}

	return createBvh(list, 0, len(list))
}

func createBvh(list []Hittable, start, end int) *bvh {
	numObjects := end - start
	var left Hittable
	var right Hittable
	var bBox aabb

	if numObjects == 1 {
		left = list[start]
		right = list[start]
		bBox = left.BoundingBox()

	} else if numObjects == 2 {
		left = list[start]
		right = list[start+1]
		bBox = combineAabbs(left.BoundingBox(), right.BoundingBox())

	} else {
		mid := sortHittablesSliceByMostSpreadAxis(list, start, end)

		// Could not split with objects on both sides. Just split up the middle index
		if mid == start || mid == end {
			mid = start + numObjects/2
		}

		left = createBvh(list, start, mid)
		right = createBvh(list, mid, end)
		bBox = combineAabbs(left.BoundingBox(), right.BoundingBox())
	}

	return &bvh{left: &left, right: &right, bBox: bBox}
}

func xAxis(c geo.Vec3) float64 {
	return c.X
}

func yAxis(c geo.Vec3) float64 {
	return c.Y
}

func zAxis(c geo.Vec3) float64 {
	return c.Z
}

func sortHittablesSliceByMostSpreadAxis(list []Hittable, start, end int) int {
	slice := list[start:end]

	xSpread, xCenter := boundingBoxSpread(slice, xAxis)
	ySpread, yCenter := boundingBoxSpread(slice, yAxis)
	zSpread, zCenter := boundingBoxSpread(slice, zAxis)

	if xSpread >= ySpread && xSpread >= zSpread {
		return SortHittablesByCenter(slice, xCenter, xAxis) + start
	} else if ySpread >= xSpread && ySpread >= zSpread {
		return SortHittablesByCenter(slice, yCenter, yAxis) + start
	} else {
		return SortHittablesByCenter(slice, zCenter, zAxis) + start
	}
}

func boundingBoxSpread(list []Hittable, axisFunc func(h geo.Vec3) float64) (float64, float64) {
	min := util.Infinity
	max := -util.Infinity
	for _, h := range list {
		min = math.Min(min, axisFunc(h.Center()))
		max = math.Max(max, axisFunc(h.Center()))
	}
	return max - min, (min + max) * .5
}

func SortHittablesByCenter(list []Hittable, center float64, axisFunc func(c geo.Vec3) float64) int {

	i := 0
	j := len(list) - 1

	for i <= j {
		if axisFunc(list[i].Center()) < center {
			i++
		} else {
			tmpI := list[i]
			tmpJ := list[j]
			list[i] = tmpJ
			list[j] = tmpI
			j--
		}
	}

	return i
}

func (b *bvh) Hit(r geo.Ray, rayLength util.Interval) (bool, *material.HitRecord) {
	if !b.bBox.hit(r, rayLength) {
		return false, nil
	}

	hitLeft, rec := (*b.left).Hit(r, rayLength)
	if hitLeft {
		rayLength = util.Interval{Min: rayLength.Min, Max: rec.RayLength}
	}

	hitRight, recRight := (*b.right).Hit(r, rayLength)
	if hitRight {
		rec = recRight
	}

	return hitLeft || hitRight, rec
}

func (b *bvh) BoundingBox() aabb {
	return b.bBox
}

func (b *bvh) IsLight() bool {
	return false
}

func (b *bvh) Center() geo.Vec3 {
	return b.bBox.center()
}
