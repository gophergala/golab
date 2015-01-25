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

// LoopDelay is the delay between the iterations of the main loop of the game engine, in milliseconds.
var LoopDelay = 50 // ~20 FPS

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

// Delta time since our last iteration
var dt float64

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

		// First erase target images. We have to do this before handling mouse clicks
		// as they may change the target positions
		eraseDrawTargetPoss(true)

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

		// Next clear moving objects from the lab image:
		model.Gopher.EraseImg()
		for _, bd := range model.Bulldogs {
			bd.EraseImg()
		}

		model.DrawImgAt(model.ExitImg, model.ExitPos.X, model.ExitPos.X)

		if !model.Dead {
			// Draw target positions
			eraseDrawTargetPoss(false)
		}

		now := time.Now().UnixNano()
		dt = float64(now-t) / 1e9

		// Now step moving objects

		stepGopher()
		stepBulldogs()

		// Check if Gopher reached the exit point
		if int(model.Gopher.Pos.X) == model.ExitPos.X && int(model.Gopher.Pos.Y) == model.ExitPos.Y {
			handleWinning()
		}

		t = now

		// Sleep some time.
		// Iterations might not be exact, but we don't rely on it:
		// We calculate delta time and calculate moving and next positions
		// based on the delta time.

		model.Mutex.Unlock() // While sleeping, clients can request view images
		if model.Won {
			// If won, nothing has to be done, just wait for a new game signal
			<-model.NewGameCh // Blocking receive
			// Send back value to detect it at the proper place
			model.NewGameCh <- 1
		}
		time.Sleep(time.Millisecond * time.Duration(LoopDelay))
		model.Mutex.Lock() // We will modify model now, labyrinth image might change so lock.
	}
}

// handleClick handles a mouse click
func handleClick(c model.Click) {
	if model.Dead {
		return
	}

	Gopher := model.Gopher

	if c.Btn == model.MouseBtnRight {
		model.TargetPoss = model.TargetPoss[0:0]
	}

	// If target buffer is full, do nothing:
	if len(model.TargetPoss) == cap(model.TargetPoss) {
		return
	}

	// Last target pos:
	var TargetPos image.Point
	if len(model.TargetPoss) == 0 {
		TargetPos = Gopher.TargetPos
	} else {
		TargetPos = model.TargetPoss[len(model.TargetPoss)-1]
	}

	// Check if new desired target is in the same row/column as the last target and if there is a free passage to there.
	pCol, pRow := TargetPos.X/model.BlockSize, TargetPos.Y/model.BlockSize
	tCol, tRow := c.X/model.BlockSize, c.Y/model.BlockSize

	// sorted simply returns its parameters in ascendant order:
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
	model.TargetPoss = append(model.TargetPoss, image.Pt(tCol*model.BlockSize+model.BlockSize/2, tRow*model.BlockSize+model.BlockSize/2))
}

// eraseDrawTargetPoss either erases or draws target positions of Gopher, both the current and the buffered ones.
func eraseDrawTargetPoss(erase bool) {
	var img image.Image
	if erase {
		img = model.EmptyImg
	} else {
		img = model.TargetImg
	}
	// dtp: drawTargetPos
	dtp := func(TargetPos image.Point) {
		rect := model.TargetImg.Bounds()
		rect = rect.Add(image.Pt(TargetPos.X-rect.Dx()/2, TargetPos.Y-rect.Dy()/2))
		draw.Draw(model.LabImg, rect, img, image.Point{}, draw.Over)
	}

	dtp(model.Gopher.TargetPos)
	for _, TargetPos := range model.TargetPoss {
		dtp(TargetPos)
	}
}

// stepGopher handles moving the Gopher and also handles the multiple target positions of Gopher.
func stepGopher() {
	Gopher := model.Gopher

	if model.Dead {
		Gopher.DrawWithImg(model.DeadImg)
		return
	}

	// Check if reached current target position:
	if int(Gopher.Pos.X) == Gopher.TargetPos.X && int(Gopher.Pos.Y) == Gopher.TargetPos.Y {
		// Check if we have more target positions in our path:
		if len(model.TargetPoss) > 0 {
			// Set the next target as the current
			Gopher.TargetPos = model.TargetPoss[0]
			// and remove it from the targets:
			model.TargetPoss = model.TargetPoss[:copy(model.TargetPoss, model.TargetPoss[1:])]
		}
	}

	// Step Gopher
	stepMovingObj(Gopher)
}

// stepBulldogs iterates over all Bulldogs, generates new random target if they reached their current, and steps them.
func stepBulldogs() {
	// Gopher's position:
	gpos := model.Gopher.Pos

	for _, bd := range model.Bulldogs {
		x, y := int(bd.Pos.X), int(bd.Pos.Y)

		if bd.TargetPos.X == x && bd.TargetPos.Y == y {
			row, col := y/model.BlockSize, x/model.BlockSize
			// Generate new, random target.
			// For this we shuffle all the directions, and check them sequentially.
			// Firts one in which direction there is a free path wins (such path surely exists).

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
					// Direction is good, check if we can even step 2 bocks in this way:
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

		stepMovingObj(bd)

		if !model.Dead {
			// Check if this Bulldog reached Gopher
			if math.Abs(gpos.X-bd.Pos.X) < model.BlockSize*0.75 && math.Abs(gpos.Y-bd.Pos.Y) < model.BlockSize*0.75 {
				handleDying()
			}
		}
	}
}

// handleDying handles the death of Gopher event.
func handleDying() {
	model.Dead = true
}

// handleWinning handles the winning of game event.
func handleWinning() {
	model.Won = true

	r := model.WonImg.Bounds()
	r = r.Add(image.Point{view.Pos.X + view.ViewWidth/2 - r.Dx()/2, view.Pos.Y + view.ViewHeight/2 - r.Dy()/2})
	draw.Draw(model.LabImg, r, model.WonImg, image.Point{}, draw.Over)
}

// stepMovingObj steps the specified MovingObj and draws its image to its new position onto the LabImg.
func stepMovingObj(m *model.MovingObj) {
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
