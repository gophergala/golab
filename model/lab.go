package model

import (
	"image"
	"image/draw"
	"math/rand"
	"sync"
)

// Mutex to be used to synchronize model modifications
var Mutex sync.Mutex

// The model/data of the labyrinth
var Lab [][]Block

// MovingObj is a struct describing a moving object.
type MovingObj struct {
	// The position in the labyrinth in pixel coordinates
	Pos struct {
		X, Y float64
	}

	// Direction where the object is facing toward
	Direction Dir

	// Target position the object is moving to
	TargetPos image.Point

	// Images for each direction, each has zero Min point
	Imgs []*image.RGBA
}

// DrawImg draws the image of the MovingObj to the LabImg.
func (m *MovingObj) DrawImg() {
	m.DrawWithImg(m.Imgs[m.Direction])
}

// EraseImg erases the image of the MovingObj from the LabImg by drawing empty block to it.
func (m *MovingObj) EraseImg() {
	m.DrawWithImg(EmptyImg)
}

// DrawWithImage draws the specified image at the position of the moving object onto the LabImg.
func (m *MovingObj) DrawWithImg(img image.Image) {
	DrawImgAt(img, int(m.Pos.X), int(m.Pos.Y))
}

// DrawImgAt draws the specified image at the specified position which specifies the center of the area to draw.
// The size of the image draw is the block size.
func DrawImgAt(img image.Image, x, y int) {
	r := image.Rect(0, 0, BlockSize, BlockSize).Add(image.Point{x - BlockSize/2, y - BlockSize/2})
	draw.Draw(LabImg, r, img, image.Point{}, draw.Over)
}

// Gopher is our hero, the moving object the user can control.
var Gopher = new(MovingObj)

// Dead tells if Gopher died
var Dead bool

// Tells if we won
var Won bool

// For Gopher we maintain multiple target positions which define a path on which Gopher will move along
var TargetPoss = make([]image.Point, 0, 20)

// Slice of Bulldogs, the ancient enemy of Gophers.
var Bulldogs []*MovingObj

// Exit position
var ExitPos = image.Point{}

// Channel to signal new game
var NewGameCh = make(chan int, 1)

// Constant for the right Mouse button value in the Click struct.
// Button value for left and middle may not be the same for older browsers, but right button always has this value.
const MouseBtnRight = 2

// Click describes a mouse click.
type Click struct {
	// X, Y are the mouse coordinates in pixel, in the coordinate system of the Labyrinth
	X, Y int
	// Btn is the mouse button
	Btn int
}

// Channel to receive mouse clicks on (view package sends, ctrl package (engine) processes them)
var ClickCh = make(chan Click, 10)

// InitNew initializes a new game.
func InitNew() {
	LabImg = image.NewRGBA(image.Rect(0, 0, LabWidth, LabHeight))

	Bulldogs = make([]*MovingObj, int(float64(Rows*Cols)*BulldogDensity/1000))

	Dead = false
	Won = false

	initLab()

	initGopher()

	initBulldogs()

	initLabImg()

	ExitPos.X, ExitPos.Y = (Cols-2)*BlockSize+BlockSize/2, (Rows-2)*BlockSize+BlockSize/2
}

// initLab initializes and generates a new Labyrinth.
func initLab() {
	Lab = make([][]Block, Rows)
	for i := range Lab {
		Lab[i] = make([]Block, Cols)
	}

	// Zero value of the labyrinth is full of empty blocks

	// generate labyrinth
	genLab()
}

// initGopher initializes Gopher.
func initGopher() {
	// Position Gopher to top left corner
	Gopher.Pos.X = BlockSize + BlockSize/2
	Gopher.Pos.Y = Gopher.Pos.X
	Gopher.Direction = DirRight
	Gopher.TargetPos.X, Gopher.TargetPos.Y = int(Gopher.Pos.X), int(Gopher.Pos.Y)
	Gopher.Imgs = GopherImgs

	// Throw away queued targets
	TargetPoss = TargetPoss[0:0]
}

