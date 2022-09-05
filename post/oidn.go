package post

import (
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/DanielPettersson/solstrale/geo"
	im "github.com/DanielPettersson/solstrale/internal/image"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/mdouchement/hdr"
	"github.com/mdouchement/hdr/codec/pfm"
)

// oidnPostProcessor delegates postprocessing of the rendered image to the Intel Open Image Denoise library
type oidnPostProcessor struct {
	OidnDenoiseExecutablePath string
}

func NewOidn(oidnDenoiseExecutablePath string) (PostProcessor, error) {

	f, err := os.Stat(oidnDenoiseExecutablePath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Oidn path does not exist: %v", err.Error()))
	}

	if !util.IsExecOwner(f.Mode()) {
		return nil, errors.New(fmt.Sprintf("Oidn path is not executable: %v", oidnDenoiseExecutablePath))
	}

	return oidnPostProcessor{
		OidnDenoiseExecutablePath: oidnDenoiseExecutablePath,
	}, nil
}

// PostProcess runs an external command to oidn in a separate process
func (p oidnPostProcessor) PostProcess(
	pixelColors []geo.Vec3,
	albedoColors []geo.Vec3,
	normalColors []geo.Vec3,
	width, height, numSamples int,
) (image.Image, error) {

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
		"--clean_aux",
	)
	err := oidnCmd.Run()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Oidn failed with: %v", err))
	}

	image, _, err := image.Decode(outFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Oidn output decode failed with: %v", err))
	}
	return image, nil
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
