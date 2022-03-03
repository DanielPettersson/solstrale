package hittable

import (
	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/material"
)

type translation struct {
	object Hittable
	offset geo.Vec3
	bBox   aabb
}

// NewTranslation creates a hittable object that translates the given hittable by the givn offset vector
func NewTranslation(
	object Hittable,
	offset geo.Vec3,
) Hittable {

	boundingBox := object.BoundingBox().add(offset)

	return translation{
		object,
		offset,
		boundingBox,
	}
}

func (t translation) Hit(r geo.Ray, rayLength util.Interval) (bool, *material.HitRecord) {

	offsetRay := geo.Ray{
		Origin:    r.Origin.Sub(t.offset),
		Direction: r.Direction,
		Time:      r.Time,
	}

	hit, record := t.object.Hit(offsetRay, rayLength)
	if record != nil {
		record.HitPoint = record.HitPoint.Add(t.offset)
	}

	return hit, record
}

func (t translation) BoundingBox() aabb {
	return t.bBox
}

func (t translation) PdfValue(o, v geo.Vec3) float64 {
	return 0.0
}

func (t translation) Random(o geo.Vec3) geo.Vec3 {
	return geo.NewVec3(1, 0, 0)
}
