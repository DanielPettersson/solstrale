package trace

import (
	"fmt"
	"math"
	"sort"
)

// Bounding Volume Hierarchy
type bvh struct {
	left  *hittable
	right *hittable
	bBox  aabb
}

func createBvh(hl hittableList) bvh {
	return _createBvh(hl.list, 0, len(hl.list))
}

func _createBvh(list []hittable, start, end int) bvh {

	numObjects := end - start
	var left hittable
	var right hittable
	var bBox aabb

	if numObjects == 1 {
		left = list[start]
		right = list[start]
		bBox = left.boundingBox()

	} else if numObjects == 2 {
		sortHittablesSliceByMostSpreadAxis(list, start, end)
		left = list[start]
		right = list[start+1]
		bBox = combineAabbs(left.boundingBox(), right.boundingBox())

	} else {
		sortHittablesSliceByMostSpreadAxis(list, start, end)
		mid := start + numObjects/2
		left = _createBvh(list, start, mid)
		right = _createBvh(list, mid, end)
		bBox = combineAabbs(left.boundingBox(), right.boundingBox())
	}

	return bvh{&left, &right, bBox}
}

func sortHittablesSliceByMostSpreadAxis(list []hittable, start, end int) {
	slice := list[start:end]

	xSpread := boundingBoxSpread(slice, func(h hittable) interval { return h.boundingBox().x })
	ySpread := boundingBoxSpread(slice, func(h hittable) interval { return h.boundingBox().y })
	zSpread := boundingBoxSpread(slice, func(h hittable) interval { return h.boundingBox().z })

	if xSpread > ySpread && xSpread > zSpread {
		sortHittablesByBoundingBox(slice, func(h hittable) interval { return h.boundingBox().x })
	} else if ySpread > xSpread && ySpread > zSpread {
		sortHittablesByBoundingBox(slice, func(h hittable) interval { return h.boundingBox().y })
	} else {
		sortHittablesByBoundingBox(slice, func(h hittable) interval { return h.boundingBox().z })
	}
}

func boundingBoxSpread(list []hittable, boundingIntervalFunc func(h hittable) interval) float64 {
	min := infinity
	max := -infinity
	for _, h := range list {
		min = math.Min(min, boundingIntervalFunc(h).min)
		max = math.Max(max, boundingIntervalFunc(h).min)
	}
	return max - min
}

func sortHittablesByBoundingBox(list []hittable, boundingIntervalFunc func(h hittable) interval) {
	sort.Slice(list, func(i, j int) bool { return boundingIntervalFunc(list[i]).min < boundingIntervalFunc(list[j]).min })
}

func (b bvh) hit(r ray, rayT interval) (bool, *hitRecord) {
	if !b.bBox.hit(r, rayT) {
		return false, nil
	}

	hitLeft, rec := (*b.left).hit(r, rayT)
	if hitLeft {
		rayT = interval{rayT.min, rec.t}
	}

	hitRight, recRight := (*b.right).hit(r, rayT)
	if hitRight {
		rec = recRight
	}

	return hitLeft || hitRight, rec
}

func (b bvh) boundingBox() aabb {
	return b.bBox
}

func (b bvh) String() string {
	return fmt.Sprintf("BoundingVolumeHierarchy(left: %v, right: %v)", *b.left, *b.right)
}
