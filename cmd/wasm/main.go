package main

import (
	"fmt"
	"syscall/js"

	"github.com/DanielPettersson/wasm-trace/pkg/trace"
)

func doTrace(width int, height int, callback js.Value) {

	output := make(chan trace.TraceProgress)
	go trace.RayTrace(width, height, output)

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
		callback := args[2]

		go doTrace(width, height, callback)
		return nil
	})
}

func main() {
	fmt.Println("Initialized WebAssembly")

	WASMTrace := js.ValueOf(make(map[string]interface{}))
	WASMTrace.Set("raytrace", raytraceWrapper())
	js.Global().Set("WASMTrace", WASMTrace)
	<-make(chan bool)
}
