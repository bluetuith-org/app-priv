package app

import (
	tea "charm.land/bubbletea/v2"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/ui/app/views"
	"github.com/bluetuith-org/bluetuith/ui/config"
)

// Application represents an application.
type Application struct{}

// NewApplication returns a new app.
func NewApplication() *Application {
	return &Application{}
}

// Start starts the application.
func (a *Application) Start(session bluetooth.Session, featureSet *appfeatures.FeatureSet, cfg *config.Config) error {
	ab := &appBinder{
		session:    session,
		featureSet: featureSet,
		config:     cfg,
	}

	view, err := views.NewViewModel(ab)
	if err != nil {
		return err
	}

	ab.app = tea.NewProgram(view)
	_, err = ab.app.Run()

	return err
}

// Authorizer returns the session authorizer.
// TODO: Setup.
func (a *Application) Authorizer() bluetooth.SessionAuthorizer {
	return nil
}

type appBinder struct {
	session    bluetooth.Session
	featureSet *appfeatures.FeatureSet
	config     *config.Config

	app *tea.Program
}

func (a *appBinder) Session() bluetooth.Session {
	return a.session
}

func (a *appBinder) Features() *appfeatures.FeatureSet {
	return a.featureSet
}

func (a *appBinder) Configuration() *config.Config {
	return a.config
}

// SendMsg sends a [tea.Msg] to the program.
func (a *appBinder) SendMsg(msg tea.Msg) {
	if a.app == nil {
		return
	}

	a.app.Send(msg)
}
