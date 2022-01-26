package trace

import (
	"fmt"
	"math"
)

type quad struct {
	Q      vec3
	u      vec3
	v      vec3
	normal vec3
	D      float64
	w      vec3
	mat    material
	bBox   aabb
}

func createQuad(Q vec3, u vec3, v vec3, mat material) quad {
	bBox := createAabbFromPoints(Q, Q.add(u).add(v)).padIfNeeded()
	n := u.cross(v)
	normal := n.unit()
	D := normal.dot(Q)
	w := n.divS(n.dot(n))

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

func createBox(a vec3, b vec3, mat material) hittable {

	sides := emptyHittableList()

	min := vec3{math.Min(a.x, b.x), math.Min(a.y, b.y), math.Min(a.z, b.z)}
	max := vec3{math.Max(a.x, b.x), math.Max(a.y, b.y), math.Max(a.z, b.z)}

	dx := vec3{max.x - min.x, 0, 0}
	dy := vec3{0, max.y - min.y, 0}
	dz := vec3{0, 0, max.z - min.z}

	sides.add(createQuad(vec3{min.x, min.y, max.z}, dx, dy, mat))
	sides.add(createQuad(vec3{max.x, min.y, max.z}, dz.neg(), dy, mat))
	sides.add(createQuad(vec3{max.x, min.y, min.z}, dx.neg(), dy, mat))
	sides.add(createQuad(vec3{min.x, min.y, min.z}, dz, dy, mat))
	sides.add(createQuad(vec3{min.x, max.y, max.z}, dx, dz.neg(), mat))
	sides.add(createQuad(vec3{min.x, min.y, min.z}, dx, dz, mat))

	return &sides
}

func (q quad) hit(r ray, rayLength interval) (bool, *hitRecord) {
	denom := q.normal.dot(r.direction)

	// No hit if the ray is parallell to the plane
	if math.Abs(denom) < almostZero {
		return false, nil
	}

	// Return false if the hit point parameter t is outside the ray length interval.
	t := (q.D - q.normal.dot(r.origin)) / denom
	if !rayLength.contains(t) {
		return false, nil
	}

	// Determine the hit point lies within the planar shape using its plane coordinates.
	intersection := r.at(t)
	planarHitPointVector := intersection.sub(q.Q)
	alpha := q.w.dot(planarHitPointVector.cross(q.v))
	beta := q.w.dot(q.u.cross(planarHitPointVector))

	// Is hit point outside of primitive
	if (alpha < 0) || (1 < alpha) || (beta < 0) || (1 < beta) {
		return false, nil
	}

	normal := q.normal
	frontFace := r.direction.dot(normal) < 0
	if !frontFace {
		normal = normal.neg()
	}
	rec := hitRecord{
		hitPoint:  intersection,
		normal:    normal,
		material:  q.mat,
		rayLength: t,
		u:         alpha,
		v:         beta,
		frontFace: frontFace,
	}

	return true, &rec
}

func (q quad) boundingBox() aabb {
	return q.bBox
}

func (q quad) String() string {
	return fmt.Sprintf("Quad(c:%v, u:%v, v:%v)", q.Q, q.u, q.v)
}
