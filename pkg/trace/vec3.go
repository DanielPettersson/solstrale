package trace

import (
	"fmt"
	"math"
)

type vec3 struct {
	x float64
	y float64
	z float64
}

func randomVec3(min float64, max float64) vec3 {
	return vec3{
		randomFloat(min, max),
		randomFloat(min, max),
		randomFloat(min, max),
	}
}

func randomInUnitSphere() vec3 {
	for {
		p := randomVec3(-1, 1)
		if p.lengthSquared() < 1 {
			return p
		}
	}
}

func randomUnitVector() vec3 {
	return randomInUnitSphere().unit()
}

func randomInHemisphere(normal vec3) vec3 {
	inUnitSphere := randomInUnitSphere()
	if inUnitSphere.dot(normal) > 0 {
		return inUnitSphere
	} else {
		return inUnitSphere.neg()
	}
}

func randomInUnitDisc() vec3 {
	for {
		p := vec3{randomFloat(-1, 1), randomFloat(-1, 1), 0}
		if p.lengthSquared() < 1 {
			return p
		}
	}
}

func (v vec3) neg() vec3 {
	return vec3{-v.x, -v.y, -v.z}
}

func (v vec3) add(w vec3) vec3 {
	return vec3{v.x + w.x, v.y + w.y, v.z + w.z}
}

func (v vec3) sub(w vec3) vec3 {
	return vec3{v.x - w.x, v.y - w.y, v.z - w.z}
}

func (v vec3) mul(w vec3) vec3 {
	return vec3{v.x * w.x, v.y * w.y, v.z * w.z}
}

func (v vec3) mulS(t float64) vec3 {
	return vec3{v.x * t, v.y * t, v.z * t}
}

func (v vec3) divS(t float64) vec3 {
	return vec3{v.x / t, v.y / t, v.z / t}
}

func (v vec3) dot(w vec3) float64 {
	return v.x*w.x + v.y*w.y + v.z*w.z
}

func (v vec3) cross(w vec3) vec3 {
	return vec3{
		v.y*w.z - v.z*w.y,
		v.z*w.x - v.x*w.z,
		v.x*w.y - v.y*w.x,
	}
}

func (v vec3) lengthSquared() float64 {
	return v.x*v.x + v.y*v.y + v.z*v.z
}

func (v vec3) length() float64 {
	return math.Sqrt(v.lengthSquared())
}

func (v vec3) unit() vec3 {
	return v.divS(v.length())
}

func (v vec3) nearZero() bool {
	s := 1e-8
	return math.Abs(v.x) < s && math.Abs(v.y) < s && math.Abs(v.z) < s
}

func (v vec3) reflect(n vec3) vec3 {
	return v.sub(n.mulS(v.dot(n) * 2))
}

func (v vec3) refract(n vec3, indexOfRefraction float64) vec3 {
	cosTheta := math.Min(v.neg().dot(n), 1)
	rOutPerp := n.mulS(cosTheta).add(v).mulS(indexOfRefraction)
	rOutParallel := n.mulS(-math.Sqrt(math.Abs(1 - rOutPerp.lengthSquared())))
	return rOutPerp.add(rOutParallel)
}

func (v vec3) String() string {
	return fmt.Sprintf("[%f, %f, %f]", v.x, v.y, v.z)
}
