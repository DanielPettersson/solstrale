package camera

import (
	"math"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/spec"
)

type Camera struct {
	origin          geo.Vec3
	lowerLeftCorner geo.Vec3
	horizontal      geo.Vec3
	vertical        geo.Vec3
	u               geo.Vec3
	v               geo.Vec3
	lensRadius      float64
}

func New(
	spec spec.TraceSpecification,
	verticalFovDegrees float64,
	aperture float64,
	focusDistance float64,
	lookFrom geo.Vec3,
	lookAt geo.Vec3,
	vup geo.Vec3,
) Camera {
	aspectRatio := float64(spec.ImageWidth) / float64(spec.ImageHeight)
	theta := util.DegreesToRadians(verticalFovDegrees)
	h := math.Tan(theta / 2)
	viewPortHeight := 2.0 * h
	viewPortWidth := aspectRatio * viewPortHeight

	w := lookFrom.Sub(lookAt).Unit()
	u := vup.Cross(w).Unit()
	v := w.Cross(u)

	origin := lookFrom
	horizontal := u.MulS(viewPortWidth).MulS(focusDistance)
	vertical := v.MulS(viewPortHeight).MulS(focusDistance)
	lowerLeftCorner := origin.Sub(horizontal.DivS(2)).Sub(vertical.DivS(2)).Sub(w.MulS(focusDistance))

	return Camera{
		origin,
		lowerLeftCorner,
		horizontal,
		vertical,
		u,
		v,
		aperture / 2,
	}
}

func (c Camera) GetRay(u float64, v float64) geo.Ray {
	rd := geo.RandomInUnitDisc().MulS(c.lensRadius)
	offset := c.u.MulS(rd.X).Add(c.v.MulS(rd.Y))

	rDir := c.lowerLeftCorner.Add(c.horizontal.MulS(u))
	rDir = rDir.Add(c.vertical.MulS(v))
	rDir = rDir.Sub(c.origin).Sub(offset)
	return geo.Ray{
		Origin:    c.origin.Add(offset),
		Direction: rDir,
		Time:      util.RandomNormalFloat(),
	}
}
