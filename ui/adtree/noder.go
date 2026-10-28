package adtree

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/text"
	"github.com/bluetuith-org/bluetuith/ui"
)

// adNoder describes an interface to get information from nodes and update node states.
type adNoder interface {
	Information() []text.Span

	// AdNode returns the base node data.
	AdNode() *adNode

	// HandleKeys handles a [tview.KeyMsg] and returns a message.
	HandleKeys(p tview.KeyMsg) (ui.Router, bool)

	// Refresh refreshes the content of the node.
	Refresh()
}

type emptyNoder struct{}

// AdNode returns the base node data.
func (e *emptyNoder) AdNode() *adNode {
	return nil
}

// HandleKeys handles a [tview.KeyMsg] and returns a message.
func (e *emptyNoder) HandleKeys(tview.KeyMsg) (ui.Router, bool) {
	return ui.EmptyRouterMsg(), false
}

// Information implements [adNoder].
func (e *emptyNoder) Information() []text.Span {
	return nil
}

// Refresh implements [adNoder].
func (e *emptyNoder) Refresh() {
}

var _ adNoder = (*emptyNoder)(nil)
