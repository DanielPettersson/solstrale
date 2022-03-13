// Package pdf provides probability density functions
package pdf

import (
	"math"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/random"
)

// Pdf is the common interface for the probability density functions
type Pdf interface {
	Value(direction geo.Vec3) float64
	Generate() geo.Vec3
}

// CosinePdf is a probability density functions with a cosine distribution
type CosinePdf struct {
	uvw geo.Onb
}

// NewCosinePdf creates a new instance of a CosinePdf
func NewCosinePdf(w geo.Vec3) Pdf {
	return CosinePdf{
		uvw: geo.BuildOnbFromVec3(w),
	}
}

// Value returns the pdf value for a given vector for the CosinePdf
func (p CosinePdf) Value(direction geo.Vec3) float64 {
	cosineTheta := direction.Unit().Dot(p.uvw.W)
	return math.Max(0, cosineTheta/math.Pi)
}

// Generates random direction for the CosinePdf shape
func (p CosinePdf) Generate() geo.Vec3 {
	return p.uvw.Local(geo.RandomCosineDirection())
}

// SpherePdf is a probability density functions with a sphere distribution
type SpherePdf struct{}

// NewSpherePdf creates a new instance of a SpherePdf
func NewSpherePdf() Pdf {
	return SpherePdf{}
}

// Value returns the pdf value for a given vector for the SpherePdf
func (p SpherePdf) Value(direction geo.Vec3) float64 {
	return 1 / (4 * math.Pi)
}

// Generates random direction for the SpherePdf shape
func (p SpherePdf) Generate() geo.Vec3 {
	return geo.RandomUnitVector()
}

// MixturePdf is for generating a mixture of two different probability density functions
type MixturePdf struct {
	p0 *Pdf
	p1 *Pdf
}

// NewMixturePdf creates a new instance of a MixturePdf
func NewMixturePdf(p0, p1 *Pdf) Pdf {
	return MixturePdf{
		p0: p0,
		p1: p1,
	}
}

// Value returns the pdf value for a given vector for the MixturePdf.
// Which is the average of the two base pdfs
func (p MixturePdf) Value(direction geo.Vec3) float64 {
	return .5*(*p.p0).Value(direction) + .5*(*p.p1).Value(direction)
}

// Generates random direction for the MixturePdf shape.
// Which is randomly chosen between the two base pdfs.
func (p MixturePdf) Generate() geo.Vec3 {
	if random.RandomNormalFloat() < .5 {
		return (*p.p0).Generate()
	}
	return (*p.p1).Generate()
}
