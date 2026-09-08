package theme

import (
	"fmt"
	"strings"
	"unicode"

	"charm.land/lipgloss/v2"
	tint "github.com/lrstanley/bubbletint/v2"
)

//go:generate go run ../../cmd/structgen theme

const (
	fgKeyword        = "fg"
	bgKeyword        = "bg"
	attrKeyword      = "attr"
	fromThemeKeyword = "from-theme"

	propSegmentSep = ";"
	propValSep     = ":"
	multiValSep    = ","
)

const (
	minBufferLen = 20
	maxBuffenLen = 200
)

// Configuration represents the app's theme configuration.
type Configuration struct {
	Tint string `themetype:"primer"`

	Border        string `themedef:"bg:from-theme; fg:white"`
	BorderFocused string `themedef:"bg:from-theme; fg:green; attr:bold"`

	TitleBar string `themedef:"bg:purple;fg:white;attr:bold"`

	ADTree struct {
		Headers string `themedef:"bg:from-theme; fg:from-theme"`

		ExpandedIndicator string `themedef:"bg:from-theme; fg:from-theme"`
		ClosedIndicator   string `themedef:"bg:from-theme; fg:from-theme"`
		Selection         string `themedef:"bg:from-theme; fg:blue; attr:reverse"`

		Adapter struct {
			Present string `themedef:"bg:from-theme; fg:from-theme"`
		}

		Device struct {
			Discovered string `themedef:"bg:from-theme; fg:from-theme"`
			Paired     string `themedef:"bg:from-theme; fg:from-theme"`
		}

		DevicesList struct {
			Nodes string `themedef:"bg:from-theme; fg:from-theme"`
		}

		ActionsList struct {
			Nodes string `themedef:"bg:from-theme; fg:from-theme"`
		}
	}

	TabsPane struct {
		Style      string `themedef:"bg:from-theme; fg:from-theme"`
		Tab        string `themedef:"bg:from-theme; fg:from-theme"`
		FocusedTab string `themedef:"bg:from-theme; fg:brightcyan"`
	}

	Info struct {
		Style   string `themedef:"bg:blue; fg:white"`
		Heading string `themedef:"bg:from-theme; fg:from-theme"`
	}

	Operations struct {
		Style   string `themedef:"bg:blue; fg:white"`
		Heading string `themedef:"bg:from-theme; fg:from-theme"`
	}

	Log struct {
		Style   string `themedef:"bg:from-theme; fg:from-theme"`
		Heading string `themedef:"bg:from-theme; fg:from-theme"`

		Info  string `themedef:"bg:from-theme; fg:blue; attr:bold,underline"`
		Debug string `themedef:"bg:from-theme; fg:brightpurple; attr:bold,underline"`
		Error string `themedef:"bg:from-theme; fg:red; attr:bold,underline"`
	}

	StatusBar struct {
		Style string `themedef:"bg:purple;fg:white;attr:bold"`
	}
}

// RootConfiguration represents the root configuration.
type RootConfiguration Configuration

// Merge merges the provided configuration with the root configuration.
func (r *RootConfiguration) Merge(cfg *Configuration) error {
	tintName := cfg.Tint

	if t := tint.DefaultTintsByID(tintName); t == nil {
		return fmt.Errorf("the specified theme was not found: %s", tintName)
	}

	r.Tint = tintName

	for parseInfo := range iterProperties(r, cfg, nil) {
		rootCfg := *parseInfo.rootCfg

		cmpParseInfo, parseErr := r.checkItem(parseInfo.cmpCfg)
		if parseErr != nil {
			return parseErr
		}

		rootParseInfo, parseErr := r.checkItem(rootCfg)
		if parseErr != nil {
			return parseErr
		}

		*parseInfo.rootCfg = rootParseInfo.merge(cmpParseInfo)
	}

	return nil
}

// convertToTheme converts the configuration to a [Theme].
func (r *RootConfiguration) convertToTheme() Theme {
	tintSpec := tint.DefaultTintsByID(r.Tint)
	tint.Register(tintSpec)
	tint.SetTint(tintSpec)

	th := Theme{}

	for parseInfo := range iterProperties(r, nil, &th) {
		r.applyTheme(parseInfo)
	}

	return th
}

