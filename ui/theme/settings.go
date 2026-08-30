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

	rootCfg RootConfiguration
	iconSet *IconSet
}

// ParseTheme parses and applies the theme configuration.
func (s *ConfigSettings) ParseTheme() error {
	cmpCfg := Configuration(s.rootCfg)

	if err := s.rootCfg.Merge(&cmpCfg); err != nil {
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
	rootCfg: RootConfiguration{
		Tint:          "primer",
		Border:        "bg:from-theme; fg:white",
		BorderFocused: "bg:from-theme; fg:green; attr:bold",
		ADTree: struct {
			Headers           string
			ExpandedIndicator string
			ClosedIndicator   string
			Selection         string
			Adapter           struct{ Present string }
			Device            struct {
				Discovered string
				Paired     string
			}
			DevicesList struct{ Nodes string }
			ActionsList struct{ Nodes string }
		}{
			Headers:           "bg:from-theme; fg:from-theme",
			ExpandedIndicator: "bg:from-theme; fg:from-theme",
			ClosedIndicator:   "bg:from-theme; fg:from-theme",
			Selection:         "bg:from-theme; fg:blue; attr:reverse",
			Adapter: struct{ Present string }{
				Present: "bg:from-theme; fg:from-theme",
			},
			Device: struct {
				Discovered string
				Paired     string
			}{
				Discovered: "bg:from-theme; fg:from-theme",
				Paired:     "bg:from-theme; fg:from-theme",
			},
			DevicesList: struct{ Nodes string }{
				Nodes: "bg:from-theme; fg:from-theme",
			},
			ActionsList: struct{ Nodes string }{
				Nodes: "bg:from-theme; fg:from-theme",
			},
		},
		TabsPane: struct {
			Style      string
			Tab        string
			FocusedTab string
		}{
			Style:      "bg:from-theme; fg:from-theme",
			Tab:        "bg:from-theme; fg:from-theme",
			FocusedTab: "bg:from-theme; fg:brightpurple",
		},
		Info: struct {
			Style   string
			Heading string
		}{
			Style:   "bg:from-theme; fg:from-theme",
			Heading: "bg:blue; fg:white",
		},
		Operations: struct {
			Style   string
			Heading string
		}{
			Style:   "bg:from-theme; fg:from-theme",
			Heading: "bg:blue; fg:white",
		},
		Log: struct {
			Style   string
			Heading string
		}{
			Style:   "bg:from-theme; fg:from-theme",
			Heading: "bg:blue; fg:white",
		},
		StatusBar: struct{ Style string }{
			Style: "bg:from-theme; fg:black",
		},
	},
	Theme: &Theme{},
}
