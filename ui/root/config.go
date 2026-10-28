package root

import (
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/config"
)

// Config represents an app configuration.
type Config struct {
	session  bluetooth.Session
	features *appfeatures.FeatureSet
	cfg      *config.Config
}

// NewConfig returns a new root app configuration.
func NewConfig(session bluetooth.Session, features *appfeatures.FeatureSet, cfg *config.Config) *Config {
	return &Config{
		session:  session,
		features: features,
		cfg:      cfg,
	}
}
