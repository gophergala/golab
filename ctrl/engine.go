package ctrl

import (
	"github.com/gophergala/golab/model"
	"github.com/gophergala/golab/view"
	"image"
	"image/draw"
	"math"
	"math/rand"
	"time"
)

// InitNew initializes a new game.
func initNew() {
	// Initialize random number generator
	rand.Seed(time.Now().Unix())

	model.InitNew()
	view.InitNew()
}

// StartEngine starts the game engine in a new goroutine and returns as soon as possible.
func StartEngine() {
	model.NewGameCh <- 1 // Cannot block as application was just started, no incoming requests processed yet

	model.Mutex.Lock()

	go simulate()
}

// simulate implements the game cycle
func simulate() {
	t := time.Now().UnixNano()

	for {
		// Check if we have to start a new game
		select {
		case <-model.NewGameCh:
			initNew()
		default:
		}

		// Process mouse clicks
	clickLoop:
		for {
			select {
			case click := <-model.ClickCh:
				handleClick(click)
			default:
				break clickLoop
			}
		}

		// First clear moving objects from the lab image:
		model.Gopher.EraseImg()
		for _, bd := range model.Bulldogs {
			bd.EraseImg()
		}

		now := time.Now().UnixNano()
		dt := float64(now-t) / 1e9

		// Now step moving objects

		stepMovingObj(model.Gopher, dt)

		// Move Bulldogs
		for _, bd := range model.Bulldogs {
			x, y := int(bd.Pos.X), int(bd.Pos.Y)
			if bd.TargetPos.X == x && bd.TargetPos.Y == y {
				row, col := y/model.BlockSize, x/model.BlockSize
				// Generate new, random target
				// Shuffle the directions slice:
				for i := len(directions) - 1; i > 0; i-- { // last is already random, no use switching with itself
					r := rand.Intn(i + 1)
					directions[i], directions[r] = directions[r], directions[i]
				}
				var drow, dcol int
				for _, dir := range directions {
					switch dir {
					case model.DirLeft:
						dcol = -1
					case model.DirRight:
						dcol = 1
					case model.DirUp:
						drow = -1
					case model.DirDown:
						drow = 1
					}
					if model.Lab[row+drow][col+dcol] == model.BlockEmpty {
						// Direction is good, check if we can even step this way 2 blocks:
						if model.Lab[row+drow*2][col+dcol*2] == model.BlockEmpty {
							drow *= 2
							dcol *= 2
						}
						break
					}
					drow, dcol = 0, 0
				}
				bd.TargetPos.X += dcol * model.BlockSize
				bd.TargetPos.Y += drow * model.BlockSize
			}
			stepMovingObj(bd, dt)
		}

		t = now

		// Sleep some time.
		// Iterations might not be exact, but we don't rely on it:
		// We calculate delta time and calculate moving and next positions
		// based on the delta time.

		model.Mutex.Unlock()              // While sleeping, clients can request view images
		time.Sleep(time.Millisecond * 50) // ~20 FPS
		model.Mutex.Lock()                // We will modify model now, labyrinth image might change so lock.
	}
}

// handleClick handles a mouse click
func handleClick(c model.Click) {
	Gopher := model.Gopher

	// If still moving, wait for it:
	if int(Gopher.Pos.X) != Gopher.TargetPos.X || int(Gopher.Pos.Y) != Gopher.TargetPos.Y {
		return
	}

	// Check if new desired target is in the same row/column and if there is a free passage to there.
	pCol, pRow := int(Gopher.Pos.X)/model.BlockSize, int(Gopher.Pos.Y)/model.BlockSize
	tCol, tRow := c.X/model.BlockSize, c.Y/model.BlockSize

	sorted := func(a, b int) (int, int) {
		if a < b {
			return a, b
		} else {
			return b, a
		}
	}

	if pCol == tCol { // Same column
		for row, row2 := sorted(pRow, tRow); row <= row2; row++ {
			if model.Lab[row][tCol] == model.BlockWall {
				return // Wall in the route
			}
		}
	} else if pRow == tRow { // Same row
		for col, col2 := sorted(pCol, tCol); col <= col2; col++ {
			if model.Lab[tRow][col] == model.BlockWall {
				return // Wall in the route
			}
		}
	} else {
		return // Only the same row or column can be commanded
	}

	// Target pos is allowed and reachable.
	// Use target position rounded to the center of the target block:
	Gopher.TargetPos.X, Gopher.TargetPos.Y = tCol*model.BlockSize+model.BlockSize/2, tRow*model.BlockSize+model.BlockSize/2

	// Mark target position visually for the player if it is not the current block:
	if pRow != tRow || pCol != tCol {
		rect := image.Rect(0, 0, model.BlockSize/4, model.BlockSize/4)
		rect = rect.Add(image.Pt(Gopher.TargetPos.X-rect.Dx()/2, Gopher.TargetPos.Y-rect.Dy()/2))
		draw.Draw(model.LabImg, rect, model.TargetImg, image.Point{}, draw.Over)
	}
}

// stepMovingObj steps the specified MovingObj, properly updating the LabImg.
func stepMovingObj(m *model.MovingObj, dt float64) {
	x, y := int(m.Pos.X), int(m.Pos.Y)

	// Only horizontal or vertical movement is allowed!
	if x != m.TargetPos.X {
		dx := math.Min(dt*model.V, math.Abs(float64(m.TargetPos.X)-m.Pos.X))
		if x > m.TargetPos.X {
			dx = -dx
			m.Direction = model.DirLeft
		} else {
			m.Direction = model.DirRight
		}
		m.Pos.X += dx
	} else if y != m.TargetPos.Y {
		dy := math.Min(dt*model.V, math.Abs(float64(m.TargetPos.Y)-m.Pos.Y))
		if y > m.TargetPos.Y {
			dy = -dy
			m.Direction = model.DirUp
		} else {
			m.Direction = model.DirDown
		}
		m.Pos.Y += dy
	}

	// Draw image at new position
	m.DrawImg()
}

// directions is a reused slice of all directions
var directions = make([]model.Dir, model.DirLength)

func init() {
	// Populate the directions slice
	for i := model.Dir(0); i < model.DirLength; i++ {
		directions[i] = i
	}
}
