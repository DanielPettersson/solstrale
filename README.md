# Solstråle
A ray tracer implemented in Golang compiled to WebAssembly with a web interface to start and view the render as it progresses.

The rendering is multithreaded in the browser by executing multiple web workers each rendering a part of the image.

Example of output from renderer at about 2000 samples per pixel:
![nedladdning](https://user-images.githubusercontent.com/3603911/150698759-52881c0d-1ae0-4e6d-8a2a-1f0723102bff.png)

## Credits
The ray tracing is based on the excellent [_Ray Tracing in One Weekend_](https://raytracing.github.io/books/RayTracingInOneWeekend.html) and [_Ray Tracing: The Next Week_](https://raytracing.github.io/books/RayTracingTheNextWeek.html) by Peter Shirley

## Running

```
script/run.sh
```

and then point your browser to localhost:8080
