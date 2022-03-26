// Package post provides post processor for the rendered image
package post

import (
	_ "embed"
	"image"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/DanielPettersson/solstrale/geo"
	im "github.com/DanielPettersson/solstrale/internal/image"
	"github.com/mdouchement/hdr"
	"github.com/mdouchement/hdr/codec/pfm"
	_ "github.com/mdouchement/hdr/codec/pfm"
)

// PostProcessor is responsible for taking the rendered image and transforming it
type PostProcessor interface {
	PostProcess(
		pixelColors []geo.Vec3,
		albedoColors []geo.Vec3,
		normalColors []geo.Vec3,
		width, height, numSamples int,
	) *image.Image
}

// OidnPostProcessor delegates postprocessing of the rendered image to the Intel Open Image Denoise library
type OidnPostProcessor struct {
	OidnDenoiseExecutablePath string
}

// PostProcess runs an external command to oidn in a separate process
func (p OidnPostProcessor) PostProcess(
	pixelColors []geo.Vec3,
	albedoColors []geo.Vec3,
	normalColors []geo.Vec3,
	width, height, numSamples int,
) *image.Image {

	ldrFile, _ := ioutil.TempFile("", "*.pfm")
	defer os.Remove(ldrFile.Name())
	pfm.Encode(ldrFile, toHdrImage(pixelColors, width, height, numSamples))

	albFile, _ := ioutil.TempFile("", "*.pfm")
	defer os.Remove(albFile.Name())
	pfm.Encode(albFile, toHdrImage(albedoColors, width, height, numSamples))

	nrmFile, _ := ioutil.TempFile("", "*.pfm")
	defer os.Remove(nrmFile.Name())
	pfm.Encode(nrmFile, toHdrImage(normalColors, width, height, numSamples))

	outFile, _ := ioutil.TempFile("", "*.pfm")
	defer os.Remove(outFile.Name())

	oidnCmd := exec.Command(
		p.OidnDenoiseExecutablePath,
		"--ldr", ldrFile.Name(),
		"--alb", albFile.Name(),
		"--nrm", nrmFile.Name(),
		"-o", outFile.Name(),
	)
	oidnCmd.Run()

	image, _, _ := image.Decode(outFile)
	return &image
}

func toHdrImage(pixels []geo.Vec3, width, height, numSamples int) hdr.Image {
	pixelData := make([]float64, len(pixels)*3)
	for i := 0; i < len(pixels); i++ {
		pix := im.ToFloat(pixels[i], numSamples)
		pixelData[i*3] = pix.X
		pixelData[i*3+1] = pix.Y
		pixelData[i*3+2] = pix.Z
	}

	return &hdr.RGB64{
		Pix:    pixelData,
		Stride: 3 * width,
		Rect: image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: width, Y: height},
		},
	}
}
