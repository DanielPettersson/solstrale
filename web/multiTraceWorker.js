let traceWorkers = {};
let width;
let height;
let numWorkers;

onmessage = function(e) {

    if (e.data.type === "init") {

        // Number of workers should be available concurrency subtracted with 
        // this worker thread and browser main thread. But a minimun of 1
        numWorkers = Math.max(e.data.availableConcurrency - 2, 1)
        console.log("Multi trace worker will initialize " + numWorkers + " workers")

        width = e.data.width;
        height = e.data.height;
        
        drawHeight = Math.floor(height / numWorkers)
        drawHeightRemainder = height % numWorkers

        for (let i = 0; i < numWorkers; i++) {

            let traceWorker = new Worker('traceWorker.js');
            traceWorker.postMessage({
                type: "init",
                id: i,
                specification: {
                    imageWidth: width,
                    imageHeight: height,
                    drawOffsetX: 0,
                    drawOffsetY: i * drawHeight,
                    drawWidth: width,
                    drawHeight: drawHeight + (i == numWorkers - 1 ? drawHeightRemainder : 0)
                }
            })
            traceWorker.onerror = function(event) {
                throw event;
            }
            traceWorker.onmessage = function(e) {
                if (e.data.type === "init") {
                    traceWorkers[e.data.id].ready = true;

                    if (Object.keys(traceWorkers).length === numWorkers && 
                        Object.values(traceWorkers).every(w => w.ready)) {
                        postMessage({type: "init"})
                    }
        
                } else if (e.data.type === "image") {

                    traceWorkers[e.data.id].progress = e.data.progress;
                    let progress = Object.values(traceWorkers).map(v => v.progress).reduce((a, b) => a + b, 0) / numWorkers
        
                    postMessage(
                        {
                            type: "image",
                            progress: progress,
                            buffer: e.data.buffer,
                            specification: e.data.specification
                        },
                        [
                            e.data.buffer
                        ]
                    )
                }
            }
        
            traceWorkers[i] = {
                worker: traceWorker,
                ready: false,
                progress: 0.0
            }
        }

    } else if (e.data.type === "start") {

        for (let traceWorker of Object.values(traceWorkers)) {
            traceWorker.worker.postMessage({
                type: "start"
            })			
        }
    }

};


