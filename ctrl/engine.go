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
		
		model.Mutex.Unlock() // While sleeping, clients can request view images
		time.Sleep(time.Millisecond * 50) // ~20 FPS
		model.Mutex.Lock() // We will modify model now, labyrinth image might change so lock.

		now := time.Now().UnixNano()
		dt := float64(now-t) / 1e9

		x, y := int(model.Pos.X), int(model.Pos.Y)

		moved := false

		// Only horizontal or vertical movement is allowed!
		if x != model.TargetPos.X {
			dx := math.Min(dt*model.V, math.Abs(float64(model.TargetPos.X)-model.Pos.X))
			if x > model.TargetPos.X {
				dx = -dx
			}
			model.Pos.X += dx
			moved = true
		} else if y != model.TargetPos.Y {
			dy := math.Min(dt*model.V, math.Abs(float64(model.TargetPos.Y)-model.Pos.Y))
			if y > model.TargetPos.Y {
				dy = -dy
			}
			model.Pos.Y += dy
			moved = true
		}

		if moved {
			// Update lab image

			// Clear gopher image from old pos
			b := model.GopherImg.Bounds()
			rect := model.GopherImg.Bounds().Add(image.Pt(x-b.Dx()/2, y-b.Dy()/2))
			draw.Draw(model.LabImg, rect, model.EmptyImg, image.Point{}, draw.Over)

			// Draw gopher at new position
			x, y = int(model.Pos.X), int(model.Pos.Y)
			rect = model.GopherImg.Bounds().Add(image.Pt(x-b.Dx()/2, y-b.Dy()/2))
			draw.Draw(model.LabImg, rect, model.GopherImg, image.Point{}, draw.Over)

		}

		t = now
	}
}
