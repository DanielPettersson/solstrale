package main

import (
	"bytes"
	"fmt"
	"syscall/js"

	"github.com/DanielPettersson/wasm-trace/pkg/trace"
)

func doTrace(width int, height int, imageCallback js.Value, progressCallback js.Value) {

	var buffer bytes.Buffer
	progress := make(chan float32)
	go trace.RayTrace(width, height, progress, &buffer)

	for p := range progress {
		progressCallback.Invoke(p)
	}

	jsBytes := js.Global().Get("Uint8ClampedArray").New(len(buffer.Bytes()))
	js.CopyBytesToJS(jsBytes, buffer.Bytes())
	imageCallback.Invoke(jsBytes)
}

func raytraceWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		width := args[0].Int()
		height := args[1].Int()
		imageCallback := args[2]
		progressCallback := args[3]

		go doTrace(width, height, imageCallback, progressCallback)
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
