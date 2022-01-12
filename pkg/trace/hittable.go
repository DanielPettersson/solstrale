package trace

type hitRecord struct {
	p         vec3
	normal    vec3
	t         float64
	frontFace bool
}

func createHitRecord(r ray, p vec3, normal vec3, t float64) hitRecord {
	frontFace := r.dir.dot(normal) < 0
	n := normal
	if !frontFace {
		n = normal.neg()
	}
	return hitRecord{
		p:         p,
		normal:    n,
		t:         t,
		frontFace: frontFace,
	}
}

type hittable interface {
	hit(r ray, rayT interval) (bool, *hitRecord)
}
