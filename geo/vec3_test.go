package geo

import (
	"fmt"
	"math"
	"testing"

	"github.com/DanielPettersson/solstrale/internal/util"
)

func TestZeroVector(t *testing.T) {
	len := ZeroVector.Length()

	if len != 0 {
		t.Errorf("Length of zero vector should be 0. Is %v", len)
	}
}

func TestNewVec3(t *testing.T) {
	vec := NewVec3(2, 3, 4)

	if vec.X != 2 {
		t.Error("X != 2")
	}

	if vec.Y != 3 {
		t.Error("Y != 2")
	}

	if vec.Z != 4 {
		t.Error("Z != 2")
	}
}

func TestRandomVec3(t *testing.T) {
	interval := util.Interval{Min: -2, Max: 2}

	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("Iteration: %v", i), func(t *testing.T) {
			t.Parallel()
			vec := RandomVec3(interval.Min, interval.Max)
			if !interval.Contains(vec.X) {
				t.Errorf("X %v is out of interval %x", vec.X, interval)
			}
			if !interval.Contains(vec.Y) {
				t.Errorf("Y %v is out of interval %x", vec.Y, interval)
			}
			if !interval.Contains(vec.Z) {
				t.Errorf("Z %v is out of interval %x", vec.Z, interval)
			}
		})
	}
}

func TestRandomInUnitSphere(t *testing.T) {
	vec := RandomInUnitSphere()

	if vec.Length() > 1 {
		t.Errorf("vec %v is not in unit sphere as length is %v", vec, vec.Length())
	}
}

func TestRandomUnitVector(t *testing.T) {
	vec := RandomUnitVector()

	if vec.Length() != 1 {
		t.Errorf("vec %v is not a unit vector as length is %v", vec, vec.Length())
	}
}

func TestRandomInHemisphere(t *testing.T) {
	normal := RandomUnitVector()
	vec := RandomInHemisphere(normal)

	if vec.Length() > 1 {
		t.Errorf("vec %v is not in unit sphere as length is %v", vec, vec.Length())
	}

	if vec.Dot(normal) < 0 {
		t.Errorf("vec %v is not is not pointing in same general direction as normal %v", vec, normal)
	}
}

func TestRandomInUnitDisc(t *testing.T) {
	vec := RandomInUnitDisc()

	if vec.Length() > 1 {
		t.Errorf("vec %v is not in unit sphere as length is %v", vec, vec.Length())
	}

	if vec.Z != 0 {
		t.Errorf("vec Z should be 0 but is %v", vec.Z)
	}

}

func TestNeg(t *testing.T) {
	vec := RandomInUnitSphere()
	negVec := vec.Neg()

	if -vec.X != negVec.X {
		t.Errorf("%v is not negated to %v", vec.X, negVec.X)
	}
	if -vec.Y != negVec.Y {
		t.Errorf("%v is not negated to %v", vec.Y, negVec.Y)
	}
	if -vec.Z != negVec.Z {
		t.Errorf("%v is not negated to %v", vec.Z, negVec.Z)
	}
}

func TestAdd(t *testing.T) {
	vec := RandomInUnitSphere()
	addVec := RandomInUnitSphere()
	resVec := vec.Add(addVec)

	if resVec.X != vec.X+addVec.X {
		t.Errorf("%v is not %v + %v", resVec.X, vec.X, addVec.X)
	}
	if resVec.Y != vec.Y+addVec.Y {
		t.Errorf("%v is not %v + %v", resVec.Y, vec.Y, addVec.Y)
	}
	if resVec.Z != vec.Z+addVec.Z {
		t.Errorf("%v is not %v + %v", resVec.Z, vec.Z, addVec.Z)
	}
}

func TestSub(t *testing.T) {
	vec := RandomInUnitSphere()
	subVec := RandomInUnitSphere()
	resVec := vec.Sub(subVec)

	if resVec.X != vec.X-subVec.X {
		t.Errorf("%v is not %v - %v", resVec.X, vec.X, subVec.X)
	}
	if resVec.Y != vec.Y-subVec.Y {
		t.Errorf("%v is not %v - %v", resVec.Y, vec.Y, subVec.Y)
	}
	if resVec.Z != vec.Z-subVec.Z {
		t.Errorf("%v is not %v - %v", resVec.Z, vec.Z, subVec.Z)
	}
}

