package trace

type material interface {
	scatter(rayIn ray, rec hitRecord) (bool, vec3, ray)
}

type lambertian struct {
	albedo vec3
}

func (m lambertian) scatter(rayIn ray, rec hitRecord) (bool, vec3, ray) {
	scatterDirection := rec.normal.add(randomUnitVector())
	if scatterDirection.nearZero() {
		scatterDirection = rec.normal
	}

	scatterRay := ray{rec.p, scatterDirection}
	return true, m.albedo, scatterRay
}

type metal struct {
	albedo vec3
	fuzz   float64
}

func (m metal) scatter(rayIn ray, rec hitRecord) (bool, vec3, ray) {
	reflected := rayIn.dir.unit().reflect(rec.normal)
	scatterRay := ray{rec.p, reflected.add(randomInUnitSphere().mulS(m.fuzz))}
	scatter := scatterRay.dir.dot(rec.normal) > 0

	return scatter, m.albedo, scatterRay
}
