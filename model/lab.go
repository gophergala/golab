package model

import (
	"sync"
)

// Mutex to be used to synchronize model modifications
var Mutex sync.Mutex

// The model/data of the labyrinth
var Lab [Rows][Cols]Block

// The player's position in the labyrinth in pixel coordinates
var Pos struct {
	X, Y float64
}

// InitNew initializes a new game.
func InitNew() {
	// Zero the labyrinth
	for _, row := range Lab {
		for ci := range row {
			row[ci] = BlockEmpty
		}
	}

	// TODO generate labyrinth

	// Position player to top left corner
	Pos.X = BlockSize + BlockSize/2
	Pos.Y = Pos.X
}
