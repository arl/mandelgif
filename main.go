package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// pointA is one of an infinity of interesting points to zoom in.
const pointA = 0.2721950 + 0.00540474i

var defaultCfg = mandelbrot{
	width: 256, height: 256, // image dimension
	nframes: 50, // number of frame in the animated GIF
	bounds: rect{ // 2D-space for the first image
		x0: -2, y0: -1,
		x1: 1, y1: 1,
	},
	zoomLevel: 0.93,   // zoom to apply between a frame and the next one
	zoomPt:    pointA, // 2D coordinates of the complex number to zoom at
	maxiter:   1024,   // number of iteration to check if a pixel is in the mandelbrot set
}

func main() {
	m := defaultCfg
	zoomPt := complexValue(m.zoomPt)
	side := m.width

	flag.IntVar(&m.nframes, "frames", m.nframes, "number of frames in final animation")
	flag.IntVar(&side, "side", side, "image width (square so width=height)")
	flag.Float64Var(&m.zoomLevel, "zoom", m.zoomLevel, "scale to apply at each frame (zoom)")
	flag.IntVar(&m.maxiter, "iter", m.maxiter, "max iterations to apply on 𝒛")
	flag.Var(&zoomPt, "point", "point to zoom in")
	flag.Parse()

	m.width, m.height = side, side
	m.zoomPt = complex128(zoomPt)

	giff, err := os.Create("out.gif")
	if err != nil {
		log.Fatalln("can't create output file", err)
	}
	defer giff.Close()

	m.renderAnimatedGif(giff)
	fmt.Println("success! out.gif")
}

type complexValue complex128

func (c *complexValue) String() string {
	return fmt.Sprint(complex128(*c))
}

func (c *complexValue) Set(s string) error {
	var real, imag float64
	if _, err := fmt.Sscanf(s, "%f+%fi", &real, &imag); err != nil {
		return fmt.Errorf("failed parsing complex number: %v", err)
	}
	*c = complexValue(complex(real, imag))
	return nil
}
