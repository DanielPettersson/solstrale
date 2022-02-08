// Package geo provides basic geometric constructs
package geo

// Ray defines a ray of light used by the ray tracer
type Ray struct {
	Origin    Vec3
	Direction Vec3
	Time      float64
}

// At returns the position at a given length of the ray
func (r Ray) At(distance float64) Vec3 {
	return r.Origin.Add(r.Direction.MulS(distance))
}
