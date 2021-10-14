package mandelgif

import (
	"image"
	gopalette "image/color/palette"
	"image/gif"
	"io"
	"log"
	"math"

	"gonum.org/v1/plot/palette"
)

type Mandelbrot struct {
	MaxIter   int        // maximum number of iterations
	ZoomLevel float64    // zoom applied at each frame
	ZoomPt    complex128 // zoom point
	Bounds    Rect
}

// compute checks if the complex number c is in the Mandelbrot Set.
//
// In theory we should apply an infinity of iterations. In pratice we know that
// if z gets bigger than a predefined number called the 'escape radius', it
// won't get back and just escape farther.
//
// details: runs fc(z) = z² + c a maximum number of 'maxiter' iterations, after
// which it consider c has being in the set (escaped = false) or outside of the
// set (escaped = true), and v quantifies how quickly the values reached the
// escape point.
func compute(c complex128, maxiter int) (v float64, escaped bool) {
	const escapeRadius = 2
	const sqEscapeRadius = escapeRadius * escapeRadius

	var (
		n       int
		z       = 0i
		modulus float64
	)
	for {
		modulus = real(z)*real(z) + imag(z)*imag(z)
		if modulus >= sqEscapeRadius {
			v = float64(n+1) - math.Log(math.Log2(modulus))
			escaped = true
			break
		}
		if n >= maxiter {
			v = float64(maxiter)
			break
		}
		z = z*z + c
		n++
	}

	return v, escaped
}

// renderFrame renders a single image corresponding to cbounds, into the
// provided image pointer.
//
// Steps:
//  1. The value of each pixel must be computed, pixel by pixel, with the
//     relatively compute-intensive 'compute' method.
//  2. Compute a color palette for the frame
//  3. Color every pixel with the palette.
func (m *Mandelbrot) renderFrame(cbounds Rect, img *image.Paletted) {
	values := make([]float64, img.Bounds().Dx()*img.Bounds().Dy())
	histogram := make(map[int]float64, img.Bounds().Dx()*img.Bounds().Dy())

	/* 1. Compute the value for each pixel and record this value into an histogram */

	// For each pixel
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			// Create the complex number corresponding to pixel (x, y)
			c := complex(
				cbounds.X0+(float64(x)*cbounds.Width()/float64(img.Bounds().Dx())),
				cbounds.Y0+(float64(y)*cbounds.Height()/float64(img.Bounds().Dy())),
			)

			value, escaped := compute(c, m.MaxIter)

			// Record the escape value for that pixel
			values[x+y*img.Bounds().Dy()] = value
			if escaped {
				histogram[int(value)]++
			}
		}
	}

	/* 2. Compute a color palette from the histogram of all pixel value */
	var total float64
	for _, v := range histogram {
		total += v
	}
	hues := make([]float64, m.MaxIter+2)
	var h float64
	i := 0
	for ; i < m.MaxIter; i++ {
		h += float64(histogram[i]) / float64(total)
		hues[i] = h
	}
	hues[i], hues[i+1] = h, h

	/*  3. Color every pixel with the palette.*/
	interpolate := func(c1, c2, t float64) float64 {
		return c1*(1-t) + c2*t
	}

	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			mu := values[x+y*img.Bounds().Dy()]
			value := float64(0)
			if mu < float64(m.MaxIter) {
				value = 1
			}

			hsva := palette.HSVA{
				H: interpolate(hues[int(math.Floor(mu))], hues[int(math.Ceil(mu))], math.Mod(mu, 1)),
				S: 1,
				V: value,
				A: 1,
			}
			img.Set(x, y, hsva)
		}
	}
}

// RenderAnimatedGif renders the zoom into the Mandelbrot set as an animated Gif
// image with nframes frames of the specific width and height, into w.
func (m *Mandelbrot) RenderAnimatedGif(w io.Writer, nframes, width, height int) {
	images := make([]*image.Paletted, nframes)
	delays := make([]int, 50)

	log.Printf("Rendering %d frames", nframes)

	// Create the slices of bounds
	bounds := make([]Rect, nframes)
	bounds[0] = m.Bounds
	for i := 1; i < nframes; i++ {
		bounds[i] = bounds[i-1]
		bounds[i].zoom(real(m.ZoomPt), imag(m.ZoomPt), m.ZoomLevel)
	}

	// Render each frame.
	for i := 0; i < nframes; i++ {
		img := image.NewPaletted(image.Rect(0, 0, width, height), gopalette.Plan9)
		m.renderFrame(bounds[i], img)
		images[i] = img
	}

	log.Println("Encoding to GIF")

	gif.EncodeAll(w, &gif.GIF{
		Image: images,
		Delay: delays,
	})
}
