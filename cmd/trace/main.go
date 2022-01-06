package main

import (
	"fmt"
	"syscall/js"
	"time"
)

func sayHello(callback js.Value) {
	for i := 0; i < 10; i++ {
		callback.Invoke(fmt.Sprintf("Hejsan svejsan från WebAssembly %d", i))
		time.Sleep(1 * time.Second)
	}
}

func raytraceWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go sayHello(args[0])
		return nil
	})
}

func main() {
	fmt.Println("Initialized WebAssembly")
	js.Global().Set("raytrace", raytraceWrapper())
	<-make(chan bool)
}
