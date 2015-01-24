package view

import (
	"image/color"
)

// Title as the application - as it appears in the Browser title
const AppTitle = "Gopher's Labyrinth - First Gopher Gala (2015)"

const (
	// Width of the client view in pixels
	ViewWidth = 600
	// Width of the client view in pixels
	ViewHeight = 600
)

// Color "constants"
var (
	Black = color.RGBA{A: 0xff}
	White = color.RGBA{0xff, 0xff, 0xff, 0xff}
)
