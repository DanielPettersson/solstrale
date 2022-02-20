package hittable

import (
	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/material"
	"github.com/DanielPettersson/solstrale/random"
)

// HittableList is a special type of hittable that is a container
// for a list of other hittable objects. Used to be able to have many
// objects in a scene
type HittableList struct {
	list []Hittable
	bBox aabb
}

// NewHittableList creates new empty HittableList
func NewHittableList() HittableList {
	return HittableList{
		[]Hittable{},
		aabb{util.EmptyInterval, util.EmptyInterval, util.EmptyInterval},
	}
}

// Add adds a new hittable object to this HittableList
func (hl *HittableList) Add(h Hittable) {
	hl.list = append(hl.list, h)
	hl.bBox = combineAabbs(hl.bBox, h.BoundingBox())
}

// Hit Checks if the given ray hits any object in this list.
// And if so, returns the properties of that ray hit
func (hl *HittableList) Hit(r geo.Ray, rayLength util.Interval, rand random.Random) (bool, *material.HitRecord) {
	var closestHitRecord *material.HitRecord
	hitAnything := false
	closestSoFar := rayLength.Max
	closestInterval := util.Interval{Min: rayLength.Min, Max: closestSoFar}

	for _, h := range hl.list {
		hit, hitRecord := h.Hit(r, closestInterval, rand)
		if hit {
			hitAnything = true
			closestSoFar = hitRecord.RayLength
			closestHitRecord = hitRecord
			closestInterval = util.Interval{Min: rayLength.Min, Max: closestSoFar}
		}
	}

	return hitAnything, closestHitRecord
}

// BoundingBox returns the bounding box that encapsulates all hittables in the list
func (hl *HittableList) BoundingBox() aabb {
	return hl.bBox
}
