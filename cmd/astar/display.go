package main

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// The Display of the environment.
type Display struct {
	side     int     // Number of cells on each side.
	gap      float32 // Gap between cells.
	cellSize fyne.Size
	markSize fyne.Size
	cells    []*canvas.Rectangle
	marks    []*canvas.Circle
}

// NewDisplay returns a blank display with the given characteristics.
func NewDisplay(side int, size, gap float32) (display *Display, objects []fyne.CanvasObject) {
	display = &Display{
		side:     side,
		gap:      gap,
		cellSize: fyne.NewSize(size, size),
		markSize: fyne.NewSize(size/2, size/2),
	}
	square := side * side
	objects = make([]fyne.CanvasObject, 0, 2*square)
	display.cells = make([]*canvas.Rectangle, square)
	for i := 0; i < square; i++ {
		display.cells[i] = display.newCell(i)
		objects = append(objects, display.cells[i])
	}
	display.marks = make([]*canvas.Circle, square)
	for i := 0; i < square; i++ {
		display.marks[i] = display.newMark(i)
		objects = append(objects, display.marks[i])
	}
	return
}

func (d *Display) newCell(i int) (cell *canvas.Rectangle) {
	cell = canvas.NewRectangle(color.White)
	cell.Resize(d.cellSize)
	cell.Move(d.indexToPosition(i))
	return
}

func (d *Display) newMark(i int) (mark *canvas.Circle) {
	mark = canvas.NewCircle(color.White)
	mark.Resize(d.markSize)
	pos := d.indexToPosition(i)
	pos.X += d.cellSize.Width / 4
	pos.Y += d.cellSize.Height / 4
	mark.Move(pos)
	mark.Hide()
	return
}

// Count returns the total number of cells.
func (d *Display) Count() int { return d.side * d.side }

// Randomise the display.
func (d *Display) Randomise(p *Palette) {
	for i := range d.cells {
		d.cells[i].FillColor = p.Choose()
		d.cells[i].Refresh()
	}
}

// Path puts a mark on the route.
func (d *Display) Path(i int, c color.Color) {
	p := d.marks[i]
	p.FillColor = c
	p.Show()
	p.Refresh()
}

// Destination marks the final destination.
func (d *Display) Destination(i int, c color.Color) {
	p := d.marks[i]
	p.FillColor = c
	p.Show()
	p.Refresh()
}

// ClearRoute clears all the marks.
func (d *Display) ClearRoute() (done bool) {
	for _, v := range d.marks {
		if v.Hidden {
			continue
		}
		done = true
		v.Hide()
		v.Refresh()
	}
	return
}

func (d *Display) indexToColumnRow(i int) (col int, row int) {
	return i % d.side, i / d.side
}

func (d *Display) indexToPosition(i int) fyne.Position {
	c, r := d.indexToColumnRow(i)
	unit := d.cellSize.Width + d.gap
	return fyne.Position{X: float32(c) * unit, Y: float32(r) * unit}
}

func (d *Display) columnRowToIndex(c, r int) int {
	return r*d.side + c
}

// Search is used for planning a route.
type Search struct {
	display *Display
	palette *Palette
	adj     map[int]map[int]float64
}

// SearchWith returns a new search using the palette.
func (d *Display) SearchWith(p *Palette) *Search {
	return &Search{
		display: d,
		palette: p,
		adj:     map[int]map[int]float64{},
	}
}

func (s *Search) adjacent(i int) (adj map[int]float64) {
	adj = make(map[int]float64)
	c, r := s.display.indexToColumnRow(i)
	for rx := r - 1; rx < r+2; rx++ {
		if rx < 0 || rx >= s.display.side {
			continue
		}
		for cx := c - 1; cx < c+2; cx++ {
			if cx < 0 || cx >= s.display.side {
				continue
			}
			if rx == 0 && cx == 0 {
				continue
			}
			//
			// Within bounds rx,cx.
			//
			ix := s.display.columnRowToIndex(cx, rx)
			adj[ix] = 0.0
		}
	}
	return
}

// Adjacent implements the A* interface.
func (s *Search) Adjacent(i int) map[int]float64 {
	m, ok := s.adj[i]
	if ok {
		return m
	}
	//
	// Find all adjacent cells.
	//
	m = s.adjacent(i)
	s.adj[i] = m
	for k := range m {
		//
		// The cell colour is equivalent to the cost.
		//
		c := s.display.cells[k].FillColor
		m[k] = float64(s.palette.Index(c))
	}
	return m
}

// Estimate implements the A* interface.
func (s *Search) Estimate(i, j int) float64 {
	ci, ri := s.display.indexToColumnRow(i)
	cj, rj := s.display.indexToColumnRow(j)
	manhattan := math.Abs(float64(ri-rj)) + math.Abs(float64(ci-cj))
	//
	// Assume an average cost all the way.
	//
	return 2.0 * manhattan
}
