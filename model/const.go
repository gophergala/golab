package model

const (
	// BlockSize is the size of the labyrinth unit in pixels.
	BlockSize = 40

	// Rows is the number of rows in the labyrinth
	Rows = 33
	// Cols is the number of columns in the labyrinth
	Cols = 33

	// LabWidth is the width of the labyrinth's image in pixels.
	LabWidth = Cols * BlockSize

	// LabHeight is the height of the labyrinth's image in pixels.
	LabHeight = Rows * BlockSize
)

// V is the moving speed of Gopher and the Buddlogs in pixel/sec.
var V float64

// Type of the unit of the labyrinth
type Block int

// Block types of the labyrinth
const (
	// Empty block (free passage)
	BlockEmpty Block = iota
	// Wall block
	BlockWall
)

type Dir int

// Directions of Gopher (facing directions)
const (
	DirRight Dir = iota
	DirLeft
	DirUp
	DirDown

	// Not a valid direction: just to tell how many directions are there
	DirLength
)

func (d Dir) String() string {
	switch d {
	case DirRight:
		return "right"
	case DirLeft:
		return "left"
	case DirUp:
		return "up"
	case DirDown:
		return "down"
	}
	return ""
}
