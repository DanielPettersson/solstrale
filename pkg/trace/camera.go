package trace

import "math"

type camera struct {
	origin          vec3
	lowerLeftCorner vec3
	horizontal      vec3
	vertical        vec3
}

func createCamera(
	spec TraceSpecification,
	verticalFovDegrees float64,
	lookFrom vec3,
	lookAt vec3,
	vup vec3,
) camera {
	aspectRatio := float64(spec.ImageWidth) / float64(spec.ImageHeight)
	theta := degreesToRadians(verticalFovDegrees)
	h := math.Tan(theta / 2)
	viewPortHeight := 2.0 * h
	viewPortWidth := aspectRatio * viewPortHeight

	w := lookFrom.sub(lookAt).unit()
	u := vup.cross(w).unit()
	v := w.cross(u)

	origin := lookFrom
	horizontal := u.mulS(viewPortWidth)
	vertical := v.mulS(viewPortHeight)
	lowerLeftCorner := origin.sub(horizontal.divS(2)).sub(vertical.divS(2)).sub(w)

	return camera{
		origin,
		lowerLeftCorner,
		horizontal,
		vertical,
	}
}

func (c camera) getRay(u float64, v float64) ray {
	rDir := c.lowerLeftCorner.add(c.horizontal.mulS(u))
	rDir = rDir.add(c.vertical.mulS(v))
	rDir = rDir.sub(c.origin)
	return ray{c.origin, rDir}
}
