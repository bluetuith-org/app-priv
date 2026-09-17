package app

import (
	"github.com/ayn2op/tview"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/ui/app/views"
	"github.com/bluetuith-org/bluetuith/ui/config"
)

// Application represents an application.
type Application struct {
	app *tview.Application

	session  bluetooth.Session
	features *appfeatures.FeatureSet
	cfg      *config.Config
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

	v, err := views.NewViews(a)
	if err != nil {
		return err
	}

	a.app = tview.NewApplication(v)

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

var _ views.AppBinder = (*Application)(nil)
