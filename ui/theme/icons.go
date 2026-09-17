package theme

// IconVariant holds both Unicode and ASCII representations of an icon.
type IconVariant struct {
	IsASCII        bool
	Unicode, ASCII string
}

// GetIcon returns the Unicode or ASCII icon depending on whether IsASCII is true.
func (i *IconVariant) GetIcon() string {
	if i.IsASCII {
		return i.ASCII
	}

	return i.Unicode
}

// IconSet represents a set of icons.
type IconSet struct {
	Info          *IconVariant
	ArrowLeft     *IconVariant
	ArrowRight    *IconVariant
	TriangleRight *IconVariant
	TriangleDown  *IconVariant
	CogWheel      *IconVariant
	Log           *IconVariant
}

// NewIconSet returns a new icon set.
func NewIconSet(isASCII bool) *IconSet {
	ic := &IconSet{}

	return ic.populateIcons(isASCII)
}

func (i *IconSet) populateIcons(isASCII bool) *IconSet {
	ic := &IconSet{
		Info: &IconVariant{
			Unicode: "\U0001F4C4",
			ASCII:   "[I]",
			IsASCII: isASCII,
		},
		ArrowLeft: &IconVariant{
			Unicode: "\u2190",
			ASCII:   "<-",
			IsASCII: isASCII,
		},
		ArrowRight: &IconVariant{
			Unicode: "\u2192",
			ASCII:   "->",
			IsASCII: isASCII,
		},
		TriangleRight: &IconVariant{
			Unicode: "\u25b6",
			ASCII:   ">",
			IsASCII: isASCII,
		},
		TriangleDown: &IconVariant{
			Unicode: "\u25bc",
			ASCII:   "v",
			IsASCII: isASCII,
		},
		CogWheel: &IconVariant{
			Unicode: "\u2699",
			ASCII:   "[-O-]",
		},
		Log: &IconVariant{
			Unicode: "\U0001F50D",
			ASCII:   "",
			IsASCII: isASCII,
		},
	}

	return ic
}
