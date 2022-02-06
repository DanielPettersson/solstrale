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
	renderImage = image.NewRGBA(image.Rect(0, 0, 400, 300))

	raster := canvas.NewRaster(
		func(w, h int) image.Image {
			return renderImage
		})

	runButton := widget.NewButton("Run", func() {

		renderProgress := make(chan trace.TraceProgress)

		go trace.RayTrace(trace.TraceSpecification{
			ImageWidth:      int(raster.Size().Width),
			ImageHeight:     int(raster.Size().Height),
			DrawOffsetX:     0,
			DrawOffsetY:     0,
			DrawWidth:       int(raster.Size().Width),
			DrawHeight:      int(raster.Size().Height),
			SamplesPerPixel: 100,
			RandomSeed:      123456,
		}, renderProgress)

		for p := range renderProgress {
			renderImage = p.RenderImage
			raster.Refresh()
		}

	})

	topBar := container.New(layout.NewHBoxLayout(), runButton)

	container := container.New(layout.NewBorderLayout(topBar, nil, nil, nil),
		topBar, raster)

	window.SetContent(container)
	window.ShowAndRun()
}
