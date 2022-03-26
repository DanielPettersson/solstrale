// Package post provides post processor for the rendered image
package post

import (
	_ "embed"
	"image"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/mdouchement/hdr"
	"github.com/mdouchement/hdr/codec/pfm"
	_ "github.com/mdouchement/hdr/codec/pfm"
)

// PostProcessor is responsible for taking the rendered image and transforming it
type PostProcessor interface {
	PostProcess(input []geo.Vec3, width, height int) (bool, *image.Image)
}

// NopPostProcessor is a post processor that does nothing
type NopPostProcessor struct{}

// PostProcess just returns the input
func (NopPostProcessor) PostProcess(input []geo.Vec3, width, height int) (bool, *image.Image) {
	return false, nil
}

// OidnPostProcessor delegates postprocessing of the rendered image to the Intel Open Image Denoise library
type OidnPostProcessor struct {
	OidnDenoiseExecutablePath string
}

// PostProcess runs an external command to oidn in a separate process
func (p OidnPostProcessor) PostProcess(input []geo.Vec3, width, height int) (bool, *image.Image) {

	inputData := make([]float64, len(input)*3)
	for i := 0; i < len(input); i++ {
		inputData[i*3] = input[i].X
		inputData[i*3+1] = input[i].Y
		inputData[i*3+2] = input[i].Z
	}

	hdrImg := hdr.RGB64{
		Pix:    inputData,
		Stride: 3 * width,
		Rect: image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: width, Y: height},
		},
	}

	oidnInputFile, err := ioutil.TempFile("", "*.pfm")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(oidnInputFile.Name())
	oidnOutputFile, err := ioutil.TempFile("", "*.pfm")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(oidnOutputFile.Name())

	pfm.Encode(oidnInputFile, &hdrImg)

	oidnCmd := exec.Command(p.OidnDenoiseExecutablePath, "--ldr", oidnInputFile.Name(), "-o", oidnOutputFile.Name())
	err = oidnCmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	image, _, err := image.Decode(oidnOutputFile)
	if err != nil {
		log.Fatal(err)
	}

	return true, &image
}