func TestMul(t *testing.T) {
	vec := RandomInUnitSphere()
	mulVec := RandomInUnitSphere()
	resVec := vec.Mul(mulVec)

	if resVec.X != vec.X*mulVec.X {
		t.Errorf("%v is not %v * %v", resVec.X, vec.X, mulVec.X)
	}
	if resVec.Y != vec.Y*mulVec.Y {
		t.Errorf("%v is not %v * %v", resVec.Y, vec.Y, mulVec.Y)
	}
	if resVec.Z != vec.Z*mulVec.Z {
		t.Errorf("%v is not %v * %v", resVec.Z, vec.Z, mulVec.Z)
	}
}

func TestMulS(t *testing.T) {
	vec := RandomInUnitSphere()
	mul := util.RandomFloat(-1, 1)
	resVec := vec.MulS(mul)

	if resVec.X != vec.X*mul {
		t.Errorf("%v is not %v * %v", resVec.X, vec.X, mul)
	}
	if resVec.Y != vec.Y*mul {
		t.Errorf("%v is not %v * %v", resVec.Y, vec.Y, mul)
	}
	if resVec.Z != vec.Z*mul {
		t.Errorf("%v is not %v * %v", resVec.Z, vec.Z, mul)
	}
}

func TestDivS(t *testing.T) {
	vec := RandomInUnitSphere()
	div := util.RandomFloat(0.5, 1)
	resVec := vec.DivS(div)

	if resVec.X != vec.X/div {
		t.Errorf("%v is not %v / %v", resVec.X, vec.X, div)
	}
	if resVec.Y != vec.Y/div {
		t.Errorf("%v is not %v / %v", resVec.Y, vec.Y, div)
	}
	if resVec.Z != vec.Z/div {
		t.Errorf("%v is not %v / %v", resVec.Z, vec.Z, div)
	}
}

func TestDot(t *testing.T) {
	dot := NewVec3(1, 1, 1).Dot(NewVec3(1, 1, 1))
	if dot != 3 {
		t.Errorf("Unexpected dot %v", dot)
	}
	dot = NewVec3(0, 0, 1).Dot(NewVec3(0, 0, 0))
	if dot != 0 {
		t.Errorf("Unexpected dot %v", dot)
	}
	dot = NewVec3(1, 2, 3).Dot(NewVec3(-2, -3, -4))
	if dot != -20 {
		t.Errorf("Unexpected dot %v", dot)
	}
}

func TestCross(t *testing.T) {
	cross := NewVec3(2, 3, 4).Cross(NewVec3(5, 6, 7))

	if cross != NewVec3(-3, 6, -3) {
		t.Errorf("Unexpected cross %v", cross)
	}
}

func TestLengthSquared(t *testing.T) {
	vec := RandomInUnitSphere()
	if vec.LengthSquared() != math.Pow(vec.Length(), 2) {
		t.Fail()
	}
}

func TestLength(t *testing.T) {
	len := NewVec3(0, 3, 4).Length()
	if len != 5 {
		t.Errorf("UNexpected length %v", len)
	}
}

func TestUnit(t *testing.T) {

	vec := RandomVec3(-10, 10)
	unitVec := vec.Unit()

	if unitVec.Length() != 1 {
		t.Error("Unit vector should have length 1")
	}

	if vec.Dot(unitVec) < 0 {
		t.Error("Unit vector should point in same direction")
	}
}

func TestNearZero(t *testing.T) {

	if !ZeroVector.NearZero() {
		t.Error("Zero vector should be near zero")
	}

	if RandomVec3(1, 2).NearZero() {
		t.Error("Vector width length is not near zero")
	}
}

func TestReflect(t *testing.T) {

	ref := NewVec3(3, -4, 0).Reflect(NewVec3(0, 1, 0))
	if ref != NewVec3(3, 4, 0) {
		t.Errorf("Not excpected reflection vector %v", ref)
	}

	ref = NewVec3(3, -4, 0).Reflect(NewVec3(1, 0, 0))
	if ref != NewVec3(-3, -4, 0) {
		t.Errorf("Not excpected reflection vector %v", ref)
	}
}

func TestRefract(t *testing.T) {
	ref := NewVec3(-3, -3, 0).Unit().Refract(NewVec3(0, 1, 0), 1.0)
	exp := NewVec3(-3, -3, 0).Unit()
	if !ref.Sub(exp).NearZero() {
		t.Errorf("Not expected refraction vector %v, exp: %v", ref, exp)
	}
}
