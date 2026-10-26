package popup

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/inert"
	"github.com/ayn2op/tview/layout"
	"github.com/ayn2op/tview/stack"
	"github.com/bluetuith-org/bluetuith/ui"
)

// Overlay defines a popup that is displayed over the main view.
type Overlay interface {
	ui.Model

	// PopupID returns the unique ID for the popup,
	// for messages to be routed to.
	PopupID() string

	// Width returns the width to be occupied by the popup.
	Width() layout.Length

	// Height returns the height to be occupied by the popup.
	Height() layout.Length
}

// Widget describes a view which draws popups over the main view.
type Widget struct {
	stk stack.Widget

	pm *Manager
}

// New returns a new popup view.
func New(base tview.Widget, pm *Manager) Widget {
	p := Widget{
		stk: stack.New(inert.New(base)),
	}

	for popup := range pm.Popups() {
		p.stk = p.stk.Push(popup)
	}

	return p
}

// Draw implements [tview.Widget].
func (w Widget) Draw(screen tview.Screen, area tview.Rectangle) {
	w.stk.Draw(screen, area)
}

// Handle implements [tview.Widget].
func (w Widget) Handle(msg tview.Msg, area tview.Rectangle) tview.Msg {
	switch ms := msg.(type) {
	case tview.KeyMsg:
		if ui.Kb().CloseItem.Matches(ms) {
			return newDeleteMsg()
		}

	default:
	}

	return w.stk.Handle(msg, area)
}

// Layout implements [tview.Widget].
func (w Widget) Layout(limits layout.Limits) layout.Size {
	return w.stk.Layout(limits)
}

// Size implements [tview.Widget].
func (w Widget) Size() (width layout.Length, height layout.Length) {
	return w.stk.Size()
}