// initBulldogs creates and initializes the Bulldogs.
func initBulldogs() {
	for i := 0; i < len(Bulldogs); i++ {
		bd := new(MovingObj)
		Bulldogs[i] = bd

		// Place bulldog at a random position
		var row, col = int(Gopher.Pos.Y) / BlockSize, int(Gopher.Pos.X) / BlockSize
		// Give some space to Gopher: do not generate Bulldogs too close:
		for gr, gc := row, col; (row-gr)*(row-gr) <= 16 && (col-gc)*(col-gc) <= 16; row, col = rPassPos(0, Rows), rPassPos(0, Cols) {
		}

		bd.Pos.X = float64(col*BlockSize + BlockSize/2)
		bd.Pos.Y = float64(row*BlockSize + BlockSize/2)

		bd.TargetPos.X, bd.TargetPos.Y = int(bd.Pos.X), int(bd.Pos.Y)
		bd.Imgs = BulldogImgs
	}
}

// initLabImg initializes and draws the image of the Labyrinth.
func initLabImg() {
	// Clear the labyrinth image
	draw.Draw(LabImg, LabImg.Bounds(), EmptyImg, image.Pt(0, 0), draw.Over)

	// Draw walls
	zeroPt := image.Point{}
	for ri, row := range Lab {
		for ci, block := range row {
			if block == BlockWall {
				x, y := ci*BlockSize, ri*BlockSize
				rect := image.Rect(x, y, x+BlockSize, y+BlockSize)
				draw.Draw(LabImg, rect, WallImg, zeroPt, draw.Over)
			}
		}
	}
}

// genLab generates a random labyrinth.
func genLab() {
	// Create a "frame":
	for ri := range Lab {
		Lab[ri][0] = BlockWall
		Lab[ri][Cols-1] = BlockWall
	}
	for ci := range Lab[0] {
		Lab[0][ci] = BlockWall
		Lab[Rows-1][ci] = BlockWall
	}

	genLabArea(0, 0, Rows-1, Cols-1)
}

// genLabArea generates a random labyrinth inside the specified area, borders exclusive.
// This is a recursive implementation, each iteration divides the area into 2 parts.
func genLabArea(x1, y1, x2, y2 int) {
	dx, dy := x2-x1, y2-y1

	// Exit condition from the recursion:
	if dx <= 2 || dy <= 2 {
		return
	}

	// Decide if we do a veritcal or horizontal split
	var vert bool
	if dy > dx {
		vert = false
	} else if dx > dy {
		vert = true
	} else if rand.Intn(2) == 0 { // Area is square, choose randomly
		vert = true
	}

	if vert {
		// Add vertical split
		var x int
		if dx > 6 { // To avoid long straight paths, only use random in smaller areas
			x = midWallPos(x1, x2)
		} else {
			x = rWallPos(x1, x2)
		}
		// A whole in it:
		y := rPassPos(y1, y2)
		for i := y1; i <= y2; i++ {
			if i != y {
				Lab[i][x] = BlockWall
			}
		}

		genLabArea(x1, y1, x, y2)
		genLabArea(x, y1, x2, y2)
	} else {
		// Add horizontal split
		var y int
		if dy > 6 { // To avoid long straight paths, only use random in smaller areas
			y = midWallPos(y1, y2)
		} else {
			y = rWallPos(y1, y2)
		}
		// A whole in it:
		x := rPassPos(x1, x2)
		for i := x1; i <= x2; i++ {
			if i != x {
				Lab[y][i] = BlockWall
			}
		}

		genLabArea(x1, y1, x2, y)
		genLabArea(x1, y, x2, y2)
	}
}

// rWallPos returns a random wall position which is an even number between the specified min and max.
func rWallPos(min, max int) int {
	return min + (rand.Intn((max-min)/2-1)+1)*2
}

// midWallPos returns the wall position being at the middle of the specified min and max.
func midWallPos(min, max int) int {
	n := (min + max) / 2
	// make sure it's even
	if n&0x01 == 1 {
		n--
	}
	return n
}

// rPassPos returns a random passage position which is an odd number between the specified min and max.
func rPassPos(min, max int) int {
	return rWallPos(min, max+2) - 1
}
