package trace

import (
	"fmt"
	"math"
)

type rotationY struct {
	object   hittable
	sinTheta float64
	cosTheta float64
	bBox     aabb
}

func createRotationY(
	object hittable,
	angle float64,
) rotationY {
	radians := degreesToRadians(angle)
	sinTheta := math.Sin(radians)
	cosTheta := math.Cos(radians)
	bBox := object.boundingBox()

	min := vec3{infinity, infinity, infinity}
	max := vec3{-infinity, -infinity, -infinity}

	for i := 0.; i < 2; i++ {
		for j := 0.; j < 2; j++ {
			for k := 0.; k < 2; k++ {
				x := i*bBox.x.max + (1-i)*bBox.x.min
				y := j*bBox.y.max + (1-j)*bBox.y.min
				z := k*bBox.z.max + (1-k)*bBox.z.min

				newx := cosTheta*x + sinTheta*z
				newz := -sinTheta*x + cosTheta*z

				tester := vec3{newx, y, newz}

				min.x = math.Min(min.x, tester.x)
				min.y = math.Min(min.y, tester.y)
				min.z = math.Min(min.z, tester.z)

				max.x = math.Max(max.x, tester.x)
				max.y = math.Max(max.y, tester.y)
				max.z = math.Max(max.z, tester.z)
			}
		}
	}

	boundingBox := createAabbFromPoints(min, max)

	return rotationY{
		object,
		sinTheta,
		cosTheta,
		boundingBox,
	}
}

func (ry rotationY) hit(r ray, rayLength interval) (bool, *hitRecord) {

	origin := r.origin
	direction := r.direction

	origin.x = ry.cosTheta*r.origin.x - ry.sinTheta*r.origin.z
	origin.z = ry.sinTheta*r.origin.x + ry.cosTheta*r.origin.z

	direction.x = ry.cosTheta*r.direction.x - ry.sinTheta*r.direction.z
	direction.z = ry.sinTheta*r.direction.x + ry.cosTheta*r.direction.z

	rotatedR := ray{origin, direction, r.time}

	hit, rec := ry.object.hit(rotatedR, rayLength)
	if !hit {
		return hit, rec
	}

	hitPoint := rec.hitPoint
	hitPoint.x = ry.cosTheta*rec.hitPoint.x + ry.sinTheta*rec.hitPoint.z
	hitPoint.z = -ry.sinTheta*rec.hitPoint.x + ry.cosTheta*rec.hitPoint.z

	normal := rec.normal
	normal.x = ry.cosTheta*rec.normal.x + ry.sinTheta*rec.normal.z
	normal.z = -ry.sinTheta*rec.normal.x + ry.cosTheta*rec.normal.z

	rec.hitPoint = hitPoint
	rec.normal = normal

	return hit, rec
}

func (r rotationY) boundingBox() aabb {
	return r.bBox
}

func (r rotationY) String() string {
	return fmt.Sprintf("%v", r.object)
}
