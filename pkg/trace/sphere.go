package trace

import (
	"fmt"
	"math"
)

type sphere struct {
	center vec3
	radius float64
	mat    material
	bBox   aabb
}

func createSphere(
	center vec3,
	radius float64,
	mat material,
) sphere {

	rVec := vec3{radius, radius, radius}
	boundingBox := createAabbFromPoints(center.sub(rVec), center.add(rVec))

	return sphere{
		center,
		radius,
		mat,
		boundingBox,
	}
}

func (s sphere) hit(r ray, rayLength interval) (bool, *hitRecord) {

	oc := r.origin.sub(s.center)
	a := r.direction.lengthSquared()
	halfB := oc.dot(r.direction)
	c := oc.lengthSquared() - s.radius*s.radius

	discriminant := halfB*halfB - a*c
	if discriminant < 0 {
		return false, nil
	}
	sqrtd := math.Sqrt(discriminant)

	root := (-halfB - sqrtd) / a
	if !rayLength.contains(root) {
		root = (-halfB + sqrtd) / a
		if !rayLength.contains(root) {
			return false, nil
		}
	}

	hitPoint := r.at(root)
	normal := hitPoint.sub(s.center).divS(s.radius)
	u, v := calculateSphereUv(normal)

	frontFace := r.direction.dot(normal) < 0
	if !frontFace {
		normal = normal.neg()
	}
	rec := hitRecord{
		hitPoint:  hitPoint,
		normal:    normal,
		material:  s.mat,
		rayLength: root,
		u:         u,
		v:         v,
		frontFace: frontFace,
	}

	return true, &rec

}

func calculateSphereUv(pointOnSphere vec3) (float64, float64) {
	theta := math.Acos(-pointOnSphere.y)
	phi := math.Atan2(-pointOnSphere.z, pointOnSphere.x) + math.Pi
	u := phi / (2 * math.Pi)
	v := theta / math.Pi
	return u, v
}

func (s sphere) boundingBox() aabb {
	return s.bBox
}

func (s sphere) String() string {
	return fmt.Sprintf("Sphere(%v, r:%f)", s.center, s.radius)
}
