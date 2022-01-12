package trace

type hittableList []hittable

func (hl *hittableList) clear() {
	hl = nil
}

func (hl *hittableList) add(h hittable) {
	*hl = append(*hl, h)
}

func (hl *hittableList) hit(r ray, rayT interval) (bool, *hitRecord) {
	var closestHitRecord *hitRecord
	hitAnything := false
	closestSoFar := rayT.max

	for _, h := range *hl {
		hit, hitRecord := h.hit(r, interval{rayT.min, closestSoFar})
		if hit {
			hitAnything = true
			closestSoFar = hitRecord.t
			closestHitRecord = hitRecord
		}
	}

	return hitAnything, closestHitRecord
}
