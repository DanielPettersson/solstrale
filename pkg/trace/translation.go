package trace

import "fmt"

type translation struct {
	object hittable
	offset vec3
	bBox   aabb
}

func createTranslation(
	object hittable,
	offset vec3,
) translation {

	boundingBox := object.boundingBox().add(offset)

	return translation{
		object,
		offset,
		boundingBox,
	}
}

func (t translation) hit(r ray, rayLength interval) (bool, *hitRecord) {

	offsetRay := ray{
		origin:    r.origin.sub(t.offset),
		direction: r.direction,
		time:      r.time,
	}

	hit, record := t.object.hit(offsetRay, rayLength)
	if record != nil {
		record.hitPoint = record.hitPoint.add(t.offset)
	}

	return hit, record
}

func (t translation) boundingBox() aabb {
	return t.bBox
}

func (t translation) String() string {
	return fmt.Sprintf("%v", t.object)
}
