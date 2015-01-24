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
	// generate labyrinth
	genLab()

	// Position player to top left corner
	Pos.X = BlockSize + BlockSize/2
	Pos.Y = Pos.X
}

// genLab generates a random labyrinth.
func genLab() {
	// Zero the labyrinth
	for _, row := range Lab {
		for ci := range row {
			row[ci] = BlockEmpty
		}
	}
	
	// Create a "frame":
	for ri := range Lab {
		Lab[ri][0] = BlockWall
		Lab[ri][Cols-1] = BlockWall
	}
	for ci := range Lab[0] {
		Lab[0][ci] = BlockWall
		Lab[Rows-1][ci] = BlockWall
	}
}
