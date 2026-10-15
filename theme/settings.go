package theme

import (
	"sync/atomic"

	tint "github.com/lrstanley/bubbletint/v2"
)

func init() {
	tint.DefaultRegistry = tint.NewRegistry(
		tint.TintMoonlightIi,
		tint.TintAlabaster,
		tint.TintITerm2DarkBackground,
		tint.TintITerm2LightBackground,
	)
}

// Settings represents the current settings for themes and other elements.
type Settings struct {
	theme atomic.Pointer[Theme]

	rootCfg *RootConfiguration
	iconSet *IconSet
}

// ParseTheme parses and applies the theme configuration.
func (s *Settings) ParseTheme(asciiIcons bool, config *Configuration) error {
	s.rootCfg = defaultConfig()
	s.iconSet = NewIconSet(asciiIcons)

	if config != nil {
		th, err := s.rootCfg.Merge(config)
		if err != nil {
			return err
		}

		s.theme.Store(th)
		return nil
	}

	th, err := s.rootCfg.convertToTheme()
	if err != nil {
		return err
	}

	s.theme.Store(th)

	return nil
}

// CurrentSettings returns the current theme configuration and settings.
func CurrentSettings() *Settings {
	return _current
}

// Current returns the current theme.
func Current() *Theme {
	return _current.theme.Load()
}

// Icons returns the preconfigured icons.
func Icons() *IconSet {
	return _current.iconSet
}

// _current returns the current theme settings.
var _current = &Settings{}
