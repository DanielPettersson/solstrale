package trace

import (
	"bytes"
	"math"
)

var (
	white        vec3    = vec3{1, 1, 1}
	lightBlue    vec3    = vec3{0.5, 0.7, 1}
	sphereCenter vec3    = vec3{0, 0, -1}
	sphereRadius float64 = 0.5
)

func hitSphere(center vec3, radius float64, r ray) float64 {
	oc := r.orig.sub(center)
	a := r.dir.lengthSquared()
	halfB := oc.dot(r.dir)
	c := oc.lengthSquared() - radius*radius
	discriminant := halfB*halfB - a*c

	if discriminant < 0 {
		return -1
	} else {
		return (-halfB - math.Sqrt(discriminant)) / a
	}

}

func rayColor(r ray) vec3 {
	t := hitSphere(sphereCenter, sphereRadius, r)
	if t > 0 {
		n := r.at(t).sub(sphereCenter).unit()
		return vec3{n.x + 1, n.y + 1, n.z + 1}.mulS(0.5)
	}

	t = 0.5 * (r.dir.unit().y + 1)

	whiteness := white.mulS(1 - t)
	blueness := lightBlue.mulS(t)

	return whiteness.add(blueness)
}

func RayTrace(imageWidth int, imageHeight int, progress chan float32, byteBuffer *bytes.Buffer) {

	ret := make([]byte, imageWidth*imageHeight*4)

	aspectRatio := float64(imageWidth) / float64(imageHeight)

	viewPortHeight := 2.0
	viewPortWidth := aspectRatio * viewPortHeight
	focalLength := 1.0

	origin := vec3{0, 0, 0}
	horizontal := vec3{viewPortWidth, 0, 0}
	vertical := vec3{0, viewPortHeight, 0}
	lowerLeftCorner := origin.sub(horizontal.divS(2))
	lowerLeftCorner = lowerLeftCorner.sub(vertical.divS(2))
	lowerLeftCorner = lowerLeftCorner.sub(vec3{0, 0, focalLength})

	for y := 0; y < imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			i := (((imageHeight-1)-y)*imageWidth + x) * 4

			u := float64(x) / float64(imageWidth-1)
			v := float64(y) / float64(imageHeight-1)

			rDir := lowerLeftCorner.add(horizontal.mulS(u))
			rDir = rDir.add(vertical.mulS(v))
			rDir = rDir.sub(origin)
			r := ray{origin, rDir}
			col := rayColor(r).toRgba()

			ret[i] = col.r
			ret[i+1] = col.g
			ret[i+2] = col.b
			ret[i+3] = col.a
		}
		progress <- float32(y) / float32(imageHeight)
	}

	byteBuffer.Write(ret)
	close(progress)
}
