package mandelgif

// Rect is a rectangle defined by 2 points in the plane.
//
//    x0              x1
// y0 +---------------+
//    |               |
//    |               |
//    |               |
// y1 +---------------+
type Rect struct {
	X0, Y0 float64
	X1, Y1 float64
}

func (r *Rect) Width() float64  { return r.X1 - r.X0 }
func (r *Rect) Height() float64 { return r.Y1 - r.Y0 }

func (r *Rect) center() (x, y float64) {
	return (r.X0 + r.X1) / 2, (r.Y0 + r.Y1) / 2
}

func (r *Rect) translate(x, y float64) {
	r.X0 += x
	r.X1 += x
	r.Y0 += y
	r.Y1 += y
}

func (r *Rect) scale(factor float64) {
	r.X0 *= factor
	r.Y0 *= factor
	r.X1 *= factor
	r.Y1 *= factor
}

func (r *Rect) zoom(x, y, zfactor float64) {
	orgx, orgy := r.center()
	// translate to the origin
	r.translate(-orgx, -orgy)
	// zoom/scale
	r.scale(zfactor)
	// translate back to the origin
	r.translate(x, y)
}
