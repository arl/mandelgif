# mandelgif
Generate animated gifs of zoom into the [Mandelbrot Set](https://en.wikipedia.org/wiki/Mandelbrot_set), like this one: 

![](example.gif)


## Installation

To clone the source code:
```
git clone git@github.com:arl/mandelgif.git
```

To import the module as a dependency to your project (library):
```
go get github.com/arl/mandelgif@latest
```

To build and install the `mandelgif` executable on your system:
```
go install github.com/arl/mandelgif/cmd/mandelgif@latest
```


## Usage:
```
mandelgif: renders a zoom of the Mandelbrot set into an animated Gif.

Usage:
	./mandelgif [options] [OUTFILE]

General Options:
	-help                Prints this help message
	-f -frames NUM       Produce an animation with NUM frames. default 50
	-z -zoom   FACTOR    Apply this zoom factor between successive frames. default 0.93)
	-p -point  COMPLEX   Zoom on this point in the complex plane. default "A" (i.e. '0.272195+0.00540474i')
	-w -width  PIXELS    Width of the output GIF image. default 256
	-h -height PIXELS    Height of the output GIF image. default 256
	-i -iter   ITER      Apply a maximum of ITER iterations on 𝒛. default 1024

Notes:
	* OUTFILE defaults to out.gif
	* The zoom point option "--point" accepts one of the following:
		- a complex numbers in the form "x+yi"
		- a letter from "A" to "F". Each letter represents a predefined interesting zoom point
```


## Disclaimer

This program is voluntarily **not** optimized. Making it go faster is let as an exercice to the reader!


## [MIT license](./LICENSE)