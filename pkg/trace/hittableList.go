package trace

import "fmt"

type hittableList struct {
	list []hittable
	bBox aabb
}

func emptyHittableList() hittableList {
	return hittableList{
		[]hittable{},
		aabb{empty_interval, empty_interval, empty_interval},
	}
}

func (hl *hittableList) clear() {
	hl = nil
}

func (hl *hittableList) add(h hittable) {
	hl.list = append(hl.list, h)
	hl.bBox = combineAabbs(hl.bBox, h.boundingBox())
}

func (hl *hittableList) hit(r ray, rayLength interval) (bool, *hitRecord) {
	var closestHitRecord *hitRecord
	hitAnything := false
	closestSoFar := rayLength.max

	for _, h := range hl.list {
		hit, hitRecord := h.hit(r, interval{rayLength.min, closestSoFar})
		if hit {
			hitAnything = true
			closestSoFar = hitRecord.rayLength
			closestHitRecord = hitRecord
		}
	}

	return hitAnything, closestHitRecord
}

func (hl *hittableList) boundingBox() aabb {
	return hl.bBox
}

func (hl hittableList) String() string {
	return fmt.Sprintf("HittableList(%v)", hl.list)
}
