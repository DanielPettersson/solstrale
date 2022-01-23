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

func (q quad) hit(r ray, rayLength interval) (bool, *hitRecord) {
	denom := q.normal.dot(r.direction)

	// No hit if the ray is parallell to the plane
	if math.Abs(denom) < 1e-8 {
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
