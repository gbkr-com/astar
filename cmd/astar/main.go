package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/gbkr-com/astar"
)

func main() {
	//
	// Centre panel.
	//
	application := app.New()
	palette := NewPalette(
		color.NRGBA{102, 152, 229, 255},
		[3]uint8{13, 20, 13},
		color.NRGBA{R: 255, A: 96},
		color.NRGBA{R: 255, A: 255},
	)
	display, objects := NewDisplay(30, 20, 1)
	centre := container.NewWithoutLayout(objects...)
	//
	// Right panel.
	//
	var (
		randButton  *widget.Button
		routeButton *widget.Button
		clearButton *widget.Button
		quitButton  *widget.Button
	)
	randButton = widget.NewButton(
		"Randomise",
		func() {
			display.Randomise(palette)
			centre.Refresh()
		},
	)
	routeButton = widget.NewButton(
		"Route",
		func() {
			srch := display.SearchWith(palette)
			r := astar.Find(srch, 0, display.Count()-1)
			n := len(r)
			if n == 0 {
				fmt.Println("nil")
				return
			}
			for i := 0; i < n; i++ {
				if i == n-1 {
					display.Destination(r[i], palette.Destination())
					break
				}
				display.Path(r[i], palette.Path())
			}
			centre.Refresh()
			routeButton.Disable()
			clearButton.Enable()
		},
	)
	clearButton = widget.NewButton(
		"Clear",
		func() {
			if display.ClearRoute() {
				centre.Refresh()
			}
			routeButton.Enable()
			clearButton.Disable()
		},
	)
	clearButton.Disable()
	quitButton = widget.NewButton(
		"Quit",
		func() {
			application.Quit()
		},
	)
	right := container.New(
		layout.NewVBoxLayout(),
		&layout.Spacer{FixVertical: true},
		randButton,
		&layout.Spacer{FixVertical: true},
		routeButton,
		&layout.Spacer{FixVertical: true},
		clearButton,
		&layout.Spacer{FixVertical: true},
		quitButton,
		&layout.Spacer{FixVertical: true},
	)
	//
	// Launch.
	//
	w := application.NewWindow("Path finding")
	content := container.NewBorder(nil, nil, nil, right, centre)
	w.SetContent(content)
	w.Resize(fyne.NewSize(750, 650))
	// w.SetPadded(false)
	w.ShowAndRun()
}
