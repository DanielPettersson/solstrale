package main

import (
	"syscall/js"

	"github.com/DanielPettersson/wasm-trace/pkg/trace"
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

func raytraceWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		go doTrace(trace.TraceSpecification{
			ImageWidth:  args[0].Get("imageWidth").Int(),
			ImageHeight: args[0].Get("imageHeight").Int(),
			DrawOffsetX: args[0].Get("drawOffsetX").Int(),
			DrawOffsetY: args[0].Get("drawOffsetY").Int(),
			DrawWidth:   args[0].Get("drawWidth").Int(),
			DrawHeight:  args[0].Get("drawHeight").Int(),
		}, args[1])
		return nil
	})
}

func main() {
	WASMTrace := js.ValueOf(make(map[string]interface{}))
	WASMTrace.Set("raytrace", raytraceWrapper())
	js.Global().Set("WASMTrace", WASMTrace)
	<-make(chan bool)
}
