package hittable

import (
	"math"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/material"
	"github.com/DanielPettersson/solstrale/random"
)

type triangle struct {
	v0     geo.Vec3
	v0v1   geo.Vec3
	v0v2   geo.Vec3
	normal geo.Vec3
	mat    material.Material
	bBox   aabb
	area   float64
}

// NewTriangle creates a new triangle flat hittable object
// A counter clockwise winding is expected
func NewTriangle(v0, v1, v2 geo.Vec3, mat material.Material) Hittable {
	bBox := createAabbFrom3Points(v0, v1, v2).padIfNeeded()
	v0v1 := v1.Sub(v0)
	v0v2 := v2.Sub(v0)
	n := v0v1.Cross(v0v2)
	normal := n.Unit()
	area := n.Length() / 2

	return triangle{
		v0,
		v0v1,
		v0v2,
		normal,
		mat,
		bBox,
		area,
	}
}

func (t triangle) Hit(r geo.Ray, rayLength util.Interval) (bool, *material.HitRecord) {

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
		U:         u,
		V:         v,
		FrontFace: frontFace,
	}

	return true, &rec
}

func (t triangle) BoundingBox() aabb {
	return t.bBox
}

func (t triangle) PdfValue(origin, direction geo.Vec3) float64 {
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

func (t triangle) RandomDirection(origin geo.Vec3) geo.Vec3 {
	p := t.v0.Add(t.v0v1.MulS(random.RandomNormalFloat())).Add(t.v0v2.MulS(random.RandomNormalFloat()))
	return p.Sub(origin)
}

func (t triangle) IsLight() bool {
	switch t.mat.(type) {
	case material.DiffuseLight:
		return true
	default:
		return false
	}
}
