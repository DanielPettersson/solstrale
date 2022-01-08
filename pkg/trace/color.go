package trace

type color vec3

type rgbaColor struct {
	r byte
	g byte
	b byte
	a byte
}

func (c color) toRgba() rgbaColor {
	return rgbaColor{
		byte(c.x * 255.999),
		byte(c.y * 255.999),
		byte(c.z * 255.999),
		255,
	}
}