func (r *RootConfiguration) applyTheme(p parseThemeInfo) {
	style := lipgloss.NewStyle()

	for seg := range strings.SplitSeq(*p.rootCfg, propSegmentSep) {
		prop, val, _ := strings.Cut(seg, propValSep)

		switch prop {
		case fgKeyword:
			switch val {
			case fromThemeKeyword:
				style = style.Foreground(tint.Current().Fg)

			default:
				color, _ := parseColor(val)
				style = style.Foreground(color)
			}

		case bgKeyword:
			switch val {
			case fromThemeKeyword:
				style = style.Background(tint.Current().Bg)

			default:
				color, _ := parseColor(val)
				style = style.Background(color)
			}

		case attrKeyword:
			for attr := range strings.SplitSeq(val, multiValSep) {
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
				}
			}
		}
	}

	*p.style = style
}

func (r *RootConfiguration) checkItem(cfgItem string) (cfgParseInfo parseConfigInfo, err error) {
	if cfgItem == "" {
		return cfgParseInfo, nil
	}

	for seg := range strings.SplitSeq(removeSpaces(cfgItem), propSegmentSep) {
		prop, val, ok := strings.Cut(seg, propValSep)
		if !ok {
			return cfgParseInfo, fmt.Errorf("invalid property sequence: %s -> %s (%s)", prop, val, cfgItem)
		}

		if val == "" {
			return cfgParseInfo, fmt.Errorf("no value was specified for property %s (%s)", prop, cfgItem)
		}

		switch prop {
		case fgKeyword, bgKeyword:
			switch val {
			case fromThemeKeyword:

			default:
				if _, ok := parseColor(val); !ok {
					return cfgParseInfo, fmt.Errorf("invalid color specifier %s for property %s (%s)", val, prop, cfgItem)
				}
			}

			if prop == fgKeyword {
				cfgParseInfo.fg = val
				continue
			}

			cfgParseInfo.bg = val

		case attrKeyword:
			const (
				totalAttrs = 5
			)

			var attrCount int

			for attr := range strings.SplitSeq(val, multiValSep) {
				if attr == "" {
					continue
				}

				switch attr {
				case "bold", "italic", "underline", "blink", "reverse":
					attrCount++
					if attrCount > totalAttrs {
						return cfgParseInfo, fmt.Errorf("invalid property length of %s in %s (%s)", val, prop, cfgItem)
					}

					continue

				default:
					return cfgParseInfo, fmt.Errorf("invalid attribute %s specified for property %s (%s)", attr, prop, cfgItem)
				}
			}

			cfgParseInfo.attrs = val

		default:
			return cfgParseInfo, fmt.Errorf("property is not supported: %s (%s)", prop, cfgItem)
		}
	}

	return cfgParseInfo, nil
}

type parseConfigInfo struct {
	fg, bg string
	attrs  string
}

func (p *parseConfigInfo) merge(cmpCfg parseConfigInfo) string {
	const (
		numSemicolons = 2
	)

	if cmpCfg.fg != "" {
		p.fg = cmpCfg.fg
	}

	if cmpCfg.bg != "" {
		p.bg = cmpCfg.bg
	}

	if cmpCfg.attrs != "" {
		p.attrs = cmpCfg.attrs
	}

	fgPrefixLen, bgPrefixLen, attrPrefixLen := len(fgKeyword), len(bgKeyword), len(attrKeyword)

	var sb strings.Builder

	sb.Grow(
		(min(len(p.fg), 10) + fgPrefixLen) +
			(min(len(p.bg), 10) + bgPrefixLen) +
			(min(len(p.attrs), 30) + attrPrefixLen) +
			numSemicolons,
	)

	if p.fg != "" {
		sb.WriteString(fgKeyword)
		sb.WriteString(propValSep)
		sb.WriteString(p.fg)
	}

	if sb.Len() > 0 {
		sb.WriteString(propSegmentSep)
	}

	if p.bg != "" {
		sb.WriteString(bgKeyword)
		sb.WriteString(propValSep)
		sb.WriteString(p.bg)
	}

	if sb.Len() > 0 {
		sb.WriteString(propSegmentSep)
	}

	if p.attrs != "" {
		sb.WriteString(attrKeyword)
		sb.WriteString(propValSep)
		sb.WriteString(p.attrs)
	}

	return sb.String()
}

func removeSpaces(s string) string {
	if !strings.ContainsFunc(s, unicode.IsSpace) {
		return s
	}

	var sb strings.Builder

	sb.Grow(max(minBufferLen, min(len(s), maxBuffenLen)))

	for field := range strings.FieldsSeq(s) {
		sb.WriteString(field)
	}

	return sb.String()
}
