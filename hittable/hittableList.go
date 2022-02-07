package hittable

import (
	"fmt"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/material"
)

type hittableList struct {
	list []Hittable
	bBox aabb
}

func NewHittableList() hittableList {
	return hittableList{
		[]Hittable{},
		aabb{util.EmptyInterval, util.EmptyInterval, util.EmptyInterval},
	}
}

func (hl *hittableList) Clear() {
	hl = nil
}

func (hl *hittableList) Add(h Hittable) {
	hl.list = append(hl.list, h)
	hl.bBox = combineAabbs(hl.bBox, h.BoundingBox())
}

func (hl *hittableList) Hit(r geo.Ray, rayLength util.Interval) (bool, *material.HitRecord) {
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

func (hl *hittableList) BoundingBox() aabb {
	return hl.bBox
}

func (hl hittableList) String() string {
	return fmt.Sprintf("HittableList(%v)", hl.list)
}
