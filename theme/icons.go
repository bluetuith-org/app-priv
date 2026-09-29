package theme

// IconVariant holds both Unicode and ASCII representations of an icon.
type IconVariant struct {
	IsASCII        bool
	Unicode, ASCII string
}

// String returns the Unicode or ASCII icon depending on whether IsASCII is true.
func (i IconVariant) String() string {
	if i.IsASCII {
		return i.ASCII
	}

	return i.Unicode
}

// IconSet represents a set of icons.
type IconSet struct {
	Adapters, Adapter IconVariant
	Devices           IconVariant
	Actions           IconVariant

	Camera, VideoCamera     IconVariant
	Printer                 IconVariant
	Speaker, Headphones     IconVariant
	Mouse, Tablet, Keyboard IconVariant
	Gaming                  IconVariant
	Network                 IconVariant
	Computer, Phone         IconVariant
	Uncategorized           IconVariant

	Info       IconVariant
	Operations IconVariant
	Log        IconVariant

	ArrowLeft     IconVariant
	ArrowRight    IconVariant
	TriangleRight IconVariant
	TriangleDown  IconVariant
}

// NewIconSet returns a new icon set.
func NewIconSet(isASCII bool) *IconSet {
	ic := &IconSet{}

	return ic.populateIcons(isASCII)
}

func (i *IconSet) populateIcons(isASCII bool) *IconSet {
	ic := &IconSet{
		Info: IconVariant{
			Unicode: "\U0001F4C4",
			ASCII:   "[I]",
			IsASCII: isASCII,
		},
		ArrowLeft: IconVariant{
			Unicode: "\u2190",
			ASCII:   "<-",
			IsASCII: isASCII,
		},
		ArrowRight: IconVariant{
			Unicode: "\u2192",
			ASCII:   "->",
			IsASCII: isASCII,
		},
		TriangleRight: IconVariant{
			Unicode: "\u25b6",
			ASCII:   ">",
			IsASCII: isASCII,
		},
		TriangleDown: IconVariant{
			Unicode: "\u25bc",
			ASCII:   "v",
			IsASCII: isASCII,
		},
		Operations: IconVariant{
			Unicode: "\u2699",
			ASCII:   "[O]",
		},
		Log: IconVariant{
			Unicode: "\U0001F50D",
			ASCII:   "[L]",
			IsASCII: isASCII,
		},
		Adapters: IconVariant{
			Unicode: "",
			ASCII:   "",
			IsASCII: isASCII,
		},
		Adapter: IconVariant{
			Unicode: "",
			ASCII:   "[Adapter]",
			IsASCII: isASCII,
		},
		Devices: IconVariant{
			Unicode: "\U0001F4F1\U0001F3A7",
			ASCII:   "",
			IsASCII: isASCII,
		},
		Actions: IconVariant{
			Unicode: "\U00002794",
			ASCII:   "[A]",
			IsASCII: isASCII,
		},
		Camera: IconVariant{
			IsASCII: isASCII,
			Unicode: "\U0001F4F7",
			ASCII:   "[Camera]",
		},
		VideoCamera: IconVariant{
			IsASCII: isASCII,
			Unicode: "\U0001F4F9",
			ASCII:   "[VideoCamera]",
		},
		Printer: IconVariant{
			IsASCII: isASCII,
			Unicode: "\U0001F5A8",
			ASCII:   "[Printer]",
		},
		Speaker: IconVariant{
			IsASCII: isASCII,
			Unicode: "\U0001F50A",
			ASCII:   "[Speaker]",
		},
		Headphones: IconVariant{
			IsASCII: isASCII,
			Unicode: "\U0001F3A7",
			ASCII:   "[Headphones]",
		},
		Mouse: IconVariant{
			IsASCII: isASCII,
			Unicode: "\U00001F5B1",
			ASCII:   "[Mouse]",
		},
		Tablet: IconVariant{
			IsASCII: isASCII,
			Unicode: "\U0001F524",
			ASCII:   "[Tablet]",
		},
		Keyboard: IconVariant{
			IsASCII: isASCII,
			Unicode: "\U00002328",
			ASCII:   "[Keyboard]",
		},
		Gaming: IconVariant{
			IsASCII: isASCII,
			Unicode: "\U0001F3AE",
			ASCII:   "[Gaming]",
		},
		Network: IconVariant{
			IsASCII: isASCII,
			Unicode: "\U0001F310",
			ASCII:   "[Network]",
		},
		Computer: IconVariant{
			IsASCII: isASCII,
			Unicode: "\U0001F4BB",
			ASCII:   "[Computer]",
		},
		Phone: IconVariant{
			IsASCII: isASCII,
			Unicode: "\U0001F4F1",
			ASCII:   "[Phone]",
		},
		Uncategorized: IconVariant{
			IsASCII: isASCII,
			Unicode: "\U00002753",
			ASCII:   "[Unknown]",
		},
	}

	return ic
}

// DeviceClassToIcon converts a device's COD to an icon.
//
// Taken from: https://github.com/bluez/bluez/blob/master/src/dbus-common.c
func DeviceClassToIcon(i *IconSet, class uint32) IconVariant {
	switch (class & 0x1f00) >> 8 {
	case 0x01:
		return i.Computer

	case 0x02:
		switch (class & 0xfc) >> 2 {
		case 0x01, 0x02, 0x03, 0x05:
			return i.Phone

		case 0x04:
			return i.Network
		}

	case 0x03:
		return i.Network

	case 0x04:
		switch (class & 0xfc) >> 2 {
		case 0x01, 0x02:
			return i.Headphones

		case 0x06:
			return i.Headphones

		case 0x0b, 0x0c, 0x0d:
			return i.VideoCamera

		default:
			return i.Speaker
		}

	case 0x05:
		switch (class & 0xc0) >> 6 {
		case 0x00:
			switch (class & 0x1e) >> 2 {
			case 0x01, 0x02:
				return i.Gaming
			}

		case 0x01:
			return i.Keyboard

		case 0x02:
			switch (class & 0x1e) >> 2 {
			case 0x05:
				return i.Tablet

			default:
				return i.Mouse
			}
		}

	case 0x06:
		if (class & 0x80) != 0 {
			return i.Printer
		}
		if (class & 0x20) != 0 {
			return i.Camera
		}
	}

	return i.Uncategorized
}
