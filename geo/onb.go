package geo

import "math"

// Onb is an Orthonormal Basis
type Onb struct {
	U Vec3
	V Vec3
	W Vec3
}

func BuildOnbFromVec3(w Vec3) Onb {
	unitW := w.Unit()

	var a Vec3
	if math.Abs(unitW.X) > .9 {
		a = NewVec3(0, 1, 0)
	} else {
		a = NewVec3(1, 0, 0)
	}
	v := unitW.Cross(a).Unit()
	u := unitW.Cross(v)

	return Onb{
		U: u,
		V: v,
		W: unitW,
	}
}

func (o Onb) Local(a Vec3) Vec3 {
	return o.U.MulS(a.X).Add(o.V.MulS(a.Y)).Add(o.W.MulS(a.Z))
}
