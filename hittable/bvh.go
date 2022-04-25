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
	left          *bvh
	right         *bvh
	leftTriangle  Triangle
	rightTriangle Triangle
	bBox          aabb
}

// NewBoundingVolumeHierarchy creates a new hittable object from the given hittable list
// The bounding Volume Hierarchy sorts the hittables in a binary tree
// where each node has a bounding box.
// This is to optimize the ray intersection search when having many hittable objects.
func NewBoundingVolumeHierarchy(list []Triangle) Hittable {

	if len(list) == 0 {
		panic("Cannot create a Bvh with empty list of objects")
	}

	bvhChan := make(chan *bvh)
	go createBvhAsync(list, 0, len(list), bvhChan)
	return <-bvhChan
}

func createBvh(list []Triangle, start, end int) *bvh {
	numObjects := end - start
	var left *bvh
	var right *bvh
	var leftTriangle Triangle
	var rightTriangle Triangle
	var bBox aabb

	if numObjects == 1 {
		leftTriangle = list[start]
		rightTriangle = list[start]
		bBox = leftTriangle.BoundingBox()

	} else if numObjects == 2 {
		leftTriangle = list[start]
		rightTriangle = list[start+1]
		bBox = combineAabbs(leftTriangle.BoundingBox(), rightTriangle.BoundingBox())

	} else {
		mid := sortHittablesSliceByMostSpreadAxis(list, start, end)
		left = createBvh(list, start, mid)
		right = createBvh(list, mid, end)
		bBox = combineAabbs(left.BoundingBox(), right.BoundingBox())
	}

	return &bvh{left: left, right: right, leftTriangle: leftTriangle, rightTriangle: rightTriangle, bBox: bBox}
}

func createBvhAsync(list []Triangle, start, end int, bvhChan chan<- *bvh) {
	numObjects := end - start

	if numObjects < asyncCountThreshold {
		bvhChan <- createBvh(list, start, end)
	} else {
		mid := sortHittablesSliceByMostSpreadAxis(list, start, end)

		leftChan := make(chan *bvh)
		rightChan := make(chan *bvh)
		go createBvhAsync(list, start, mid, leftChan)
		go createBvhAsync(list, mid, end, rightChan)

		left := <-leftChan
		right := <-rightChan
		bBox := combineAabbs(left.BoundingBox(), right.BoundingBox())
		bvhChan <- &bvh{left: left, right: right, bBox: bBox}
	}

}

func sortHittablesSliceByMostSpreadAxis(list []Triangle, start, end int) int {
	slice := list[start:end]

	xSpread, xCenter := boundingBoxSpread(slice, 0)
	ySpread, yCenter := boundingBoxSpread(slice, 1)
	zSpread, zCenter := boundingBoxSpread(slice, 2)

	var center int
	if xSpread >= ySpread && xSpread >= zSpread {
		center = SortTrianglesByCenter(slice, xCenter, 0)
	} else if ySpread >= xSpread && ySpread >= zSpread {
		center = SortTrianglesByCenter(slice, yCenter, 1)
	} else {
		center = SortTrianglesByCenter(slice, zCenter, 2)
	}

	center += start

	// Could not split with objects on both sides. Just split up the middle index
	if center == start || center == end {
		center = start + (end-start)/2
	}
	return center
}

func boundingBoxSpread(list []Triangle, axis int) (float64, float64) {
	min := util.Infinity
	max := -util.Infinity
	listLen := len(list)
	for i := 0; i < listLen; i++ {
		c := list[i].Center(axis)
		min = math.Min(min, c)
		max = math.Max(max, c)
	}
	return max - min, (min + max) * .5
}

func SortTrianglesByCenter(list []Triangle, center float64, axis int) int {

	i := 0
	j := len(list) - 1

	for i <= j {
		if list[i].Center(axis) < center {
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

func (b *bvh) HitLeft(r geo.Ray, rayLength util.Interval) (bool, *material.HitRecord) {
	if b.left != nil {
		return (*b.left).Hit(r, rayLength)
	} else {
		return b.leftTriangle.Hit(r, rayLength)
	}
}

func (b *bvh) HitRight(r geo.Ray, rayLength util.Interval) (bool, *material.HitRecord) {
	if b.right != nil {
		return (*b.right).Hit(r, rayLength)
	} else {
		return b.rightTriangle.Hit(r, rayLength)
	}
}

func (b *bvh) Hit(r geo.Ray, rayLength util.Interval) (bool, *material.HitRecord) {
	if !b.bBox.hit(r, rayLength) {
		return false, nil
	}

	hitLeft, rec := b.HitLeft(r, rayLength)
	if hitLeft {
		rayLength = util.Interval{Min: rayLength.Min, Max: rec.RayLength}
	}

	hitRight, recRight := b.HitRight(r, rayLength)
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
