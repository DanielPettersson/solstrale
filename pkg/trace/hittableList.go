package trace

type hittableList []hittable

func (hl *hittableList) clear() {
	hl = nil
}

func (hl *hittableList) add(h hittable) {
	*hl = append(*hl, h)
}

func (hl *hittableList) hit(r ray, rayTmin float64, rayTmax float64) (bool, *hitRecord) {
	var closestHitRecord *hitRecord
	hitAnything := false
	closestSoFar := rayTmax

	for _, h := range *hl {
		hit, hitRecord := h.hit(r, rayTmin, closestSoFar)
		if hit {
			hitAnything = true
			closestSoFar = hitRecord.t
			closestHitRecord = hitRecord
		}
	}

	return hitAnything, closestHitRecord
}
