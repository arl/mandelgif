package mandelgif

import "testing"

func Test_rect_translate(t *testing.T) {
	tests := []struct {
		org    Rect
		tx, ty float64
		want   Rect
	}{
		{org: Rect{0, 0, 0, 0}, tx: 1, ty: 1, want: Rect{1, 1, 1, 1}},
	}
	for _, tt := range tests {
		r := tt.org
		r.translate(tt.tx, tt.ty)
		if r != tt.want {
			t.Errorf("translating %v by (%v,%v) = %v, want %v", tt.org, tt.tx, tt.ty, tt.want, r)
		}
	}
}

func Test_rect_scale(t *testing.T) {
	tests := []struct {
		org    Rect
		factor float64
		want   Rect
	}{
		{org: Rect{0, 0, 0, 0}, factor: 10, want: Rect{0, 0, 0, 0}},
		{org: Rect{-1, -2, 0.5, 0.25}, factor: 10, want: Rect{-10, -20, 5, 2.5}},
	}
	for _, tt := range tests {
		r := tt.org
		r.scale(tt.factor)
		if r != tt.want {
			t.Errorf("scaling %v by %v = %v, want %v", tt.org, tt.factor, tt.want, r)
		}
	}
}

func Test_rect_center(t *testing.T) {
	tests := []struct {
		r            Rect
		wantx, wanty float64
	}{
		{r: Rect{0, 0, 0, 0}, wantx: 0, wanty: 0},
		{r: Rect{0, 0, 1, 1}, wantx: 0.5, wanty: 0.5},
		{r: Rect{-1, -1, 1, 1}, wantx: 0, wanty: 0},
	}
	for _, tt := range tests {
		x, y := tt.r.center()
		if x != tt.wantx || y != tt.wanty {
			t.Errorf("center of %v = (%v,%v), want (%v,%v)", tt.r, x, y, tt.wantx, tt.wanty)
		}
	}
}
