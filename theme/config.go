package theme

// Configuration represents the app's theme configuration.
type Configuration struct {
	Tint string `themetype:"primer"`

	Global string `themedef:"-"`

	App struct {
		TitleBar string `themedef:"bg:purple; attr:bold"`

		Border        string `themedef:"-"`
		BorderFocused string `themedef:"fg:green; attr:bold"`

		PlaceHolder     string `themedef:"fg:grey; attr:italic"`
		CurrentViewName string `themedef:"bg:brightblack; attr:italic"`
	}

	ADTree struct {
		TitleBarIndicator string `themedef:"bg:#5f5ffd; attr:bold"`
		Indicators        string `themedef:"-"`
		Selection         string `themedef:"bg:blue; attr:bold"`

		Adapter struct {
			Name string `themedef:"-"`

			PoweredOnPill    string `themedef:"bg:green; fg:black"`
			PoweredOffPill   string `themedef:"bg:red"`
			DiscoverablePill string `themedef:"bg:blue"`
			ScanningPill     string `themedef:"bg:yellow; fg:black"`
			PairablePill     string `themedef:"bg:cyan; fg:black"`
		}

		Device struct {
			Discovered string `themedef:"fg:orange"`
			Paired     string `themedef:"fg:cyan"`
			Unknown    string `themedef:"-"`

			NewDevicePill string `themedef:"bg:orange; fg:black"`
			PairedPill    string `themedef:"bg:cyan; fg:black"`
			TrustedPill   string `themedef:"bg:purple"`
			ConnectedPill string `themedef:"bg:green; fg:black"`
			BlockedPill   string `themedef:"bg:red"`
		}

		DevicesList string `themedef:"-"`

		ActionsList struct {
			Style string `themedef:"-"`
			Nodes string `themedef:"-"`
		}

		Status string `themedef:"-"`
	}

	TabsPane struct {
		ArrowStyle string `themedef:"-"`
		Tab        string `themedef:"-"`
		FocusedTab string `themedef:"fg:brightcyan; attr:bold,underline"`
	}

	Info struct {
		Title         string `themedef:"bg:purple"`
		PropertyName  string `themedef:"fg:brightpurple; attr:bold"`
		PropertyValue string `themedef:"-"`
	}

	Operations struct {
		TitleBarIndicator string `themedef:"bg:yellow; fg:black"`
		Selection         string `themedef:"fg:blue; attr:bold"`

		OperationNumber string `themedef:"fg:brightpurple; attr:bold,italic,underline"`
		Time            string `themedef:"fg:grey; attr:italic"`
		PropertyName    string `themedef:"fg:brightpurple; attr:bold,underline"`
		Error           string `themedef:"fg:red; attr:bold"`
		Completed       string `themedef:"fg:green"`
		InProgress      string `themedef:"fg:yellow"`
		Message         string `themedef:"attr:bold"`
	}

	Log struct {
		FollowIndicator string `themedef:"fg:grey; attr:italic"`
		Selection       string `themedef:"fg:blue"`

		Time    string `themedef:"fg:grey"`
		Info    string `themedef:"fg:blue(lighten:0.2); attr:bold"`
		Debug   string `themedef:"fg:brightpurple; attr:bold"`
		Error   string `themedef:"fg:red; attr:bold,underline"`
		Message string `themedef:"-"`
		Attr    string `themedef:"attr:bold"`
	}

	StatusBar struct {
		Style string `themedef:"bg:purple; attr:bold"`

		Info  string `themedef:"bg:blue(lighten:0.2); fg:black; attr:bold"`
		Debug string `themedef:"bg:brightpurple; fg:black; attr:bold"`
		Error string `themedef:"bg:red; attr:bold"`
	}
}
