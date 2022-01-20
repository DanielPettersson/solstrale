package trace

import "math"

type sphere struct {
	center vec3
	radius float64
	mat    material
}

func (s sphere) hit(r ray, rayT interval) (bool, *hitRecord) {

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
	if !rayT.contains(root) {
		root = (-halfB + sqrtd) / a
		if !rayT.contains(root) {
			return false, nil
		}
	}

	p := r.at(root)
	normal := p.sub(s.center).divS(s.radius)
	hitRecord := createHitRecord(r, p, normal, root, s.mat)
	return true, &hitRecord

}
