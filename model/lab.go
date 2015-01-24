package model

import (
	"math/rand"
	"sync"
)

// Mutex to be used to synchronize model modifications
var Mutex sync.Mutex

// The model/data of the labyrinth
var Lab [][]Block

// The player's position in the labyrinth in pixel coordinates
var Pos struct {
	X, Y float64
}

// Channel to signal new game
var NewGameCh = make(chan int, 1)

// InitNew initializes a new game.
func InitNew() {
	Lab = make([][]Block, Rows)
	for i := range Lab {
		Lab[i] = make([]Block, Cols)
	}
	
	// Zero value of the labyrinth is full of empty blocks

	// generate labyrinth
	genLab()

	// Position player to top left corner
	Pos.X = BlockSize + BlockSize/2
	Pos.Y = Pos.X
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
