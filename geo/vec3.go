package geo

import (
	"fmt"
	"math"

	"github.com/DanielPettersson/solstrale/internal/util"
)

var (
	AlmostZero float64 = 1e-8
	ZeroVector         = Vec3{}
)

type Vec3 struct {
	X float64
	Y float64
	Z float64
}

func NewVec3(x, y, z float64) Vec3 {
	return Vec3{
		X: x,
		Y: y,
		Z: z,
	}
}

func RandomVec3(min float64, max float64) Vec3 {
	return Vec3{
		util.RandomFloat(min, max),
		util.RandomFloat(min, max),
		util.RandomFloat(min, max),
	}
}

func RandomInUnitSphere() Vec3 {

	var p Vec3
	for {
		p.X = util.RandomFloat(-1, 1)
		p.Y = util.RandomFloat(-1, 1)
		p.Z = util.RandomFloat(-1, 1)

		if p.LengthSquared() < 1 {
			return p
		}
	}
}

func RandomUnitVector() Vec3 {
	return RandomInUnitSphere().Unit()
}

func RandomInHemisphere(normal Vec3) Vec3 {
	inUnitSphere := RandomInUnitSphere()
	if inUnitSphere.Dot(normal) > 0 {
		return inUnitSphere
	} else {
		return inUnitSphere.Neg()
	}
}

func RandomInUnitDisc() Vec3 {

	var p Vec3
	for {
		p.X = util.RandomFloat(-1, 1)
		p.Y = util.RandomFloat(-1, 1)
		if p.LengthSquared() < 1 {
			return p
		}
	}
}

func (v Vec3) Neg() Vec3 {
	return Vec3{-v.X, -v.Y, -v.Z}
}

func (v Vec3) Add(w Vec3) Vec3 {
	return Vec3{v.X + w.X, v.Y + w.Y, v.Z + w.Z}
}

func (v Vec3) Sub(w Vec3) Vec3 {
	return Vec3{v.X - w.X, v.Y - w.Y, v.Z - w.Z}
}

func (v Vec3) Mul(w Vec3) Vec3 {
	return Vec3{v.X * w.X, v.Y * w.Y, v.Z * w.Z}
}

func (v Vec3) MulS(t float64) Vec3 {
	return Vec3{v.X * t, v.Y * t, v.Z * t}
}

func (v Vec3) DivS(t float64) Vec3 {
	return Vec3{v.X / t, v.Y / t, v.Z / t}
}

func (v Vec3) Dot(w Vec3) float64 {
	return v.X*w.X + v.Y*w.Y + v.Z*w.Z
}

func (v Vec3) Cross(w Vec3) Vec3 {
	return Vec3{
		v.Y*w.Z - v.Z*w.Y,
		v.Z*w.X - v.X*w.Z,
		v.X*w.Y - v.Y*w.X,
	}
}

func (v Vec3) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

func (v Vec3) Unit() Vec3 {
	return v.DivS(v.Length())
}

func (v Vec3) NearZero() bool {
	return math.Abs(v.X) < AlmostZero && math.Abs(v.Y) < AlmostZero && math.Abs(v.Z) < AlmostZero
}

func (v Vec3) Reflect(n Vec3) Vec3 {
	return v.Sub(n.MulS(v.Dot(n) * 2))
}

func (v Vec3) Refract(n Vec3, indexOfRefraction float64) Vec3 {
	cosTheta := math.Min(v.Neg().Dot(n), 1)
	rOutPerp := n.MulS(cosTheta).Add(v).MulS(indexOfRefraction)
	rOutParallel := n.MulS(-math.Sqrt(math.Abs(1 - rOutPerp.LengthSquared())))
	return rOutPerp.Add(rOutParallel)
}

func (v Vec3) String() string {
	return fmt.Sprintf("[%f, %f, %f]", v.X, v.Y, v.Z)
}
