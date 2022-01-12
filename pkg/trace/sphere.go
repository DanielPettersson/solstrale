package trace

import "math"

type sphere struct {
	center vec3
	radius float64
}

func (s sphere) hit(r ray, rayT interval) (bool, *hitRecord) {

	oc := r.orig.sub(s.center)
	a := r.dir.lengthSquared()
	halfB := oc.dot(r.dir)
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
	hitRecord := createHitRecord(r, p, normal, root)
	return true, &hitRecord

}
