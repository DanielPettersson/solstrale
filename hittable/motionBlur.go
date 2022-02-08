package hittable

import (
	"fmt"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/material"
)

type motionBlur struct {
	blurredHittable Hittable
	blurDirection   geo.Vec3
	bBox            aabb
}

func NewMotionBlur(
	blurredHittable Hittable,
	blurDirection geo.Vec3,
) Hittable {

	boundingBox1 := blurredHittable.BoundingBox()
	boundingBox2 := blurredHittable.BoundingBox().add(blurDirection)
	boundingBox := combineAabbs(boundingBox1, boundingBox2)

	return motionBlur{
		blurredHittable,
		blurDirection,
		boundingBox,
	}
}

func (m motionBlur) Hit(r geo.Ray, rayLength util.Interval) (bool, *material.HitRecord) {

	offset := m.blurDirection.MulS(r.Time)

	offsetRay := geo.Ray{
		Origin:    r.Origin.Sub(offset),
		Direction: r.Direction,
		Time:      r.Time,
	}

	hit, record := m.blurredHittable.Hit(offsetRay, rayLength)
	if record != nil {
		record.HitPoint = record.HitPoint.Add(offset)
	}

	return hit, record
}

func (m motionBlur) BoundingBox() aabb {
	return m.bBox
}

func (m motionBlur) String() string {
	return fmt.Sprintf("%v", m.blurredHittable)
}
