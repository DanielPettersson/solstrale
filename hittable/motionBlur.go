package hittable

import (
	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/material"
)

type motionBlur struct {
	blurredHittable Hittable
	blurDirection   geo.Vec3
	bBox            aabb
}

// NewMotionBlur creates a new hittable object that adds linear interpolated translation to
// its hittable based on the time of the ray. This gives the appearence of the object moving.
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
