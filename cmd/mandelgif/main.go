package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/arl/mandelgif"
)

// pointA is one of an infinity of interesting points to zoom in.
const pointA = 0.2721950 + 0.00540474i

var defaultCfg = mandelgif.Mandelbrot{
	Width: 256, Height: 256, // image dimension
	NFrames: 50, // number of frame in the animated GIF
	Bounds: mandelgif.Rect{ // 2D-space for the first image
		X0: -2, Y0: -1,
		X1: 1, Y1: 1,
	},
	ZoomLevel: 0.93,   // zoom to apply between a frame and the next one
	ZoomPt:    pointA, // 2D coordinates of the complex number to zoom at
	MaxIter:   1024,   // number of iteration to check if a pixel is in the mandelbrot set
}

func main() {
	m := defaultCfg
	zoomPt := complexValue(m.ZoomPt)
	side := m.Width

	flag.IntVar(&m.NFrames, "frames", m.NFrames, "number of frames in final animation")
	flag.IntVar(&side, "side", side, "image width (square so width=height)")
	flag.Float64Var(&m.ZoomLevel, "zoom", m.ZoomLevel, "scale to apply at each frame (zoom)")
	flag.IntVar(&m.MaxIter, "iter", m.MaxIter, "max iterations to apply on 𝒛")
	flag.Var(&zoomPt, "point", "point to zoom in")
	flag.Parse()

	m.Width, m.Height = side, side
	m.ZoomPt = complex128(zoomPt)

	giff, err := os.Create("out.gif")
	if err != nil {
		log.Fatalln("can't create output file", err)
	}
	defer giff.Close()

	m.RenderAnimatedGif(giff)
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
