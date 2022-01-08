package trace

import "math"

type vec3 struct {
	x float64
	y float64
	z float64
}

func neg(v vec3) vec3 {
	return vec3{-v.x, -v.y, -v.z}
}

func add(v vec3, w vec3) vec3 {
	return vec3{v.x + w.x, v.y + w.y, v.z + w.z}
}

func sub(v vec3, w vec3) vec3 {
	return vec3{v.x - w.x, v.y - w.y, v.z - w.z}
}

func mul(v vec3, w vec3) vec3 {
	return vec3{v.x * w.x, v.y * w.y, v.z * w.z}
}

func mulScalar(v vec3, t float64) vec3 {
	return vec3{v.x * t, v.y * t, v.z * t}
}

func divScalar(v vec3, t float64) vec3 {
	return vec3{v.x / t, v.y / t, v.z / t}
}

func dot(v vec3, w vec3) float64 {
	return v.x*w.x + v.y*w.y + v.z + w.z
}

func cross(v vec3, w vec3) vec3 {
	return vec3{
		w.y*v.z - w.z*v.y,
		w.z*v.x - w.x*v.z,
		w.x*v.y - w.y*v.x,
	}
}

func (v *vec3) addTo(w vec3) vec3 {
	v.x += w.x
	v.y += w.y
	v.z += w.z
	return *v
}

func (v *vec3) multiplyBy(m float64) vec3 {
	v.x *= m
	v.y *= m
	v.z *= m
	return *v
}

func (v *vec3) divideBy(d float64) vec3 {
	return v.multiplyBy(1 / d)
}

func (v *vec3) lengthSquared() float64 {
	return v.x*v.x + v.y*v.y + v.z*v.z
}

func (v *vec3) length() float64 {
	return math.Sqrt(v.lengthSquared())
}

func (v *vec3) unit() vec3 {
	return divScalar(*v, v.length())
}
