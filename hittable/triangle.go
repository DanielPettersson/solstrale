package hittable

import (
	"math"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/material"
	"github.com/DanielPettersson/solstrale/random"
)

type Triangle struct {
	v0     geo.Vec3
	v0v1   geo.Vec3
	v0v2   geo.Vec3
	tu0    float64
	tv0    float64
	tu1    float64
	tv1    float64
	tu2    float64
	tv2    float64
	normal geo.Vec3
	mat    material.Material
	bBox   aabb
	area   float64
	center geo.Vec3
}

func NewTriangle(v0, v1, v2 geo.Vec3, mat material.Material) Triangle {
	return NewTriangleWithTexCoords(v0, v1, v2, 0, 0, 0, 0, 0, 0, mat)
}

// NewTriangle creates a new triangle flat hittable object
// A counter clockwise winding is expected
func NewTriangleWithTexCoords(v0, v1, v2 geo.Vec3, tu0, tv0, tu1, tv1, tu2, tv2 float64, mat material.Material) Triangle {
	bBox := createAabbFrom3Points(v0, v1, v2).padIfNeeded()
	v0v1 := v1.Sub(v0)
	v0v2 := v2.Sub(v0)
	n := v0v1.Cross(v0v2)
	normal := n.Unit()
	area := n.Length() / 2

	center := v0.Add(v1).Add(v2).MulS(0.33333)

	return Triangle{
		v0,
		v0v1,
		v0v2,
		tu0,
		tv0,
		tu1,
		tv1,
		tu2,
		tv2,
		normal,
		mat,
		bBox,
		area,
		center,
	}
}

func (t Triangle) Hit(r geo.Ray, rayLength util.Interval) (bool, *material.HitRecord) {

	pVec := r.Direction.Cross(t.v0v2)
	det := t.v0v1.Dot(pVec)

	// No hit if the ray is parallell to the plane
	if math.Abs(det) < util.AlmostZero {
		return false, nil
	}

	invDet := 1 / det
	tVec := r.Origin.Sub(t.v0)
	qVec := tVec.Cross(t.v0v1)

	// Is hit point outside of primitive
	u := tVec.Dot(pVec) * invDet
	if u < 0 || u > 1 {
		return false, nil
	}
	v := r.Direction.Dot(qVec) * invDet
	if v < 0 || u+v > 1 {
		return false, nil
	}

	tt := t.v0v2.Dot(qVec) * invDet
	intersection := r.At(tt)

	// Return false if the hit point parameter t is outside the ray length interval.
	if !rayLength.Contains(tt) {
		return false, nil
	}

	uv0 := (1 - u - v)
	uu := uv0*t.tu0 + u*t.tu1 + v*t.tu2
	vv := uv0*t.tv0 + u*t.tv1 + v*t.tv2

	normal := t.normal
	frontFace := r.Direction.Dot(normal) < 0
	if !frontFace {
		normal = normal.Neg()
	}
	rec := material.HitRecord{
		HitPoint:  intersection,
		Normal:    normal,
		Material:  t.mat,
		RayLength: tt,
		U:         uu,
		V:         vv,
		FrontFace: frontFace,
	}

	return true, &rec
}

func (t Triangle) BoundingBox() aabb {
	return t.bBox
}

// Center returns the center point for the triangle
func (t Triangle) Center(axis int) float64 {
	return t.center.Axis(axis)
}

func (t Triangle) PdfValue(origin, direction geo.Vec3) float64 {
	ray := geo.NewRay(
		origin,
		direction,
		0,
	)

	hit, rec := t.Hit(ray, util.Interval{Min: 0.001, Max: util.Infinity})

	if !hit {
		return 0
	}

	distanceSquared := rec.RayLength * rec.RayLength * direction.LengthSquared()
	cosine := math.Abs(direction.Dot(rec.Normal) / direction.Length())

	return distanceSquared / (cosine * t.area)
}

func (t Triangle) RandomDirection(origin geo.Vec3) geo.Vec3 {
	p := t.v0.Add(t.v0v1.MulS(random.RandomNormalFloat())).Add(t.v0v2.MulS(random.RandomNormalFloat()))
	return p.Sub(origin)
}

func (t Triangle) IsLight() bool {
	return t.mat.IsLight()
}
