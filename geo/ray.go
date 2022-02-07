package geo

type Ray struct {
	Origin    Vec3
	Direction Vec3
	Time      float64
}

func (r Ray) At(distance float64) Vec3 {
	return r.Origin.Add(r.Direction.MulS(distance))
}
