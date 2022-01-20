package trace

type ray struct {
	origin    vec3
	direction vec3
	time      float64
}

func (r ray) at(distance float64) vec3 {
	return r.origin.add(r.direction.mulS(distance))
}
