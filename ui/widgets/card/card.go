package card

import (
	"github.com/ayn2op/tview/box"
)

// New returns a new card.
func New(cc *Content) box.Widget {
	return box.New(cc)
}
