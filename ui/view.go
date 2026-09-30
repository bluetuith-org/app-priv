package ui

import (
	"github.com/ayn2op/tview"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/config"
)

// AppBinder represents an interface to interact with the root application itself.
type AppBinder interface {
	// Session returns the current session.
	Session() bluetooth.Session

	// Features returns the current features of the session.
	Features() *appfeatures.FeatureSet

	// Configuration returns the app's configuration.
	Configuration() *config.Config
}

// View represents a single View.
// All views must implement this interface.
type View interface {
	// ViewID returns the view's ID.
	ViewID() ViewID

	// HandleRouterMsg handles the routed message.
	HandleRouterMsg(m RouterMsg) tview.Cmd
}

// Model represents a UI model.
type Model[T any] interface {
	// Init returns a command to run when the model starts, or nil.
	Init() tview.Cmd

	// Update receives messages when this model has focus.
	Update(msg tview.Msg) (T, tview.Cmd)

	// View draws this model onto the screen.
	View(focused bool) tview.Element
}

// RootView represents the root view. All views must inherit this interface,
// to call app-related functions and route messages to components.
type RootView interface {
	AppBinder

	// SendMsg sends a message to the program.
	SendMsg(msg tview.Msg)

	// SendRoutedUpdateMsg routes the message to the specified view.
	// Should only be called from the model's [Model.Update] function.
	SendRoutedUpdateMsg(msg RouterMsg) tview.Cmd
}
