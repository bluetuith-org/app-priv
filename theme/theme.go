package theme

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ayn2op/tview"
	"github.com/gdamore/tcell/v3"
	tc "github.com/gdamore/tcell/v3/color"
	tint "github.com/lrstanley/bubbletint/v2"
)

//go:generate go run ../cmd/structgen theme

const (
	fgKeyword   = "fg"
	bgKeyword   = "bg"
	attrKeyword = "attr"

	propSegmentSep = ";"
	propValSep     = ":"
	multiValSep    = ","

	colorSpecOpen  = "("
	colorSpecClose = ")"

	colorDarkenKeyword  = "darken"
	colorLightenKeyword = "lighten"
)

const (
	minBufferLen = 20
	maxBuffenLen = 200
)

// RootConfiguration represents the root configuration.
type RootConfiguration Configuration

// Merge merges the provided configuration with the root configuration.
func (r *RootConfiguration) Merge(cfg *Configuration) (*Theme, error) {
	tintName := cfg.Tint

	tintSpec := tint.DefaultTintsByID(tintName)
	if tintSpec == nil {
		return nil, fmt.Errorf("the specified theme was not found: %s", tintName)
	}

	tint.Register(tintSpec)
	tint.SetTint(tintSpec)

	r.Tint = tintName
	theme := &Theme{}

	for parseInfo := range iterProperties(r, cfg, theme) {
		cmpStyle, err := r.applyTheme(parseInfo.cmpCfg)
		if err != nil {
			return nil, err
		}

		rootStyle, err := r.applyTheme(*parseInfo.rootCfg)
		if err != nil {
			return nil, err
		}

		style := rootStyle
		if cmpStyle != (tcell.Style{}) {
			style = tview.MergeStyle(rootStyle, cmpStyle)
		}

		*parseInfo.style = style
	}

	return theme, nil
}

// convertToTheme converts the configuration to a [Theme].
func (r *RootConfiguration) convertToTheme() (*Theme, error) {
	theme := &Theme{}

	tintSpec := tint.DefaultTintsByID(r.Tint)
	if tintSpec == nil {
		return nil, fmt.Errorf("the specified theme was not found: %s", r.Tint)
	}

	tint.Register(tintSpec)
	tint.SetTint(tintSpec)

	for parseInfo := range iterProperties(r, nil, theme) {
		style, err := r.applyTheme(*parseInfo.rootCfg)
		if err != nil {
			return nil, err
		}

		*parseInfo.style = style
	}

	return theme, nil
}

func (r *RootConfiguration) applyTheme(cfgItem string) (tcell.Style, error) {
	style := tcell.Style{}
	fgSet, bgSet := false, false

	if cfgItem == "" {
		goto SetStyle
	}

	for seg := range strings.SplitSeq(removeSpaces(cfgItem), propSegmentSep) {
		prop, val, ok := strings.Cut(seg, propValSep)
		if !ok {
			return style, fmt.Errorf("invalid property sequence: '%s' -> '%s' ('%s')", prop, val, cfgItem)
		}

		if val == "" {
			return style, fmt.Errorf("no value was specified for property '%s' ('%s')", prop, cfgItem)
		}

		switch prop {
		case fgKeyword:
			c, err := r.parseColorSpec(val)
			if err != nil {
				return style, fmt.Errorf("%w: Invalid format '%s' ('%s')", err, prop, cfgItem)
			}
			if c == _noColor {
				c = tc.FromImageColor(tint.Current().Fg)
			}

			style = style.Foreground(c)
			fgSet = true

		case bgKeyword:
			c, err := r.parseColorSpec(val)
			if err != nil {
				return style, fmt.Errorf("%w: Invalid format '%s' ('%s')", err, prop, cfgItem)
			}

			if c == _noColor {
				c = tc.FromImageColor(tint.Current().Bg)
			}

			style = style.Background(c)
			bgSet = true

		case attrKeyword:
			const (
				totalAttrs = 6
			)

			var attrCount int

			for attr := range strings.SplitSeq(val, multiValSep) {
				if attr == "" {
					continue
				}

				switch attr {
				case "bold":
					style = style.Bold(true)

				case "italic":
					style = style.Italic(true)

				case "underline":
					style = style.Underline(true)

				case "blink":
					style = style.Blink(true)

				case "reverse":
					style = style.Reverse(true)

				case "dim":
					style = style.Dim(true)

				default:
					return style, fmt.Errorf("invalid attribute '%s' specified for property '%s' ('%s')", attr, prop, cfgItem)
				}

				attrCount++
				if attrCount > totalAttrs {
					return style, fmt.Errorf("invalid property length of '%s' in '%s' ('%s')", val, prop, cfgItem)
				}
			}
		}
	}

SetStyle:
	if !fgSet {
		style = style.Foreground(tc.FromImageColor(tint.Current().Fg))
	}

	if !bgSet {
		style = style.Background(tc.FromImageColor(tint.Current().Bg))
	}

	return style, nil
}

func (r *RootConfiguration) parseColorSpec(colorSpec string) (tc.Color, error) {
	if colorSpec == "" {
		return _noColor, fmt.Errorf("no format specified, color specification is empty")
	}

	str := colorSpec

	colorName := colorSpec
	spec := ""

	var darken, lighten float64

	openIdx := strings.Index(str, colorSpecOpen)
	switch {
	case openIdx < 0:
		goto MatchColor

	case openIdx == 0:
		return _noColor, fmt.Errorf("no color specified: '%s'", colorSpec)

	case openIdx > 0:
		colorName = colorName[:openIdx]

		closeIdx := strings.Index(colorSpec, colorSpecClose)
		if closeIdx < 0 {
			return _noColor, fmt.Errorf("no closing parenthesis for expression '%s'", colorSpec)
		}

		spec = colorSpec[openIdx+1 : closeIdx]
	}

	for spec := range strings.SplitSeq(spec, multiValSep) {
		specProp, specVal, ok := strings.Cut(spec, propValSep)
		if !ok {
			return _noColor, fmt.Errorf("invalid color specifier sequence: '%s'", spec)
		}

		switch specProp {
		case colorDarkenKeyword:
			dark, err := strconv.ParseFloat(specVal, 64)
			if err != nil {
				return _noColor, fmt.Errorf("invalid number for darken: '%s'", colorSpec)
			}

			darken = dark

		case colorLightenKeyword:
			light, err := strconv.ParseFloat(specVal, 64)
			if err != nil {
				return _noColor, fmt.Errorf("invalid number for lighten: '%s'", colorSpec)
			}

			lighten = light

		default:
			return _noColor, fmt.Errorf("invalid format: '%s' ('%s')", spec, colorSpec)
		}
	}

MatchColor:
	c, ok := parseColor(colorName)
	if !ok {
		return _noColor, fmt.Errorf("invalid color specifier '%s'", colorName)
	}

	if darken != 0 {
		c = darkenColor(c, darken)
	}

	if lighten != 0 {
		c = lightenColor(c, lighten)
	}

	return tc.FromImageColor(c), nil
}

func removeSpaces(s string) string {
	var sb strings.Builder

	sb.Grow(max(minBufferLen, min(len(s), maxBuffenLen)))

	for field := range strings.FieldsSeq(s) {
		sb.WriteString(field)
	}

	return sb.String()
}
