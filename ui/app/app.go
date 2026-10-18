package app

import (
	"github.com/ayn2op/tview"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/config"
	"github.com/bluetuith-org/bluetuith/ui/root"
	"github.com/gdamore/tcell/v3"
)

// Application represents an application.
type Application struct {
	app *tview.Application[root.Model]

	session  bluetooth.Session
	features *appfeatures.FeatureSet
	cfg      *config.Config

	screen tcell.Screen
}

// NewApplication returns a new app.
func NewApplication() *Application {
	return &Application{}
}

// Start starts the application.
func (a *Application) Start(session bluetooth.Session, featureSet *appfeatures.FeatureSet, cfg *config.Config) error {
	a.session = session
	a.features = featureSet
	a.cfg = cfg

	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}

	a.screen = screen
	if err := a.screen.Init(); err != nil {
		return err
	}

	a.screen.EnableMouse()

	a.app = tview.NewApplication(root.New(a), tview.WithScreen(screen))

	return a.app.Run()
}

// Session returns the current session.
func (a *Application) Session() bluetooth.Session {
	return a.session
}

// Features returns the current features of the session.
func (a *Application) Features() *appfeatures.FeatureSet {
	return a.features
}

// Configuration returns the app's configuration.
func (a *Application) Configuration() *config.Config {
	return a.cfg
}

// Authorizer returns the session authorizer.
func (a *Application) Authorizer() bluetooth.SessionAuthorizer {
	return nil
}

// Screen returns the screen for the application.
func (a *Application) Screen() tcell.Screen {
	return a.screen
}
