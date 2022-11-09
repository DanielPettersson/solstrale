package post

import (
	"image"
	"image/color"
	"math"

	"github.com/DanielPettersson/solstrale/geo"
	im "github.com/DanielPettersson/solstrale/internal/image"
)

type bloomPostProcessor struct {
}

// NewBloom returns a new bloom post processor instance
func NewBloom() PostProcessor {
	return bloomPostProcessor{}
}

var BLUR_KERNEL = [][]float64{
	{.10, .15, .10},
	{.15, .20, .15},
	{.10, .15, .10},
}

func imageIndex(x, y, width int) int {
	return y*width + x
}

// PostProcess applies a bloom filter to the renderered image
func (b bloomPostProcessor) PostProcess(
	pixelColors []geo.Vec3,
	albedoColors []geo.Vec3,
	normalColors []geo.Vec3,
	width, height, numSamples int,
) (image.Image, error) {

	pixelCount := width * height
	workColors := make([]geo.Vec3, pixelCount)
	maxIntensity := 0.

	// make an image of bright areas
	for i := 0; i < pixelCount; i++ {
		if pixelColors[i].Length() > float64(numSamples) {
			workColors[i] = pixelColors[i]
		} else {
			workColors[i] = geo.Vec3{}
		}
	}

	// apply blur
	blurSize := len(BLUR_KERNEL)
	for blurCount := 0; blurCount < 200; blurCount++ {
		tmpColors := make([]geo.Vec3, pixelCount)
		for y := 1; y < height-1; y++ {
			for x := 1; x < width-1; x++ {
				val := geo.Vec3{}
				for by := 0; by < blurSize; by++ {
					for bx := 0; bx < blurSize; bx++ {
						tmpVal := workColors[imageIndex(x+bx-1, y+by-1, width)].MulS(BLUR_KERNEL[by][bx])
						val = val.Add(tmpVal)
					}
				}
				tmpColors[imageIndex(x, y, width)] = val
			}
		}
		for i := 0; i < pixelCount; i++ {
			workColors[i] = tmpColors[i]
		}
	}

	// normalize

	for i := 0; i < pixelCount; i++ {
		maxIntensity = math.Max(maxIntensity, workColors[i].Length())
	}

	for i := 0; i < pixelCount; i++ {
		workColors[i] = workColors[i].DivS(maxIntensity).MulS(float64(numSamples)).MulS(3)
	}

	// Create output
	ret := make([]color.RGBA, pixelCount)
	for i := 0; i < pixelCount; i++ {
		ret[i] = im.ToRgba(pixelColors[i].Add(workColors[i]), numSamples)
	}
	img := im.RenderImage{
		ImageWidth:  width,
		ImageHeight: height,
		Data:        ret,
	}

	return img, nil
}
