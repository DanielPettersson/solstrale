package trace

import "math"

type texture interface {
	color(u, v float64, point vec3) vec3
}

type solidColor struct {
	colorValue vec3
}

func (sc solidColor) color(u, v float64, point vec3) vec3 {
	return sc.colorValue
}

type checkerTexture struct {
	scale float64
	even  texture
	odd   texture
}

func (ct checkerTexture) color(u, v float64, point vec3) vec3 {
	invScale := 1 / ct.scale
	xInt := math.Floor(point.x * invScale)
	yInt := math.Floor(point.y * invScale)
	zInt := math.Floor(point.z * invScale)

	if int(xInt+yInt+zInt)%2 == 0 {
		return ct.even.color(u, v, point)
	} else {
		return ct.odd.color(u, v, point)
	}

}
