package trace

type hitRecord struct {
	p         vec3
	normal    vec3
	mat       material
	t         float64
	frontFace bool
}

func createHitRecord(r ray, p vec3, normal vec3, t float64, mat material) hitRecord {
	frontFace := r.dir.dot(normal) < 0
	n := normal
	if !frontFace {
		n = normal.neg()
	}
	return hitRecord{
		p:         p,
		normal:    n,
		mat:       mat,
		t:         t,
		frontFace: frontFace,
	}
}

type hittable interface {
	hit(r ray, rayT interval) (bool, *hitRecord)
}
