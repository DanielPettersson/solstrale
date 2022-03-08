package hittable

import (
	"math"
	"sort"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/material"
)

// Bounding Volume Hierarchy
type bvh struct {
	NonPdfUsingHittable
	left  *Hittable
	right *Hittable
	bBox  *aabb
}

// NewBoundingVolumeHierarchy creates a new hittable object from the given hittable list
// The bounding Volume Hierarchy sorts the hittables in a binary tree
// where each node has a bounding box.
// This is to optimize the ray intersection search when having many hittable objects.
func NewBoundingVolumeHierarchy(hl HittableList) Hittable {

	if len(hl.list) == 0 {
		panic("Cannot create a Bvh with empty list of objects")
	}

	return _createBvh(hl.list, 0, len(hl.list))
}

func _createBvh(list []Hittable, start, end int) *bvh {

	numObjects := end - start
	var left Hittable
	var right Hittable
	var bBox aabb

	if numObjects == 1 {
		left = list[start]
		right = list[start]
		bBox = left.BoundingBox()

	} else if numObjects == 2 {
		sortHittablesSliceByMostSpreadAxis(list, start, end)
		left = list[start]
		right = list[start+1]
		bBox = combineAabbs(left.BoundingBox(), right.BoundingBox())

	} else {
		sortHittablesSliceByMostSpreadAxis(list, start, end)
		mid := start + numObjects/2
		left = _createBvh(list, start, mid)
		right = _createBvh(list, mid, end)
		bBox = combineAabbs(left.BoundingBox(), right.BoundingBox())
	}

	return &bvh{left: &left, right: &right, bBox: &bBox}
}

func sortHittablesSliceByMostSpreadAxis(list []Hittable, start, end int) {
	slice := list[start:end]

	xSpread := boundingBoxSpread(slice, func(h Hittable) util.Interval { return h.BoundingBox().x })
	ySpread := boundingBoxSpread(slice, func(h Hittable) util.Interval { return h.BoundingBox().y })
	zSpread := boundingBoxSpread(slice, func(h Hittable) util.Interval { return h.BoundingBox().z })

	if xSpread >= ySpread && xSpread >= zSpread {
		sortHittablesByBoundingBox(slice, func(h Hittable) util.Interval { return h.BoundingBox().x })
	} else if ySpread >= xSpread && ySpread >= zSpread {
		sortHittablesByBoundingBox(slice, func(h Hittable) util.Interval { return h.BoundingBox().y })
	} else {
		sortHittablesByBoundingBox(slice, func(h Hittable) util.Interval { return h.BoundingBox().z })
	}
}

func boundingBoxSpread(list []Hittable, boundingIntervalFunc func(h Hittable) util.Interval) float64 {
	min := util.Infinity
	max := -util.Infinity
	for _, h := range list {
		min = math.Min(min, boundingIntervalFunc(h).Min)
		max = math.Max(max, boundingIntervalFunc(h).Min)
	}
	return max - min
}

func sortHittablesByBoundingBox(list []Hittable, boundingIntervalFunc func(h Hittable) util.Interval) {
	sort.Slice(list, func(i, j int) bool { return boundingIntervalFunc(list[i]).Min < boundingIntervalFunc(list[j]).Min })
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
	return *b.bBox
}

func (b *bvh) IsLight() bool {
	return false
}
