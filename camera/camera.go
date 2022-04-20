// Package camera provides a camera used by raytracer to shoot rays into the scene
package camera

import (
	"math"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/random"
)

// CameraConfig contains all needed parameters for constructing a camera
type CameraConfig struct {
	VerticalFovDegrees float64
	ApertureSize       float64
	FocusDistance      float64
	LookFrom           geo.Vec3
	LookAt             geo.Vec3
}

// Camera contains all data needed to describe a cameras position, field of view and
// where it is pointing
type Camera struct {
	origin          geo.Vec3
	lowerLeftCorner geo.Vec3
	horizontal      geo.Vec3
	vertical        geo.Vec3
	u               geo.Vec3
	v               geo.Vec3
	lensRadius      float64
}

// New creates a new camera from image dimensions and config
func New(
	imageWidth int,
	imageHeight int,
	c CameraConfig,
) Camera {
	aspectRatio := float64(imageWidth) / float64(imageHeight)
	theta := util.DegreesToRadians(c.VerticalFovDegrees)
	h := math.Tan(theta / 2)
	viewPortHeight := 2.0 * h
	viewPortWidth := aspectRatio * viewPortHeight

	w := c.LookFrom.Sub(c.LookAt).Unit()
	u := geo.NewVec3(0, 1, 0).Cross(w).Unit()
	v := w.Cross(u)

	origin := c.LookFrom
	horizontal := u.MulS(viewPortWidth).MulS(c.FocusDistance)
	vertical := v.MulS(viewPortHeight).MulS(c.FocusDistance)
	lowerLeftCorner := origin.Sub(horizontal.DivS(2)).Sub(vertical.DivS(2)).Sub(w.MulS(c.FocusDistance))

	return Camera{
		origin,
		lowerLeftCorner,
		horizontal,
		vertical,
		u,
		v,
		c.ApertureSize / 2,
	}
}

// GetRay is a function for generating a ray for a certain u/v for the raytraced image
func (c Camera) GetRay(u float64, v float64) geo.Ray {
	rd := geo.RandomInUnitDisc().MulS(c.lensRadius)
	offset := c.u.MulS(rd.X).Add(c.v.MulS(rd.Y))

	rDir := c.lowerLeftCorner.Add(c.horizontal.MulS(u))
	rDir = rDir.Add(c.vertical.MulS(v))
	rDir = rDir.Sub(c.origin).Sub(offset)
	return geo.NewRay(
		c.origin.Add(offset),
		rDir,
		random.RandomNormalFloat(),
	)
}
