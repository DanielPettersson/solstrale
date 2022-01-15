# Solstråle
A ray tracer implemented in Golang compiled to WebAssembly with a web interface to start and view the render as it progresses.

The rendering is multithreaded in the browser by executing multiple web workers each rendering a part of the image.

Still basic rendering, as can be seen in the example. Only diffuse and metallic materials and no light sources.
![nedladdning (1)](https://user-images.githubusercontent.com/3603911/149635353-1e6892ef-72ca-4c22-a477-551fd1bb870a.png)

## Running

```
script/run.sh
```

and then point your browser to localhost:8080
