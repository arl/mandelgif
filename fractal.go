package main

import (
	"image"
	gopalette "image/color/palette"
	"image/gif"
	"io"
	"log"
	"math"

	"gonum.org/v1/plot/palette"
)

type mandelbrot struct {
	width, height int        // rendered image dimensions
	maxiter       int        // maximum number of iterations
	nframes       int        // how many frames to render
	zoomLevel     float64    // zoom applied at each frame
	zoomPt        complex128 // zoom point

	bounds rect
}

func compute(c complex128, maxiter int) (v float64, escaped bool) {
	const escapeRadius = 2
	const sqEscapeRadius = escapeRadius * escapeRadius

	// Check if c is in the mandelbrot set. In theory we should apply an
	// infinity of iterations. In pratice we know that if z gets bigger
	// than a predefined escape radius, it won't get back and just escape farther.
	var (
		n       int
		z       = 0i
		modulus float64
	)
	for {
		// modulus = mod(z)
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
func (m *mandelbrot) renderFrame(cbounds rect, img *image.Paletted) {
	values := make([]float64, img.Bounds().Dx()*img.Bounds().Dy())
	histogram := make(map[int]float64, img.Bounds().Dx()*img.Bounds().Dy())

	/* 1. Compute the value for each pixel and record this value into an histogram */

	// For each pixel
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			// Create the complex number corresponding to pixel (x, y)
			c := complex(
				cbounds.x0+(float64(x)*cbounds.width()/float64(img.Bounds().Dx())),
				cbounds.y0+(float64(y)*cbounds.height()/float64(img.Bounds().Dy())),
			)

			value, escaped := compute(c, m.maxiter)

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
	hues := make([]float64, m.maxiter+2)
	var h float64
	i := 0
	for ; i < m.maxiter; i++ {
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
			if mu < float64(m.maxiter) {
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

// renderAnimatedGifs writes the fractal zoom into w, as an animated Gif image.
// The animation has m.nframes frames (images to be computed).
func (m *mandelbrot) renderAnimatedGif(w io.Writer) {
	images := make([]*image.Paletted, m.nframes)
	delays := make([]int, 50)

	log.Printf("Rendering %d frames", m.nframes)

	// Create the slices of bounds
	bounds := make([]rect, m.nframes)
	bounds[0] = m.bounds
	for i := 1; i < m.nframes; i++ {
		bounds[i] = bounds[i-1]
		bounds[i].zoom(real(m.zoomPt), imag(m.zoomPt), m.zoomLevel)
	}

	// Render each frame.
	for i := 0; i < m.nframes; i++ {
		img := image.NewPaletted(image.Rect(0, 0, m.width, m.height), gopalette.Plan9)
		m.renderFrame(bounds[i], img)
		images[i] = img
	}

	log.Println("Encoding to GIF")

	gif.EncodeAll(w, &gif.GIF{
		Image: images,
		Delay: delays,
	})
}
