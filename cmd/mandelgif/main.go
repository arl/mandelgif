package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/arl/mandelgif"
)

const (
	// pointA and pointB are 2 interesting points to zoom in.
	pointA = 0.2721950 + 0.00540474i
	pointB = -1.24254013716898265806 + 0.413238151606368892027i
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
	flag.StringVar(&outname, "o", outname, "")
	flag.StringVar(&outname, "out", outname, "")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `
  -help
        prints this help message  
  -frames | -f int
        number of frames to render in the animation (default 50)
  -height | -h int
        image height (default 256)
  -iter | -i int
        max iterations to apply on 𝒛 (default 1024)
  -out | -o string
        output filename (default "out.gif")
  -point | -p value
        starting point to zoom in (default point A '0.272195+0.00540474i')
  -width | -w int
        image width (default 256)
  -zoom | -z float
        zoom level (i.e scale) to apply between two successive frames (default 0.93)`)
	}

	flag.Parse()

	m.ZoomPt = complex128(zoomPt)

	giff, err := os.Create("out.gif")
	if err != nil {
		log.Fatalln("can't create output file", err)
	}
	defer giff.Close()

	if err := m.Render(giff, nframes, width, height); err != nil {
		log.Fatalf("error: %v", err)
	}
}
