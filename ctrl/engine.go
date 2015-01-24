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

		Gopher := model.Gopher

		x, y := int(Gopher.Pos.X), int(Gopher.Pos.Y)

		moved := false

		// Only horizontal or vertical movement is allowed!
		if x != Gopher.TargetPos.X {
			dx := math.Min(dt*model.V, math.Abs(float64(Gopher.TargetPos.X)-Gopher.Pos.X))
			if x > Gopher.TargetPos.X {
				dx = -dx
				Gopher.Direction = model.DirLeft
			} else {
				Gopher.Direction = model.DirRight
			}
			Gopher.Pos.X += dx
			moved = true
		} else if y != Gopher.TargetPos.Y {
			dy := math.Min(dt*model.V, math.Abs(float64(Gopher.TargetPos.Y)-Gopher.Pos.Y))
			if y > Gopher.TargetPos.Y {
				dy = -dy
				Gopher.Direction = model.DirUp
			} else {
				Gopher.Direction = model.DirDown
			}
			Gopher.Pos.Y += dy
			moved = true
		}

		if moved {
			// Update lab image

			// Clear gopher image from old pos
			img := model.GopherImgs[Gopher.Direction]

			b := img.Bounds()
			rect := img.Bounds().Add(image.Pt(x-b.Dx()/2, y-b.Dy()/2))
			draw.Draw(model.LabImg, rect, model.EmptyImg, image.Point{}, draw.Over)

			// Draw gopher at new position
			Gopher.DrawImg()
		}

		t = now
	}
}
