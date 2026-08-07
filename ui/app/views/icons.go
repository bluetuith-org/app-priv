package views

import "charm.land/lipgloss/v2"

type iconSet struct {
	Info       *iconVariant
	ArrowLeft  *iconVariant
	ArrowRight *iconVariant
}

func newIconSet(isASCII bool) *iconSet {
	ic := &iconSet{}

	return ic.populateIcons(isASCII)
}

func (i *iconSet) populateIcons(isASCII bool) *iconSet {
	ic := &iconSet{
		Info: &iconVariant{
			Unicode: "\U0001F4C4",
			ASCII:   "[I]",
			IsASCII: isASCII,
		},
		ArrowLeft: &iconVariant{
			Unicode: "\u2190",
			ASCII:   "<-",
			IsASCII: isASCII,
		},
		ArrowRight: &iconVariant{
			Unicode: "\u2192",
			ASCII:   "->",
			IsASCII: isASCII,
		},
	}

	return ic
}

type iconVariant struct {
	IsASCII        bool
	Unicode, ASCII string
}

func (i *iconVariant) Render(style lipgloss.Style, text string) string {
	return style.Render(i.getIcon(), text)
}

func (i *iconVariant) getIcon() string {
	if i.IsASCII {
		return i.ASCII
	}

	return i.Unicode
}
