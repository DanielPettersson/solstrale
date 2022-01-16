# Solstråle
A ray tracer implemented in Golang compiled to WebAssembly with a web interface to start and view the render as it progresses.

The rendering is multithreaded in the browser by executing multiple web workers each rendering a part of the image.

Example of output from renderer at about 500 samples per pixel:
![nedladdning (2)](https://user-images.githubusercontent.com/3603911/149679982-4ed90b55-8556-44f6-907c-91edba0f04e1.png)

# Credits
The ray tracing is heavily inspired by the excellent [_Ray Tracing in One Weekend_](https://raytracing.github.io/books/RayTracingInOneWeekend.html)

## Running

```
script/run.sh
```

and then point your browser to localhost:8080
