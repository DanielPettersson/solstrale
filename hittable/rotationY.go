package hittable

import (
	"math"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/material"
)

type rotationY struct {
	object   Hittable
	sinTheta float64
	cosTheta float64
	bBox     aabb
}

// NewRotationY creates a hittable object that rotates the given hittable around the Y axis
func NewRotationY(
	object Hittable,
	angle float64,
) Hittable {
	radians := util.DegreesToRadians(angle)
	sinTheta := math.Sin(radians)
	cosTheta := math.Cos(radians)
	bBox := object.BoundingBox()

	min := geo.Vec3{X: util.Infinity, Y: util.Infinity, Z: util.Infinity}
	max := geo.Vec3{X: -util.Infinity, Y: -util.Infinity, Z: -util.Infinity}

	for i := 0.; i < 2; i++ {
		for j := 0.; j < 2; j++ {
			for k := 0.; k < 2; k++ {
				x := i*bBox.x.Max + (1-i)*bBox.x.Min
				y := j*bBox.y.Max + (1-j)*bBox.y.Min
				z := k*bBox.z.Max + (1-k)*bBox.z.Min

				newx := cosTheta*x + sinTheta*z
				newz := -sinTheta*x + cosTheta*z

				tester := geo.Vec3{X: newx, Y: y, Z: newz}

				min.X = math.Min(min.X, tester.X)
				min.Y = math.Min(min.Y, tester.Y)
				min.Z = math.Min(min.Z, tester.Z)

				max.X = math.Max(max.X, tester.X)
				max.Y = math.Max(max.Y, tester.Y)
				max.Z = math.Max(max.Z, tester.Z)
			}
		}
	}

	boundingBox := createAabbFromPoints(min, max)

	return rotationY{
		object:   object,
		sinTheta: sinTheta,
		cosTheta: cosTheta,
		bBox:     boundingBox,
	}
}

func (ry rotationY) Hit(r geo.Ray, rayLength util.Interval) (bool, *material.HitRecord) {

	origin := r.Origin
	direction := r.Direction

	origin.X = ry.cosTheta*r.Origin.X - ry.sinTheta*r.Origin.Z
	origin.Z = ry.sinTheta*r.Origin.X + ry.cosTheta*r.Origin.Z

	direction.X = ry.cosTheta*r.Direction.X - ry.sinTheta*r.Direction.Z
	direction.Z = ry.sinTheta*r.Direction.X + ry.cosTheta*r.Direction.Z

	rotatedR := geo.Ray{Origin: origin, Direction: direction, Time: r.Time}

	hit, rec := ry.object.Hit(rotatedR, rayLength)
	if !hit {
		return hit, rec
	}

	hitPoint := rec.HitPoint
	hitPoint.X = ry.cosTheta*rec.HitPoint.X + ry.sinTheta*rec.HitPoint.Z
	hitPoint.Z = -ry.sinTheta*rec.HitPoint.X + ry.cosTheta*rec.HitPoint.Z

	normal := rec.Normal
	normal.X = ry.cosTheta*rec.Normal.X + ry.sinTheta*rec.Normal.Z
	normal.Z = -ry.sinTheta*rec.Normal.X + ry.cosTheta*rec.Normal.Z

	rec.HitPoint = hitPoint
	rec.Normal = normal

	return hit, rec
}

func (ry rotationY) BoundingBox() aabb {
	return ry.bBox
}

func (ry rotationY) PdfValue(origin, direction geo.Vec3) float64 {
	return ry.object.PdfValue(origin, direction)
}

func (ry rotationY) RandomDirection(origin geo.Vec3) geo.Vec3 {
	return ry.object.RandomDirection(origin)
}

func (ry rotationY) IsLight() bool {
	return ry.object.IsLight()
}
