package hittable

import (
	"fmt"
	"math"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/material"
)

type quad struct {
	q      geo.Vec3
	u      geo.Vec3
	v      geo.Vec3
	normal geo.Vec3
	d      float64
	w      geo.Vec3
	mat    material.Material
	bBox   aabb
}

func NewQuad(Q geo.Vec3, u geo.Vec3, v geo.Vec3, mat material.Material) Hittable {
	bBox := createAabbFromPoints(Q, Q.Add(u).Add(v)).padIfNeeded()
	n := u.Cross(v)
	normal := n.Unit()
	D := normal.Dot(Q)
	w := n.DivS(n.Dot(n))

	return quad{
		Q,
		u,
		v,
		normal,
		D,
		w,
		mat,
		bBox,
	}
}

func NewBox(a geo.Vec3, b geo.Vec3, mat material.Material) Hittable {

	sides := NewHittableList()

	min := geo.Vec3{X: math.Min(a.X, b.X), Y: math.Min(a.Y, b.Y), Z: math.Min(a.Z, b.Z)}
	max := geo.Vec3{X: math.Max(a.X, b.X), Y: math.Max(a.Y, b.Y), Z: math.Max(a.Z, b.Z)}

	dx := geo.Vec3{X: max.X - min.X, Y: 0, Z: 0}
	dy := geo.Vec3{X: 0, Y: max.Y - min.Y, Z: 0}
	dz := geo.Vec3{X: 0, Y: 0, Z: max.Z - min.Z}

	sides.Add(NewQuad(geo.Vec3{X: min.X, Y: min.Y, Z: max.Z}, dx, dy, mat))
	sides.Add(NewQuad(geo.Vec3{X: max.X, Y: min.Y, Z: max.Z}, dz.Neg(), dy, mat))
	sides.Add(NewQuad(geo.Vec3{X: max.X, Y: min.Y, Z: min.Z}, dx.Neg(), dy, mat))
	sides.Add(NewQuad(geo.Vec3{X: min.X, Y: min.Y, Z: min.Z}, dz, dy, mat))
	sides.Add(NewQuad(geo.Vec3{X: min.X, Y: max.Y, Z: max.Z}, dx, dz.Neg(), mat))
	sides.Add(NewQuad(geo.Vec3{X: min.X, Y: min.Y, Z: min.Z}, dx, dz, mat))

	return &sides
}

func (q quad) Hit(r geo.Ray, rayLength util.Interval) (bool, *material.HitRecord) {
	denom := q.normal.Dot(r.Direction)

	// No hit if the ray is parallell to the plane
	if math.Abs(denom) < util.AlmostZero {
		return false, nil
	}

	// Return false if the hit point parameter t is outside the ray length interval.
	t := (q.d - q.normal.Dot(r.Origin)) / denom
	if !rayLength.Contains(t) {
		return false, nil
	}

	// Determine the hit point lies within the planar shape using its plane coordinates.
	intersection := r.At(t)
	planarHitPointVector := intersection.Sub(q.q)
	alpha := q.w.Dot(planarHitPointVector.Cross(q.v))
	beta := q.w.Dot(q.u.Cross(planarHitPointVector))

	// Is hit point outside of primitive
	if (alpha < 0) || (1 < alpha) || (beta < 0) || (1 < beta) {
		return false, nil
	}

	normal := q.normal
	frontFace := r.Direction.Dot(normal) < 0
	if !frontFace {
		normal = normal.Neg()
	}
	rec := material.HitRecord{
		HitPoint:  intersection,
		Normal:    normal,
		Material:  q.mat,
		RayLength: t,
		U:         alpha,
		V:         beta,
		FrontFace: frontFace,
	}

	return true, &rec
}

func (q quad) BoundingBox() aabb {
	return q.bBox
}

func (q quad) String() string {
	return fmt.Sprintf("Quad(c:%v, u:%v, v:%v)", q.q, q.u, q.v)
}