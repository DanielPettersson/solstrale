package geo

import (
	"math"

	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/random"
)

var (
	// ZeroVector is a vector with the length of 0
	ZeroVector = Vec3{}
)

// Vec3 is a 3 dimensional vector
type Vec3 struct {
	X float64
	Y float64
	Z float64
}

// NewVec3 creates a new Vec3
func NewVec3(x, y, z float64) Vec3 {
	return Vec3{
		X: x,
		Y: y,
		Z: z,
	}
}

// RandomVec3 creates a random Vec3 within the given interval
func RandomVec3(min float64, max float64) Vec3 {
	return Vec3{
		random.RandomFloat(min, max),
		random.RandomFloat(min, max),
		random.RandomFloat(min, max),
	}
}

// RandomInUnitSphere creates a random Vec3 that is shorter than 1
func RandomInUnitSphere() Vec3 {

	var p Vec3
	for {
		p.X = random.RandomFloat(-1, 1)
		p.Y = random.RandomFloat(-1, 1)
		p.Z = random.RandomFloat(-1, 1)

		if p.LengthSquared() < 1 {
			return p
		}
	}
}

// RandomUnitVector creates a random Vec3 that has the length of 1
func RandomUnitVector() Vec3 {
	return RandomInUnitSphere().Unit()
}

// RandomInHemisphere creates a random Vec3 that is shorter than 1.
// And in the same general direction as given normal.
func RandomInHemisphere(normal Vec3) Vec3 {
	inUnitSphere := RandomInUnitSphere()
	if inUnitSphere.Dot(normal) > 0 {
		return inUnitSphere
	}
	return inUnitSphere.Neg()
}

// RandomInUnitDisc creates a random Vec3 that is shorter than 1
// And that has a Z value of 0
func RandomInUnitDisc() Vec3 {

	var p Vec3
	for {
		p.X = random.RandomFloat(-1, 1)
		p.Y = random.RandomFloat(-1, 1)
		if p.LengthSquared() < 1 {
			return p
		}
	}
}

// Neg returns a Vec3 that has all values negated
func (v Vec3) Neg() Vec3 {
	return Vec3{-v.X, -v.Y, -v.Z}
}

// Add returns a Vec3 that has all values added with corresponding value in given w Vec3
func (v Vec3) Add(w Vec3) Vec3 {
	return Vec3{v.X + w.X, v.Y + w.Y, v.Z + w.Z}
}

// Sub returns a Vec3 that has all values subtracted by corresponding value in given w Vec3
func (v Vec3) Sub(w Vec3) Vec3 {
	return Vec3{v.X - w.X, v.Y - w.Y, v.Z - w.Z}
}

// Mul returns a Vec3 that has all values multiplied with corresponding value in given w Vec3
func (v Vec3) Mul(w Vec3) Vec3 {
	return Vec3{v.X * w.X, v.Y * w.Y, v.Z * w.Z}
}

// MulS returns a Vec3 that has all values multiplied with given scalar t
func (v Vec3) MulS(t float64) Vec3 {
	return Vec3{v.X * t, v.Y * t, v.Z * t}
}

// DivS returns a Vec3 that has all values divided by given scalar t
func (v Vec3) DivS(t float64) Vec3 {
	return Vec3{v.X / t, v.Y / t, v.Z / t}
}

// Dot returns the dot product with given Vec3 w
func (v Vec3) Dot(w Vec3) float64 {
	return v.X*w.X + v.Y*w.Y + v.Z*w.Z
}

// Cross returns the cross product with given Vec3 w
func (v Vec3) Cross(w Vec3) Vec3 {
	return Vec3{
		v.Y*w.Z - v.Z*w.Y,
		v.Z*w.X - v.X*w.Z,
		v.X*w.Y - v.Y*w.X,
	}
}

// LengthSquared return the length squared of the vector
func (v Vec3) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// Length return the length of the vector
func (v Vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

// Unit returns the vector but sized to a length of 1
func (v Vec3) Unit() Vec3 {
	return v.DivS(v.Length())
}

// NearZero checks if the vectors length is near zero
func (v Vec3) NearZero() bool {
	return math.Abs(v.X) < util.AlmostZero && math.Abs(v.Y) < util.AlmostZero && math.Abs(v.Z) < util.AlmostZero
}

// Reflect returns the reflection vector given the normal n
func (v Vec3) Reflect(n Vec3) Vec3 {
	return v.Sub(n.MulS(v.Dot(n) * 2))
}

// Refract returns the refraction vector given the normal n and the index of refraction of the material
func (v Vec3) Refract(n Vec3, indexOfRefraction float64) Vec3 {
	cosTheta := math.Min(v.Neg().Dot(n), 1)
	rOutPerp := n.MulS(cosTheta).Add(v).MulS(indexOfRefraction)
	rOutParallel := n.MulS(-math.Sqrt(math.Abs(1 - rOutPerp.LengthSquared())))
	return rOutPerp.Add(rOutParallel)
}
