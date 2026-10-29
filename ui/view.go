package ui

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/text"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/config"
	"github.com/bluetuith-org/bluetuith/keybindings"
	"github.com/gdamore/tcell/v3"
)

// Model represents a UI model.
type Model interface {
	// Init returns a command to run when the model starts, or nil.
	Init() tview.Cmd

	// Update receives messages when this model has focus.
	Update(msg tview.Msg, rsp *Response)

	// View draws this model onto the screen.
	View(focused bool) tview.Widget
}

// RootView represents the root view. All views must inherit this interface,
// to call app-related functions and route messages to components.
type RootView interface {
	// SetFocused stores the currently focused view.
	SetFocused(id ViewID)

	// PanelFocused returns whether the panel is focused.
	PanelFocused() bool

	// SendMsg sends a message to the program.
	SendMsg(msg tcell.Event)

	// Session returns the current session.
	Session() bluetooth.Session

	// Features returns the current features of the session.
	Features() *appfeatures.FeatureSet

	// Configuration returns the app's configuration.
	Configuration() *config.Config

	// Authorizer returns the authorizer for bluetooth-related permission prompts.
	Authorizer() bluetooth.SessionAuthorizer
}

// TabSection defines a section of the tabs widget where
// the selected view is shown.
type TabSection interface {
	Model

	// ViewID returns the view's ID.
	ViewID() ViewID

	// Title returns the title of the tab section.
	Title() string

	// Icon returns the icon associated with the tab section.
	Icon() string

	// Header returns the content to print in the title bar area.
	Header() []text.Span

	// Status returns the content to print in the status bar area.
	Status() []text.Span

	// ViewKeybind returns the keybinding used to focus the view.
	ViewKeybind() keybindings.Keybinding
}
