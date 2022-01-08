package trace

type ray struct {
	orig vec3
	dir  vec3
}

func (r ray) at(t float64) vec3 {
	return r.orig.add(r.dir.mulS(t))
}
