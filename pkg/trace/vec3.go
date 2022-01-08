package trace

import "math"

type vec3 struct {
	x float64
	y float64
	z float64
}

func (v vec3) neg() vec3 {
	return vec3{-v.x, -v.y, -v.z}
}

func (v vec3) add(w vec3) vec3 {
	return vec3{v.x + w.x, v.y + w.y, v.z + w.z}
}

func (v vec3) sub(w vec3) vec3 {
	return vec3{v.x - w.x, v.y - w.y, v.z - w.z}
}

func (v vec3) mul(w vec3) vec3 {
	return vec3{v.x * w.x, v.y * w.y, v.z * w.z}
}

func (v vec3) mulS(t float64) vec3 {
	return vec3{v.x * t, v.y * t, v.z * t}
}

func (v vec3) divS(t float64) vec3 {
	return vec3{v.x / t, v.y / t, v.z / t}
}

func (v vec3) dot(w vec3) float64 {
	return v.x*w.x + v.y*w.y + v.z + w.z
}

func (v vec3) cross(w vec3) vec3 {
	return vec3{
		w.y*v.z - w.z*v.y,
		w.z*v.x - w.x*v.z,
		w.x*v.y - w.y*v.x,
	}
}

func (v vec3) lengthSquared() float64 {
	return v.x*v.x + v.y*v.y + v.z*v.z
}

func (v vec3) length() float64 {
	return math.Sqrt(v.lengthSquared())
}

func (v vec3) unit() vec3 {
	return v.divS(v.length())
}

type rgbaColor struct {
	r byte
	g byte
	b byte
	a byte
}

func (c vec3) toRgba() rgbaColor {
	return rgbaColor{
		byte(c.x * 255.999),
		byte(c.y * 255.999),
		byte(c.z * 255.999),
		255,
	}
}
