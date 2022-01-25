package main

import (
	"syscall/js"

	"github.com/DanielPettersson/solstrale/pkg/trace"
)

func doTrace(spec trace.TraceSpecification, callback js.Value) {

	output := make(chan trace.TraceProgress)
	go trace.RayTrace(spec, output)

	for p := range output {
		jsBytes := js.Global().Get("Uint8ClampedArray").New(len(p.ImageData))
		js.CopyBytesToJS(jsBytes, p.ImageData)
		callback.Invoke(
			jsBytes,
			p.Progress,
		)
	}
}

func addTextureWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		name := args[0].Get("name").String()
		width := args[0].Get("width").Int()
		height := args[0].Get("height").Int()

		jsBytes := args[0].Get("data")
		imageBytes := make([]byte, jsBytes.Get("byteLength").Int())
		js.CopyBytesToGo(imageBytes, jsBytes)

		trace.AddTexture(name, width, height, imageBytes)
		return nil
	})
}

func raytraceWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		go doTrace(trace.TraceSpecification{
			ImageWidth:      args[0].Get("imageWidth").Int(),
			ImageHeight:     args[0].Get("imageHeight").Int(),
			DrawOffsetX:     args[0].Get("drawOffsetX").Int(),
			DrawOffsetY:     args[0].Get("drawOffsetY").Int(),
			DrawWidth:       args[0].Get("drawWidth").Int(),
			DrawHeight:      args[0].Get("drawHeight").Int(),
			SamplesPerPixel: args[0].Get("samplesPerPixel").Int(),
			RandomSeed:      args[0].Get("randomSeed").Int(),
		}, args[1])
		return nil
	})
}

func main() {
	WASMTrace := js.ValueOf(make(map[string]interface{}))
	WASMTrace.Set("raytrace", raytraceWrapper())
	WASMTrace.Set("addTexture", addTextureWrapper())
	js.Global().Set("WASMTrace", WASMTrace)
	<-make(chan bool)
}
