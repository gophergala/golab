package ctrl

import (
	"github.com/gophergala/golab/model"
	"github.com/gophergala/golab/view"
)

// InitNew initializes a new game.
func InitNew() {
	model.InitNew()
	view.InitNew()
}
