package hittable

import (
	"math"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/material"
	"github.com/DanielPettersson/solstrale/random"
)

type sphere struct {
	center geo.Vec3
	radius float64
	mat    material.Material
	bBox   aabb
}

// NewSphere creates a new sphere shaped hittable object
func NewSphere(
	center geo.Vec3,
	radius float64,
	mat material.Material,
) Hittable {

	rVec := geo.Vec3{X: radius, Y: radius, Z: radius}
	boundingBox := createAabbFromPoints(center.Sub(rVec), center.Add(rVec))

	return sphere{
		center,
		radius,
		mat,
		boundingBox,
	}
}

func (s sphere) Hit(r geo.Ray, rayLength util.Interval) (bool, *material.HitRecord) {

	oc := r.Origin.Sub(s.center)
	a := r.Direction.LengthSquared()
	halfB := oc.Dot(r.Direction)
	c := oc.LengthSquared() - s.radius*s.radius

	discriminant := halfB*halfB - a*c
	if discriminant < 0 {
		return false, nil
	}
	sqrtd := math.Sqrt(discriminant)

	root := (-halfB - sqrtd) / a
	if !rayLength.Contains(root) {
		root = (-halfB + sqrtd) / a
		if !rayLength.Contains(root) {
			return false, nil
		}
	}

	hitPoint := r.At(root)
	normal := hitPoint.Sub(s.center).DivS(s.radius)
	u, v := calculateSphereUv(normal)

	frontFace := r.Direction.Dot(normal) < 0
	if !frontFace {
		normal = normal.Neg()
	}
	rec := material.HitRecord{
		HitPoint:  hitPoint,
		Normal:    normal,
		Material:  s.mat,
		RayLength: root,
		U:         u,
		V:         v,
		FrontFace: frontFace,
	}

	return true, &rec

}

func calculateSphereUv(pointOnSphere geo.Vec3) (float64, float64) {
	theta := math.Acos(-pointOnSphere.Y)
	phi := math.Atan2(-pointOnSphere.Z, pointOnSphere.X) + math.Pi
	u := phi / (2 * math.Pi)
	v := theta / math.Pi
	return u, v
}

func randomToSphere(radius, distanceSquared float64) geo.Vec3 {
	r1 := random.RandomNormalFloat()
	r2 := random.RandomNormalFloat()
	z := 1 + r2*(math.Sqrt(1-radius*radius/distanceSquared)-1)

	phi := 2 * math.Pi * r1
	zz := math.Sqrt(1 - z*z)
	x := math.Cos(phi) * zz
	y := math.Sin(phi) * zz

	return geo.NewVec3(x, y, z)
}

func (s sphere) BoundingBox() aabb {
	return s.bBox
}

func (s sphere) PdfValue(origin, direction geo.Vec3) float64 {
	ray := geo.Ray{
		Origin:    origin,
		Direction: direction,
	}

	hit, _ := s.Hit(ray, util.Interval{Min: 0.001, Max: util.Infinity})

	if !hit {
		return 0
	}

	cosThetaMax := math.Sqrt(1 - s.radius*s.radius/s.center.Sub(origin).LengthSquared())
	solidAngle := 2 * math.Pi * (1 - cosThetaMax)

	return 1 / solidAngle
}

func (s sphere) RandomDirection(origin geo.Vec3) geo.Vec3 {
	direction := s.center.Sub(origin)
	uvw := geo.BuildOnbFromVec3(direction)
	return uvw.Local(randomToSphere(s.radius, direction.LengthSquared()))
}
