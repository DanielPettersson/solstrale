package main

import (
	"fmt"
	"syscall/js"
)

func raytraceWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return "Hejsan svejsan från WebAssembly"
	})
}

func main() {
	fmt.Println("Initialized WebAssembly")
	js.Global().Set("raytrace", raytraceWrapper())
	<-make(chan bool)
}
