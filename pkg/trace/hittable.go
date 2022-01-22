package trace

type hitRecord struct {
	hitPoint  vec3
	normal    vec3
	material  material
	rayLength float64
	u         float64
	v         float64
	frontFace bool
}

type hittable interface {
	hit(r ray, rayLength interval) (bool, *hitRecord)
	boundingBox() aabb
}
