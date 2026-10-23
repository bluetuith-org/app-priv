package popup

import (
	"github.com/ayn2op/tview/layout"
	"github.com/bluetuith-org/bluetuith/ui"
)

// Popup defines a popup that is displayed over the main view.
type Popup interface {
	ui.Model

	// PopupID returns the unique ID for the popup,
	// for messages to be routed to.
	PopupID() string

	// Width returns the width to be occupied by the popup.
	Width() layout.Length

	// Height returns the height to be occupied by the popup.
	Height() layout.Length
}
