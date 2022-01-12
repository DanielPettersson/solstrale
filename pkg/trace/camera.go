package trace

type camera struct {
	origin          vec3
	lowerLeftCorner vec3
	horizontal      vec3
	vertical        vec3
}

func createCamera(imageWidth int, imageHeight int) camera {
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

	return camera{
		origin:          origin,
		lowerLeftCorner: lowerLeftCorner,
		horizontal:      horizontal,
		vertical:        vertical,
	}
}

func (c camera) getRay(u float64, v float64) ray {
	rDir := c.lowerLeftCorner.add(c.horizontal.mulS(u))
	rDir = rDir.add(c.vertical.mulS(v))
	rDir = rDir.sub(c.origin)
	return ray{c.origin, rDir}
}
