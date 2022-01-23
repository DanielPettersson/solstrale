package trace

import "math"

type texture interface {
	color(rec hitRecord) vec3
}

type solidColor struct {
	colorValue vec3
}

func (sc solidColor) color(rec hitRecord) vec3 {
	return sc.colorValue
}

type checkerTexture struct {
	scale float64
	even  texture
	odd   texture
}

func (ct checkerTexture) color(rec hitRecord) vec3 {
	invScale := 1 / ct.scale
	uInt := math.Floor(rec.u * invScale)
	vInt := math.Floor(rec.v * invScale)

	if int(uInt+vInt)%2 == 0 {
		return ct.even.color(rec)
	} else {
		return ct.odd.color(rec)
	}

}

type imageTexture struct {
	bytes  []byte
	width  int
	height int
}

func (it imageTexture) color(rec hitRecord) vec3 {
	u := interval{0, 1}.clamp(rec.u)
	v := 1 - interval{0, 1}.clamp(rec.v)

	x := int(u * float64(it.width))
	y := int(v * float64(it.height))
	i := (y*it.width + x) * 4

	return rgbaColor{
		r: it.bytes[i],
		g: it.bytes[i+1],
		b: it.bytes[i+2],
	}.toVec3()
}
