package pdf

import (
	"math"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/random"
)

type Pdf interface {
	Value(direction geo.Vec3) float64
	Generate() geo.Vec3
}

type CosinePdf struct {
	uvw geo.Onb
}

func NewCosinePdf(w geo.Vec3) Pdf {
	return CosinePdf{
		uvw: geo.BuildOnbFromVec3(w),
	}
}

func (p CosinePdf) Value(direction geo.Vec3) float64 {
	cosineTheta := direction.Unit().Dot(p.uvw.W)
	return math.Max(0, cosineTheta/math.Pi)
}

func (p CosinePdf) Generate() geo.Vec3 {
	return p.uvw.LocalV(geo.RandomCosineDirection())
}

type SpherePdf struct{}

func NewSpherePdf() Pdf {
	return SpherePdf{}
}

func (p SpherePdf) Value(direction geo.Vec3) float64 {
	return 1 / (4 * math.Pi)
}

func (p SpherePdf) Generate() geo.Vec3 {
	return geo.RandomUnitVector()
}

type MixturePdf struct {
	p0 *Pdf
	p1 *Pdf
}

func NewMixturePdf(p0, p1 *Pdf) Pdf {
	return MixturePdf{
		p0: p0,
		p1: p1,
	}
}

func (p MixturePdf) Value(direction geo.Vec3) float64 {
	return .5*(*p.p0).Value(direction) + .5*(*p.p1).Value(direction)
}

func (p MixturePdf) Generate() geo.Vec3 {
	if random.RandomNormalFloat() < .5 {
		return (*p.p0).Generate()
	}
	return (*p.p1).Generate()
}
