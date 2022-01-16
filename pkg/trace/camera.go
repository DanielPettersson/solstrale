package trace

import "math"

type camera struct {
	origin          vec3
	lowerLeftCorner vec3
	horizontal      vec3
	vertical        vec3
	u               vec3
	v               vec3
	w               vec3
	lensRadius      float64
}

func createCamera(
	spec TraceSpecification,
	verticalFovDegrees float64,
	aperture float64,
	focusDistance float64,
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
	horizontal := u.mulS(viewPortWidth).mulS(focusDistance)
	vertical := v.mulS(viewPortHeight).mulS(focusDistance)
	lowerLeftCorner := origin.sub(horizontal.divS(2)).sub(vertical.divS(2)).sub(w.mulS(focusDistance))

	return camera{
		origin,
		lowerLeftCorner,
		horizontal,
		vertical,
		u,
		v,
		w,
		aperture / 2,
	}
}

func (c camera) getRay(u float64, v float64) ray {
	rd := randomInUnitDisc().mulS(c.lensRadius)
	offset := c.u.mulS(rd.x).add(c.v.mulS(rd.y))

	rDir := c.lowerLeftCorner.add(c.horizontal.mulS(u))
	rDir = rDir.add(c.vertical.mulS(v))
	rDir = rDir.sub(c.origin).sub(offset)
	return ray{
		c.origin.add(offset),
		rDir,
	}
}
