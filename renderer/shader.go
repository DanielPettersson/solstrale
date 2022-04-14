package renderer

import (
	"math"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/hittable"
	"github.com/DanielPettersson/solstrale/material"
	"github.com/DanielPettersson/solstrale/pdf"
)

// Shader calculates the color from a ray hitting a hittable object
type Shader interface {
	Shade(renderer *Renderer, rec *material.HitRecord, ray geo.Ray, depth int) geo.Vec3
}

// PathTracingShader is the full raytracing shader
type PathTracingShader struct {
	MaxDepth int
}

// Shade calculates the color using path tracing
func (pts PathTracingShader) Shade(renderer *Renderer, rec *material.HitRecord, ray geo.Ray, depth int) geo.Vec3 {

	if depth >= pts.MaxDepth {
		return geo.ZeroVector
	}

	emittedColor := rec.Material.Emitted(rec)
	scatter, scatterRecord := rec.Material.Scatter(ray, rec)
	if !scatter {
		return emittedColor
	}

	if scatterRecord.SkipPdf {
		rc, _, _ := renderer.rayColor(scatterRecord.SkipPdfRay, depth+1)
		return scatterRecord.Attenuation.Mul(rc)
	}

	lightPtr := hittable.NewHittablePdf(renderer.lights, rec.HitPoint)
	mixturePdf := pdf.NewMixturePdf(&lightPtr, scatterRecord.PdfPtr)

	scattered := geo.NewRay(
		rec.HitPoint,
		mixturePdf.Generate(),
		ray.Time,
	)
	pdfVal := mixturePdf.Value(scattered.Direction)
	scatteringPdf := rec.Material.ScatteringPdf(ray, rec, scattered)
	rc, _, _ := renderer.rayColor(scattered, depth+1)
	scatterColor := scatterRecord.Attenuation.MulS(scatteringPdf).Mul(rc).DivS(pdfVal)

	return filterInvalidColorValues(emittedColor.Add(scatterColor))
}

func filterInvalidColorValues(col geo.Vec3) geo.Vec3 {
	return geo.NewVec3(
		filterColorValue(col.X),
		filterColorValue(col.Y),
		filterColorValue(col.Z),
	)
}

func filterColorValue(val float64) float64 {
	if math.IsNaN(val) {
		return 0
	}
	// A subjectively chosen value that is a trade off between
	// color acne and suppressing intensity
	if val > 3 {
		return 3
	}
	return val
}

// AlbedoShader outputs flat color
type AlbedoShader struct{}

// Shade calculates the color only attenuation color.
func (AlbedoShader) Shade(renderer *Renderer, rec *material.HitRecord, ray geo.Ray, depth int) geo.Vec3 {
	emittedColor := rec.Material.Emitted(rec)
	scatter, scatterRecord := rec.Material.Scatter(ray, rec)
	if !scatter {
		return emittedColor
	}
	return scatterRecord.Attenuation
}

// NormalShader outputs the normals of the ray hitpoint
type NormalShader struct{}

// Shade calculates the color only using normal.
func (NormalShader) Shade(renderer *Renderer, rec *material.HitRecord, ray geo.Ray, depth int) geo.Vec3 {
	return rec.Normal.Unit()
}

// SimpleShader is a simple shader for quick rendering
type SimpleShader struct{}

// Shade calculates the color only using normal and attenuation color.
func (SimpleShader) Shade(renderer *Renderer, rec *material.HitRecord, ray geo.Ray, depth int) geo.Vec3 {
	emittedColor := rec.Material.Emitted(rec)
	scatter, scatterRecord := rec.Material.Scatter(ray, rec)
	if !scatter {
		return emittedColor
	}

	// Get a factor to multiply attenuation color, range between .25 -> 1.25
	// To get some decent flat shading
	normalFactor := rec.Normal.Unit().Dot(geo.NewVec3(1, 1, -1))*.5 + .75

	return scatterRecord.Attenuation.MulS(normalFactor)
}
