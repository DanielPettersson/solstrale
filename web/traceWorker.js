importScripts("wasm_exec.js")

let id;
let specification;

onmessage = function(e) {

    if (e.data.type === "init") {
       
        id = e.data.id;

        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("trace.wasm"), go.importObject).then((result) => {
            go.run(result.instance);			
            postMessage({
                type: "init",
                id: id
            })
        });

    } else if (e.data.type === "start") {
        WASMTrace.raytrace(
            e.data.specification,
            (imageBytes, progress) => {
    
                postMessage(
                    {
                        type: "image",
                        id: id,
                        progress: progress,
                        buffer: imageBytes.buffer,
                        specification: e.data.specification
                    },
                    [
                        imageBytes.buffer
                    ]
                )
            }
        )
    }    
};
