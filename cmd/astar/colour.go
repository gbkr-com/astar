package main

import (
	"image/color"
	"math/rand"
)

// Palette for the display.
type Palette struct {
	shades      [5]color.Color
	path        color.Color
	destination color.Color
}

// NewPalette returns colours for the display.
func NewPalette(first color.Color, decr [3]uint8, path, destination color.Color) (palette *Palette) {
	palette = &Palette{
		path:        path,
		destination: destination,
	}
	for i := range palette.shades {
		switch i {
		case 0:
			palette.shades[i] = first
		case len(palette.shades) - 1:
			palette.shades[i] = color.Black
		default:
			r, g, b, a := palette.shades[i-1].RGBA()
			palette.shades[i] = color.RGBA{
				R: uint8(r) - decr[0],
				G: uint8(g) - decr[1],
				B: uint8(b) - decr[2],
				A: uint8(a),
			}
		}
	}
	return
}

// Path returns the colour for path marks.
func (p *Palette) Path() color.Color { return p.path }

// Destination returns the colour for the destination mark.
func (p *Palette) Destination() color.Color { return p.destination }

// Choose a colour random from the palette.
func (p *Palette) Choose() color.Color {
	return p.shades[rand.Intn(len(p.shades))]
}

// Index returns the index of the given colour.
func (p *Palette) Index(c color.Color) int {
	r, g, b, a := c.RGBA()
	for i, v := range p.shades {
		vr, vg, vb, va := v.RGBA()
		if r == vr && g == vg && b == vb && a == va {
			return i
		}
	}
	return -1
}
