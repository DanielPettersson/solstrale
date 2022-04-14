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
		object: object,
		offset: offset,
		bBox:   boundingBox,
	}
}

func (t translation) Hit(r geo.Ray, rayLength util.Interval) (bool, *material.HitRecord) {

	offsetRay := geo.NewRay(
		r.Origin.Sub(t.offset),
		r.Direction,
		r.Time,
	)

	hit, record := t.object.Hit(offsetRay, rayLength)
	if record != nil {
		record.HitPoint = record.HitPoint.Add(t.offset)
	}

	return hit, record
}

func (t translation) BoundingBox() aabb {
	return t.bBox
}

func (t translation) PdfValue(origin, direction geo.Vec3) float64 {
	return t.object.PdfValue(origin, direction)
}

func (t translation) RandomDirection(origin geo.Vec3) geo.Vec3 {
	return t.object.RandomDirection(origin)
}

func (t translation) IsLight() bool {
	return t.object.IsLight()
}
