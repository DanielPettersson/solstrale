package hittable

import (
	"math"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/material"
)

const (
	asyncCountThreshold int = 1000
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

	bvhChan := make(chan Hittable)
	go createBvhAsync(list, 0, len(list), bvhChan)
	return <-bvhChan
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
		left = createBvh(list, start, mid)
		right = createBvh(list, mid, end)
		bBox = combineAabbs(left.BoundingBox(), right.BoundingBox())
	}

	return &bvh{left: &left, right: &right, bBox: bBox}
}

func createBvhAsync(list []Hittable, start, end int, bvhChan chan<- Hittable) {
	numObjects := end - start

	if numObjects < asyncCountThreshold {
		bvhChan <- createBvh(list, start, end)
	} else {
		mid := sortHittablesSliceByMostSpreadAxis(list, start, end)

		leftChan := make(chan Hittable)
		rightChan := make(chan Hittable)
		go createBvhAsync(list, start, mid, leftChan)
		go createBvhAsync(list, mid, end, rightChan)

		left := <-leftChan
		right := <-rightChan
		bBox := combineAabbs(left.BoundingBox(), right.BoundingBox())
		bvhChan <- &bvh{left: &left, right: &right, bBox: bBox}
	}

}

func sortHittablesSliceByMostSpreadAxis(list []Hittable, start, end int) int {
	slice := list[start:end]

	xSpread, xCenter := boundingBoxSpread(slice, 0)
	ySpread, yCenter := boundingBoxSpread(slice, 1)
	zSpread, zCenter := boundingBoxSpread(slice, 2)

	var center int
	if xSpread >= ySpread && xSpread >= zSpread {
		center = SortHittablesByCenter(slice, xCenter, 0)
	} else if ySpread >= xSpread && ySpread >= zSpread {
		center = SortHittablesByCenter(slice, yCenter, 1)
	} else {
		center = SortHittablesByCenter(slice, zCenter, 2)
	}

	center += start

	// Could not split with objects on both sides. Just split up the middle index
	if center == start || center == end {
		center = start + (end-start)/2
	}
	return center
}

func boundingBoxSpread(list []Hittable, axis int) (float64, float64) {
	min := util.Infinity
	max := -util.Infinity
	listLen := len(list)
	for i := 0; i < listLen; i++ {
		c := list[i].Center().Axis(axis)
		min = math.Min(min, c)
		max = math.Max(max, c)
	}
	return max - min, (min + max) * .5
}

func SortHittablesByCenter(list []Hittable, center float64, axis int) int {

	i := 0
	j := len(list) - 1

	for i <= j {
		if list[i].Center().Axis(axis) < center {
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
