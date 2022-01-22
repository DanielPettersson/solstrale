package trace

import "fmt"

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

func (m motionBlur) hit(r ray, rayLength interval) (bool, *hitRecord) {

	offset := m.blurDirection.mulS(r.time)

	offsetRay := ray{
		origin:    r.origin.sub(offset),
		direction: r.direction,
		time:      r.time,
	}

	hit, record := m.blurredHittable.hit(offsetRay, rayLength)
	if record != nil {
		record.hitPoint = record.hitPoint.add(offset)
	}

	return hit, record
}

func (m motionBlur) boundingBox() aabb {
	return m.bBox
}

func (m motionBlur) String() string {
	return fmt.Sprintf("%v", m.blurredHittable)
}
