package trace

type motionBlur struct {
	blurredHittable hittable
	blurDirection   vec3
}

func (m motionBlur) hit(r ray, rayT interval) (bool, *hitRecord) {

	offset := m.blurDirection.mulS(r.time)

	offsetRay := ray{
		origin:    r.origin.sub(offset),
		direction: r.direction,
		time:      r.time,
	}

	hit, record := m.blurredHittable.hit(offsetRay, rayT)
	if record != nil {
		record.p = record.p.add(offset)
	}

	return hit, record
}
