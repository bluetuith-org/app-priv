package theme

import (
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

// ConfigSettings represents the current settings for themes and other elements.
type ConfigSettings struct {
	*Theme

	rootCfg *RootConfiguration
	iconSet *IconSet
}

// ParseTheme parses and applies the theme configuration.
func (s *ConfigSettings) ParseTheme() error {
	cmpCfg := (*Configuration)(s.rootCfg)

	if err := s.rootCfg.Merge(cmpCfg); err != nil {
		return err
	}

	th := s.rootCfg.convertToTheme()
	s.Theme = &th

	s.iconSet = NewIconSet(false)

	return nil
}

// Settings returns the current theme configuration and settings.
func Settings() *ConfigSettings {
	return _current
}

// Current returns the current theme.
func Current() *Theme {
	return _current.Theme
}

// Icons returns the preconfigured icons.
func Icons() *IconSet {
	return _current.iconSet
}

// _current returns the current theme settings.
var _current = &ConfigSettings{
	rootCfg: defaultRootConfig(),
	Theme:   &Theme{},
}
