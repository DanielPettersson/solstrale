package trace

type hitRecord struct {
	hitPoint  vec3
	normal    vec3
	material  material
	rayLength float64
	u         float64
	v         float64
	frontFace bool
}

func createHitRecord(r ray, p vec3, normal vec3, t float64, mat material) hitRecord {
	frontFace := r.direction.dot(normal) < 0
	n := normal
	if !frontFace {
		n = normal.neg()
	}
	return hitRecord{
		hitPoint:  p,
		normal:    n,
		material:  mat,
		rayLength: t,
		frontFace: frontFace,
	}
}

type hittable interface {
	hit(r ray, rayLength interval) (bool, *hitRecord)
	boundingBox() aabb
}
