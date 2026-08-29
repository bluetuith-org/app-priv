package theme

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	tint "github.com/lrstanley/bubbletint/v2"
)

const (
	fgKeyword        = "fg"
	bgKeyword        = "bg"
	attrKeyword      = "attr"
	fromThemeKeyword = "from-theme"
	globalKeyword    = "global"

	propSegmentSep = ";"
	propValSep     = ":"
	multiValSep    = ","
)

//go:generate go run ../../cmd/structgen theme

// Configuration represents the app's theme configuration.
type Configuration struct {
	Tint string

	Global                string
	Border, BorderFocused string

	ADTree struct {
		Headers string

		ExpandedIndicator, ClosedIndicator string
		Selection                          string

		Adapter struct {
			Present string
		}

		Device struct {
			Discovered string
			Paired     string
		}

		DevicesList struct {
			Nodes string
		}

		ActionsList struct {
			Nodes string
		}
	}

	TabsPane struct {
		Bg              string
		Tab, FocusedTab string
	}

	Info struct {
		Bg      string
		Heading string
	}

	Operations struct {
		Bg      string
		Heading string
	}

	Log struct {
		Bg      string
		Heading string
	}

	StatusBar struct {
		Bg string
	}
}

// RootConfiguration represents the root configuration.
type RootConfiguration Configuration

// Merge merges the provided configuration with the root configuration.
func (r *RootConfiguration) Merge(cfg *Configuration) error {
	var sb strings.Builder

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

		*parseInfo.rootCfg = rootParseInfo.format(&sb, cmpParseInfo)
	}

	return nil
}

// convertToTheme converts the configuration to a [Theme].
func (r *RootConfiguration) convertToTheme() Theme {
	th := Theme{}

	for parseInfo := range iterProperties(r, nil, &th) {
		r.applyTheme(parseInfo)
	}

	return th
}

func (r *RootConfiguration) applyTheme(p *parseThemeInfo) {
	style := lipgloss.NewStyle()

	for seg := range strings.SplitSeq(*p.rootCfg, propSegmentSep) {
		prop, val, _ := strings.Cut(seg, propValSep)

		switch prop {
		case fgKeyword:
			switch val {
			case fromThemeKeyword:
				style = style.Foreground(tint.Current().Fg)

			case globalKeyword:

			default:
				color, _ := parseColor(val)
				style = style.Foreground(color)
			}

		case bgKeyword:
			switch val {
			case fromThemeKeyword:
				style = style.Background(tint.Current().Bg)

			case globalKeyword:

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
}

func (r *RootConfiguration) checkItem(cfgItem string) (cfgParseInfo parseConfigInfo, err error) {
	if cfgItem == "" {
		return cfgParseInfo, nil
	}

	for seg := range strings.SplitSeq(strings.TrimSpace(cfgItem), propSegmentSep) {
		prop, val, ok := strings.Cut(seg, propValSep)
		if !ok {
			return cfgParseInfo, fmt.Errorf("invalid property sequence: %s -> %s (%s)", prop, val, cfgItem)
		}

		if val == "" {
			return cfgParseInfo, fmt.Errorf("no value was specified for property %s (%s)", prop, cfgItem)
		}

		switch prop {
		case fgKeyword:
			switch val {
			case fromThemeKeyword:
			case globalKeyword:

			default:
				if _, ok := parseColor(val); !ok {
					return cfgParseInfo, fmt.Errorf("invalid color specifier %s for property %s (%s)", val, prop, cfgItem)
				}
			}

			cfgParseInfo.fg = val

		case bgKeyword:
			switch val {
			case fromThemeKeyword:
			case globalKeyword:

			default:
				if _, ok := parseColor(val); !ok {
					return cfgParseInfo, fmt.Errorf("invalid color specifier %s for property %s (%s)", val, prop, cfgItem)
				}
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

func (p *parseConfigInfo) format(sb *strings.Builder, cmpCfg parseConfigInfo) string {
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

	sb.Reset()
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
