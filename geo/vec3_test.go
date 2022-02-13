package geo

import (
	"math"
	"testing"

	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestZeroVector(t *testing.T) {
	len := ZeroVector.Length()
	assert.Equal(t, float64(0), len)
}

func TestNewVec3(t *testing.T) {
	vec := NewVec3(2, 3, 4)

	assert.Equal(t, float64(2), vec.X)
	assert.Equal(t, float64(3), vec.Y)
	assert.Equal(t, float64(4), vec.Z)
}

func TestRandomVec3(t *testing.T) {
	interval := util.Interval{Min: -2, Max: 2}

	for i := 0; i < 100; i++ {
		vec := RandomVec3(interval.Min, interval.Max)

		assert.True(t, interval.Contains(vec.X))
		assert.True(t, interval.Contains(vec.Y))
		assert.True(t, interval.Contains(vec.Z))
	}
}

func TestRandomInUnitSphere(t *testing.T) {
	for i := 0; i < 100; i++ {
		vec := RandomInUnitSphere()

		assert.True(t, vec.Length() <= 1)
	}
}

func TestRandomUnitVector(t *testing.T) {
	for i := 0; i < 100; i++ {
		vec := RandomUnitVector()

		assert.True(t, math.Abs(vec.Length()-1) < util.AlmostZero)
	}
}

func TestRandomInHemisphere(t *testing.T) {
	for i := 0; i < 100; i++ {
		normal := RandomUnitVector()
		vec := RandomInHemisphere(normal)

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
	for i := 0; i < 100; i++ {
		vec := RandomInUnitDisc()

		assert.True(
			t, vec.Length() <= 1,
			"vec %v is not in unit sphere as length is %v", vec, vec.Length(),
		)
		assert.Equal(t, float64(0), vec.Z)
	}

}

func TestNeg(t *testing.T) {
	vec := RandomInUnitSphere()
	negVec := vec.Neg()

	assert.Equal(t, -vec.X, negVec.X)
	assert.Equal(t, -vec.Y, negVec.Y)
	assert.Equal(t, -vec.Z, negVec.Z)
}

func TestAdd(t *testing.T) {
	vec := RandomInUnitSphere()
	addVec := RandomInUnitSphere()
	resVec := vec.Add(addVec)

	assert.Equal(t, vec.X+addVec.X, resVec.X)
	assert.Equal(t, vec.Y+addVec.Y, resVec.Y)
	assert.Equal(t, vec.Z+addVec.Z, resVec.Z)
}

func TestSub(t *testing.T) {
	vec := RandomInUnitSphere()
	subVec := RandomInUnitSphere()
	resVec := vec.Sub(subVec)

	assert.Equal(t, vec.X-subVec.X, resVec.X)
	assert.Equal(t, vec.Y-subVec.Y, resVec.Y)
	assert.Equal(t, vec.Z-subVec.Z, resVec.Z)
}

func TestMul(t *testing.T) {
	vec := RandomInUnitSphere()
	mulVec := RandomInUnitSphere()
	resVec := vec.Mul(mulVec)

	assert.Equal(t, vec.X*mulVec.X, resVec.X)
	assert.Equal(t, vec.Y*mulVec.Y, resVec.Y)
	assert.Equal(t, vec.Z*mulVec.Z, resVec.Z)
}

func TestMulS(t *testing.T) {
	vec := RandomInUnitSphere()
	mul := util.RandomFloat(-1, 1)
	resVec := vec.MulS(mul)

	assert.Equal(t, vec.X*mul, resVec.X)
	assert.Equal(t, vec.Y*mul, resVec.Y)
	assert.Equal(t, vec.Z*mul, resVec.Z)
}

func TestDivS(t *testing.T) {
	vec := RandomInUnitSphere()
	div := util.RandomFloat(0.5, 1)
	resVec := vec.DivS(div)

	assert.Equal(t, vec.X/div, resVec.X)
	assert.Equal(t, vec.Y/div, resVec.Y)
	assert.Equal(t, vec.Z/div, resVec.Z)
}

func TestDot(t *testing.T) {
	dot := NewVec3(1, 1, 1).Dot(NewVec3(1, 1, 1))
	assert.Equal(t, float64(3), dot)

	dot = NewVec3(0, 0, 1).Dot(NewVec3(0, 0, 0))
	assert.Equal(t, float64(0), dot)

	dot = NewVec3(1, 2, 3).Dot(NewVec3(-2, -3, -4))
	assert.Equal(t, float64(-20), dot)
}

func TestCross(t *testing.T) {
	cross := NewVec3(2, 3, 4).Cross(NewVec3(5, 6, 7))
	assert.Equal(t, NewVec3(-3, 6, -3), cross)
}

func TestLengthSquared(t *testing.T) {
	vec := RandomInUnitSphere()
	assert.True(t, math.Abs(vec.LengthSquared()-math.Pow(vec.Length(), 2)) < util.AlmostZero)
}

func TestLength(t *testing.T) {
	len := NewVec3(0, 3, 4).Length()
	assert.Equal(t, float64(5), len)
}

func TestUnit(t *testing.T) {

	vec := RandomVec3(-10, 10)
	unitVec := vec.Unit()

	assert.True(t, math.Abs(unitVec.Length()-1) < util.AlmostZero)
	assert.True(t, vec.Dot(unitVec) > 0)
}

func TestNearZero(t *testing.T) {
	assert.True(t, ZeroVector.NearZero())
	assert.False(t, RandomVec3(1, 2).NearZero())
}

func TestReflect(t *testing.T) {

	ref := NewVec3(3, -4, 0).Reflect(NewVec3(0, 1, 0))
	assert.Equal(t, NewVec3(3, 4, 0), ref)

	ref = NewVec3(3, -4, 0).Reflect(NewVec3(1, 0, 0))
	assert.Equal(t, NewVec3(-3, -4, 0), ref)
}

func TestRefract(t *testing.T) {
	ref := NewVec3(-3, -3, 0).Unit().Refract(NewVec3(0, 1, 0), 1.0)
	exp := NewVec3(-3, -3, 0).Unit()
	assert.True(t, ref.Sub(exp).NearZero())
}
