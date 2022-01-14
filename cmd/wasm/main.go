package main

import (
	"syscall/js"

	"github.com/DanielPettersson/wasm-trace/pkg/trace"
)

func doTrace(width, height, interlaceSize, interlaceOffset int, callback js.Value) {

	output := make(chan trace.TraceProgress)
	specification := trace.TraceSpecification{
		ImageWidth:      width,
		ImageHeight:     height,
		InterlaceSize:   interlaceSize,
		InterlaceOffset: interlaceOffset,
	}

	go trace.RayTrace(specification, output)

	for p := range output {
		jsBytes := js.Global().Get("Uint8ClampedArray").New(len(p.ImageData))
		js.CopyBytesToJS(jsBytes, p.ImageData)
		callback.Invoke(jsBytes, p.Progress)
	}
}

func raytraceWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		width := args[0].Int()
		height := args[1].Int()
		interlaceSize := args[2].Int()
		interlaceOffset := args[3].Int()
		callback := args[4]

		go doTrace(width, height, interlaceSize, interlaceOffset, callback)
		return nil
	})
}

func main() {
	WASMTrace := js.ValueOf(make(map[string]interface{}))
	WASMTrace.Set("raytrace", raytraceWrapper())
	js.Global().Set("WASMTrace", WASMTrace)
	<-make(chan bool)
}
