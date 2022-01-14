importScripts("wasm_exec.js")

const go = new Go();
WebAssembly.instantiateStreaming(fetch("trace.wasm"), go.importObject).then((result) => {
    go.run(result.instance);			
    postMessage({
        type: "init"
    })
});

onmessage = function(e) {
    WASMTrace.raytrace(
        e.data.width, e.data.height, 
        (imageBytes, progress) => {

            postMessage(
                {
                    type: "image",
                    progress: progress,
                    buffer: imageBytes.buffer
                },
                [
                    imageBytes.buffer
                ]
            )
        }
    )
};
