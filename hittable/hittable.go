package hittable

import (
	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/material"
)

type Hittable interface {
	Hit(r geo.Ray, rayLength util.Interval) (bool, *material.HitRecord)
	BoundingBox() aabb
}
