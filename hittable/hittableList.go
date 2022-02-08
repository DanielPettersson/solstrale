package hittable

import (
	"fmt"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/material"
)

type HittableList struct {
	list []Hittable
	bBox aabb
}

func NewHittableList() HittableList {
	return HittableList{
		[]Hittable{},
		aabb{util.EmptyInterval, util.EmptyInterval, util.EmptyInterval},
	}
}

func (hl *HittableList) Clear() {
	hl = nil
}

func (hl *HittableList) Add(h Hittable) {
	hl.list = append(hl.list, h)
	hl.bBox = combineAabbs(hl.bBox, h.BoundingBox())
}

func (hl *HittableList) Hit(r geo.Ray, rayLength util.Interval) (bool, *material.HitRecord) {
	var closestHitRecord *material.HitRecord
	hitAnything := false
	closestSoFar := rayLength.Max
	closestInterval := util.Interval{Min: rayLength.Min, Max: closestSoFar}

	for _, h := range hl.list {
		hit, hitRecord := h.Hit(r, closestInterval)
		if hit {
			hitAnything = true
			closestSoFar = hitRecord.RayLength
			closestHitRecord = hitRecord
			closestInterval = util.Interval{Min: rayLength.Min, Max: closestSoFar}
		}
	}

	return hitAnything, closestHitRecord
}

func (hl *HittableList) BoundingBox() aabb {
	return hl.bBox
}

func (hl HittableList) String() string {
	return fmt.Sprintf("HittableList(%v)", hl.list)
}
