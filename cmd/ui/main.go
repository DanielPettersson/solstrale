package main

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/DanielPettersson/solstrale/pkg/trace"
)

func main() {
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

		renderProgress := make(chan trace.TraceProgress, 2)

		go trace.RayTrace(trace.TraceSpecification{
			ImageWidth:      int(raster.Size().Width),
			ImageHeight:     int(raster.Size().Height),
			DrawOffsetX:     0,
			DrawOffsetY:     0,
			DrawWidth:       int(raster.Size().Width),
			DrawHeight:      int(raster.Size().Height),
			SamplesPerPixel: 10,
			RandomSeed:      123456,
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
