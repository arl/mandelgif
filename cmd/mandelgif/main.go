package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/arl/mandelgif"
)

const (
	// Predefined interesting points to zoom in.
	pointA = 0.2721950 + 0.00540474i
	pointB = -1.24254013716898265806 + 0.413238151606368892027i
	pointC = -0.743904874255535 - 0.1317119067802009i // Sea Horse
	pointD = -0.761574 - 0.0847596i                   // Spirals
	pointE = -1.62917 - 0.0203968i
	pointF = 0.42884 - 0.231345i

	usage = `mandelgif: renders a zoom of the Mandelbrot set into an animated Gif.

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
		- a letter from "A" to "F". Each letter represents a predefined interesting zoom point`
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("[mandelgif] ")

	m := mandelgif.Mandelbrot{
		Bounds: mandelgif.Rect{ // 2D-space for the first image
			X0: -2, Y0: -1,
			X1: 1, Y1: 1,
		},
		ZoomLevel: 0.93,   // zoom to apply between a frame and the next one
		ZoomPt:    pointA, // 2D coordinates of the complex number to zoom at
		MaxIter:   1024,   // number of iteration to check if a pixel is in the mandelbrot set
	}
	zoomPt := complexFlag(m.ZoomPt)
	height := 256
	width := 256
	nframes := 50
	outname := "out.gif"

	flag.IntVar(&nframes, "f", nframes, "")
	flag.IntVar(&nframes, "frames", nframes, "")
	flag.IntVar(&height, "h", height, "")
	flag.IntVar(&height, "height", height, "")
	flag.IntVar(&width, "w", width, "")
	flag.IntVar(&width, "width", width, "")
	flag.Float64Var(&m.ZoomLevel, "z", m.ZoomLevel, "")
	flag.Float64Var(&m.ZoomLevel, "zoom", m.ZoomLevel, "")
	flag.Var(&zoomPt, "p", "")
	flag.Var(&zoomPt, "point", "")
	flag.IntVar(&m.MaxIter, "i", m.MaxIter, "")
	flag.IntVar(&m.MaxIter, "iter", m.MaxIter, "")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usage)
	}

	flag.Parse()

	if flag.NArg() == 1 {
		outname = flag.Arg(0)
	} else if flag.NArg() != 0 {
		flag.Usage()
		os.Exit(1)
	}

	m.ZoomPt = complex128(zoomPt)

	giff, err := os.Create(outname)
	if err != nil {
		log.Fatalln("can't create output file", err)
	}
	defer giff.Close()

	if err := m.Render(giff, nframes, width, height); err != nil {
		log.Fatalf("error: %v", err)
	}
}
