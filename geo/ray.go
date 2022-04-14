// Package geo provides basic geometric constructs
package geo

// Ray defines a ray of light used by the ray tracer
type Ray struct {
	Origin            Vec3
	Direction         Vec3
	DirectionInverted Vec3
	Time              float64
}

func NewRay(origin, direction Vec3, time float64) Ray {
	dir := direction.Unit()
	dirInv := NewVec3(1/dir.X, 1/dir.Y, 1/dir.Z)

	return Ray{
		Origin:            origin,
		Direction:         dir,
		DirectionInverted: dirInv,
		Time:              time,
	}
}

// At returns the position at a given length of the ray
func (r Ray) At(distance float64) Vec3 {
	return r.Origin.Add(r.Direction.MulS(distance))
}
