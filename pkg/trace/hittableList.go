package trace

type hittableList struct {
	list *[]hittable
	bBox aabb
}

func emptyHittableList() hittableList {
	return hittableList{
		&[]hittable{},
		aabb{empty_interval, empty_interval, empty_interval},
	}
}

func (hl *hittableList) clear() {
	hl = nil
}

func (hl *hittableList) add(h hittable) {
	*hl.list = append(*hl.list, h)
	hl.bBox = combineAabbs(hl.bBox, h.boundingBox())
}

func (hl *hittableList) hit(r ray, rayT interval) (bool, *hitRecord) {
	var closestHitRecord *hitRecord
	hitAnything := false
	closestSoFar := rayT.max

	for _, h := range *hl.list {
		hit, hitRecord := h.hit(r, interval{rayT.min, closestSoFar})
		if hit {
			hitAnything = true
			closestSoFar = hitRecord.t
			closestHitRecord = hitRecord
		}
	}

	return hitAnything, closestHitRecord
}

func (hl *hittableList) boundingBox() aabb {
	return hl.bBox
}
