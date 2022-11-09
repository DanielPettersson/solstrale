package post

import (
	"errors"
	"fmt"
	"image"
	"image/color"

	"github.com/DanielPettersson/solstrale/geo"
	im "github.com/DanielPettersson/solstrale/internal/image"
)

type bloomPostProcessor struct {
	blurRadius      float64
	bloomMultiplier float64
}

// NewBloom returns a new bloom post processor with the following properties.
// blurRadius: controls the number of blur iterations applied and thus the size of the bloom
// bloomMultiplier: a multipler for how intensive the bloom effect is
func NewBloom(blurRadius, bloomMultiplier float64) PostProcessor {
	return bloomPostProcessor{
		blurRadius:      blurRadius,
		bloomMultiplier: bloomMultiplier,
	}
}

var BLUR_KERNEL = [][]float64{
	{.05, .15, .05},
	{.15, .20, .15},
	{.05, .15, .05},
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

	if b.blurRadius < 0 || b.blurRadius > 1 {
		return nil, errors.New(fmt.Sprintf("Invalid blurRadius %v given", b.blurRadius))
	}

	pixelCount := width * height
	workColors := make([]geo.Vec3, pixelCount)

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
	blurIterations := int(float64(width) * b.blurRadius)
	for blurCount := 0; blurCount < blurIterations; blurCount++ {
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

	for i := 0; i < pixelCount; i++ {
		workColors[i] = workColors[i].MulS(b.bloomMultiplier)
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
