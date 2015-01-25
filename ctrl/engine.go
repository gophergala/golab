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

		// Sleep some time.
		// Iterations might not be exact, but we don't rely on it:
		// We calculate delta time and calculate moving and next positions
		// based on the delta time.

		model.Mutex.Unlock()              // While sleeping, clients can request view images
		time.Sleep(time.Millisecond * 50) // ~20 FPS
		model.Mutex.Lock()                // We will modify model now, labyrinth image might change so lock.

		now := time.Now().UnixNano()
		dt := float64(now-t) / 1e9

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
	}
}

// stepMovingObj steps the specified MovingObj, properly updating the LabImg.
func stepMovingObj(m *model.MovingObj, dt float64) {
	x, y := int(m.Pos.X), int(m.Pos.Y)

	moved := false

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
		moved = true
	} else if y != m.TargetPos.Y {
		dy := math.Min(dt*model.V, math.Abs(float64(m.TargetPos.Y)-m.Pos.Y))
		if y > m.TargetPos.Y {
			dy = -dy
			m.Direction = model.DirUp
		} else {
			m.Direction = model.DirDown
		}
		m.Pos.Y += dy
		moved = true
	}

	if moved {
		// Update lab image

		// Clear image from old pos
		img := m.Imgs[m.Direction]

		b := img.Bounds()
		rect := img.Bounds().Add(image.Pt(x-b.Dx()/2, y-b.Dy()/2))
		draw.Draw(model.LabImg, rect, model.EmptyImg, image.Point{}, draw.Over)

		// Draw image at new position
		m.DrawImg()
	}
}

// directions is a reused slice of all directions
var directions = make([]model.Dir, model.DirLength)

func init() {
	// Populate the directions slice
	for i := model.Dir(0); i < model.DirLength; i++ {
		directions[i] = i
	}
}
