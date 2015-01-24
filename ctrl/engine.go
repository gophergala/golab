package ctrl

import (
	"github.com/gophergala/golab/model"
	"github.com/gophergala/golab/view"
	"image"
	"image/draw"
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

	go simulate()
}

// simulate implements the game cycle
func simulate() {
	for {
		// Check if we have to start a new game
		select {
		case <-model.NewGameCh:
			initNew()
		default:
		}

		time.Sleep(time.Millisecond * 50) // ~20 FPS

		var moved bool = false

		if moved {
			// Update lab image

			// Clear gopher image from old pos
			x, y := int(model.Pos.X)/model.BlockSize*model.BlockSize, int(model.Pos.Y)/model.BlockSize*model.BlockSize
			rect := image.Rect(x, y, x+model.BlockSize, y+model.BlockSize)
			draw.Draw(model.LabImg, rect, image.NewUniform(model.Black), image.Point{}, draw.Over)
			// Draw gopher
			draw.Draw(model.LabImg, rect, model.GopherImg, image.Point{}, draw.Over)
		}
	}
}
