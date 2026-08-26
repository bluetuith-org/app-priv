package theme

//go:generate go run ../../cmd/themegen

// Configuration represents the app's theme configuration.
// "bg:<color>, fg:<color>, attr:<attr1>, <attr2>..."
// bg, fg additional props: default, global
type Configuration struct {
	Bg, Fg                string
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
