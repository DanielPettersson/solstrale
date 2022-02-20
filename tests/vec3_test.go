package tests

import (
	"math"
	"testing"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/random"
	"github.com/stretchr/testify/assert"
)

func TestZeroVector(t *testing.T) {
	len := geo.ZeroVector.Length()
	assert.Equal(t, float64(0), len)
}

func TestNewVec3(t *testing.T) {
	vec := geo.NewVec3(2, 3, 4)

	assert.Equal(t, float64(2), vec.X)
	assert.Equal(t, float64(3), vec.Y)
	assert.Equal(t, float64(4), vec.Z)
}

func TestRandomVec3(t *testing.T) {
	interval := util.Interval{Min: -2, Max: 2}
	rand := random.NewRandom(0)

	for i := 0; i < 100; i++ {
		vec := geo.RandomVec3(rand, interval.Min, interval.Max)

		assert.True(t, interval.Contains(vec.X))
		assert.True(t, interval.Contains(vec.Y))
		assert.True(t, interval.Contains(vec.Z))
	}
}

func TestRandomInUnitSphere(t *testing.T) {
	rand := random.NewRandom(0)

	for i := 0; i < 100; i++ {
		vec := geo.RandomInUnitSphere(rand)

		assert.True(t, vec.Length() <= 1)
	}
}

func TestRandomUnitVector(t *testing.T) {
	rand := random.NewRandom(0)

	for i := 0; i < 100; i++ {
		vec := geo.RandomUnitVector(rand)

		assert.True(t, math.Abs(vec.Length()-1) < util.AlmostZero)
	}
}

func TestRandomInHemisphere(t *testing.T) {
	rand := random.NewRandom(0)

	for i := 0; i < 100; i++ {
		normal := geo.RandomUnitVector(rand)
		vec := geo.RandomInHemisphere(rand, normal)

		assert.True(
			t, vec.Length() <= 1,
			"vec %v is not in unit sphere as length is %v", vec, vec.Length(),
		)
		assert.True(
			t, vec.Dot(normal) > 0,
			"vec %v is not is not pointing in same general direction as normal %v", vec, normal,
		)
	}
}

func TestRandomInUnitDisc(t *testing.T) {
	rand := random.NewRandom(0)

	for i := 0; i < 100; i++ {
		vec := geo.RandomInUnitDisc(rand)

		assert.True(
			t, vec.Length() <= 1,
			"vec %v is not in unit sphere as length is %v", vec, vec.Length(),
		)
		assert.Equal(t, float64(0), vec.Z)
	}

}

func TestNeg(t *testing.T) {
	rand := random.NewRandom(0)
	vec := geo.RandomInUnitSphere(rand)
	negVec := vec.Neg()

	assert.Equal(t, -vec.X, negVec.X)
	assert.Equal(t, -vec.Y, negVec.Y)
	assert.Equal(t, -vec.Z, negVec.Z)
}

func TestAdd(t *testing.T) {
	rand := random.NewRandom(0)
	vec := geo.RandomInUnitSphere(rand)
	addVec := geo.RandomInUnitSphere(rand)
	resVec := vec.Add(addVec)

	assert.Equal(t, vec.X+addVec.X, resVec.X)
	assert.Equal(t, vec.Y+addVec.Y, resVec.Y)
	assert.Equal(t, vec.Z+addVec.Z, resVec.Z)
}

func TestSub(t *testing.T) {
	rand := random.NewRandom(0)
	vec := geo.RandomInUnitSphere(rand)
	subVec := geo.RandomInUnitSphere(rand)
	resVec := vec.Sub(subVec)

	assert.Equal(t, vec.X-subVec.X, resVec.X)
	assert.Equal(t, vec.Y-subVec.Y, resVec.Y)
	assert.Equal(t, vec.Z-subVec.Z, resVec.Z)
}

func TestMul(t *testing.T) {
	rand := random.NewRandom(0)
	vec := geo.RandomInUnitSphere(rand)
	mulVec := geo.RandomInUnitSphere(rand)
	resVec := vec.Mul(mulVec)

	assert.Equal(t, vec.X*mulVec.X, resVec.X)
	assert.Equal(t, vec.Y*mulVec.Y, resVec.Y)
	assert.Equal(t, vec.Z*mulVec.Z, resVec.Z)
}

func TestMulS(t *testing.T) {
	rand := random.NewRandom(0)
	vec := geo.RandomInUnitSphere(rand)
	mul := rand.RandomFloat(-1, 1)
	resVec := vec.MulS(mul)

	assert.Equal(t, vec.X*mul, resVec.X)
	assert.Equal(t, vec.Y*mul, resVec.Y)
	assert.Equal(t, vec.Z*mul, resVec.Z)
}

func TestDivS(t *testing.T) {
	rand := random.NewRandom(0)
	vec := geo.RandomInUnitSphere(rand)
	div := rand.RandomFloat(0.5, 1)
	resVec := vec.DivS(div)

	assert.Equal(t, vec.X/div, resVec.X)
	assert.Equal(t, vec.Y/div, resVec.Y)
	assert.Equal(t, vec.Z/div, resVec.Z)
}

func TestDot(t *testing.T) {
	dot := geo.NewVec3(1, 1, 1).Dot(geo.NewVec3(1, 1, 1))
	assert.Equal(t, float64(3), dot)

	dot = geo.NewVec3(0, 0, 1).Dot(geo.NewVec3(0, 0, 0))
	assert.Equal(t, float64(0), dot)

	dot = geo.NewVec3(1, 2, 3).Dot(geo.NewVec3(-2, -3, -4))
	assert.Equal(t, float64(-20), dot)
}

func TestCross(t *testing.T) {
	cross := geo.NewVec3(2, 3, 4).Cross(geo.NewVec3(5, 6, 7))
	assert.Equal(t, geo.NewVec3(-3, 6, -3), cross)
}

func TestLengthSquared(t *testing.T) {
	rand := random.NewRandom(0)
	vec := geo.RandomInUnitSphere(rand)
	assert.True(t, math.Abs(vec.LengthSquared()-math.Pow(vec.Length(), 2)) < util.AlmostZero)
}

func TestLength(t *testing.T) {
	len := geo.NewVec3(0, 3, 4).Length()
	assert.Equal(t, float64(5), len)
}

func TestUnit(t *testing.T) {
	rand := random.NewRandom(0)
	vec := geo.RandomVec3(rand, -10, 10)
	unitVec := vec.Unit()

	assert.True(t, math.Abs(unitVec.Length()-1) < util.AlmostZero)
	assert.True(t, vec.Dot(unitVec) > 0)
}

func TestNearZero(t *testing.T) {
	rand := random.NewRandom(0)
	assert.True(t, geo.ZeroVector.NearZero())
	assert.False(t, geo.RandomVec3(rand, 1, 2).NearZero())
}

func TestReflect(t *testing.T) {

	ref := geo.NewVec3(3, -4, 0).Reflect(geo.NewVec3(0, 1, 0))
	assert.Equal(t, geo.NewVec3(3, 4, 0), ref)

	ref = geo.NewVec3(3, -4, 0).Reflect(geo.NewVec3(1, 0, 0))
	assert.Equal(t, geo.NewVec3(-3, -4, 0), ref)
}

func TestRefract(t *testing.T) {
	ref := geo.NewVec3(-3, -3, 0).Unit().Refract(geo.NewVec3(0, 1, 0), 1.0)
	exp := geo.NewVec3(-3, -3, 0).Unit()
	assert.True(t, ref.Sub(exp).NearZero())
}
