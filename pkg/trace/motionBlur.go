package trace

type motionBlur struct {
	blurredHittable hittable
	blurDirection   vec3
	bBox            aabb
}

func createMotionBlur(
	blurredHittable hittable,
	blurDirection vec3,
) motionBlur {

	boundingBox1 := blurredHittable.boundingBox()
	boundingBox2 := blurredHittable.boundingBox().add(blurDirection)
	boundingBox := combineAabbs(boundingBox1, boundingBox2)

	return motionBlur{
		blurredHittable,
		blurDirection,
		boundingBox,
	}
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

func (m motionBlur) boundingBox() aabb {
	return m.bBox
}
