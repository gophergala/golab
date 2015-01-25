package model

const (
	// BlockSize is the size of the labyrinth unit in pixels.
	BlockSize = 40
)

var (
	// Rows is the number of rows in the Labyrinth
	Rows int
	// Cols is the number of columns in the Labyrinth
	Cols int

	// LabWidth is the width of the labyrinth's image in pixels.
	LabWidth int

	// LabHeight is the height of the labyrinth's image in pixels.
	LabHeight int
)

// V is the moving speed of Gopher and the Buddlogs in pixel/sec.
var V float64

// "Bulldog density", it tells how many Bulldogs to generate for average of 1,000 blocks.
// For example if this is 10.0 and rows*cols = 21*21 = 441, 10.0*441/1000 = 4.41 => 4 Bulldogs will be generated.
var BulldogDensity float64

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
