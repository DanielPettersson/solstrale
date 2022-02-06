package main

import (
	"image"
	"math"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/DanielPettersson/solstrale/pkg/trace"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	app := app.New()
	window := app.NewWindow("Solstråle")
	window.Resize(fyne.Size{
		Width:  800,
		Height: 600,
	})

	var renderImage image.Image
	renderImage = image.NewRGBA(image.Rect(0, 0, 1, 1))

	abortRender := make(chan bool, 1)

	raster := canvas.NewRaster(
		func(w, h int) image.Image {
			return renderImage
		})

	progress := widget.NewProgressBar()

	runButton := widget.Button{
		Text: "Run",
	}
	stopButton := widget.Button{
		Text: "Stop",
	}
	stopButton.Disable()

	runButton.OnTapped = func() {
		runButton.Disable()
		stopButton.Enable()

		renderProgress := make(chan trace.TraceProgress, 1)

		height := int(math.Round(float64(raster.Size().Height)))
		width := int(math.Round(float64(raster.Size().Width)))

		go trace.RayTrace(trace.TraceSpecification{
			ImageWidth:      width,
			ImageHeight:     height,
			SamplesPerPixel: 100,
			RandomSeed:      rand.Int(),
		}, renderProgress, abortRender)

		go func() {
			for p := range renderProgress {
				renderImage = p.RenderImage
				progress.SetValue(p.Progress)
				raster.Refresh()
			}
			runButton.Enable()
			stopButton.Disable()
		}()
	}

	stopButton.OnTapped = func() {
		runButton.Enable()
		stopButton.Disable()
		abortRender <- true
	}

	topBar := container.New(layout.NewHBoxLayout(), &runButton, &stopButton)

	container := container.New(layout.NewBorderLayout(topBar, progress, nil, nil),
		topBar, progress, raster)

	window.SetContent(container)
	window.ShowAndRun()

	abortRender <- true
}
