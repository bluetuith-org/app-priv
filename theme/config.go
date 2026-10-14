package theme

// Configuration represents the app's theme configuration.
type Configuration struct {
	Tint string `themetype:"primer"`

	Global string `themedef:"bg:from-theme; fg:from-theme"`

	App struct {
		TitleBar string `themedef:"bg:purple; fg:white; attr:bold"`

		Border        string `themedef:"bg:from-theme; fg:white"`
		BorderFocused string `themedef:"bg:from-theme; fg:green; attr:bold"`

		PlaceHolder     string `themedef:"bg:from-theme; fg:grey; attr:italic"`
		CurrentViewName string `themedef:"bg:brightblack; fg:from-theme; attr:italic"`
	}

	ADTree struct {
		TitleBarIndicator string `themedef:"bg:#5f5ffd; fg:white; attr:bold"`
		Indicators        string `themedef:"bg:from-theme; fg:from-theme"`
		Selection         string `themedef:"bg:from-theme; fg:blue; attr:reverse"`

		Adapter struct {
			Name string `themedef:"bg:from-theme; fg:from-theme"`

			PoweredOnPill    string `themedef:"bg:green; fg:black"`
			PoweredOffPill   string `themedef:"bg:red; fg:white"`
			DiscoverablePill string `themedef:"bg:blue; fg:white"`
			ScanningPill     string `themedef:"bg:yellow; fg:black"`
			PairablePill     string `themedef:"bg:cyan; fg:black"`
		}

		Device struct {
			Discovered string `themedef:"bg:from-theme; fg:orange"`
			Paired     string `themedef:"bg:from-theme; fg:cyan"`
			Unknown    string `themedef:"bg:from-theme; fg:white"`

			NewDevicePill string `themedef:"bg:orange; fg:black"`
			PairedPill    string `themedef:"bg:cyan; fg:black"`
			TrustedPill   string `themedef:"bg:purple; fg:white"`
			ConnectedPill string `themedef:"bg:green; fg:black"`
			BlockedPill   string `themedef:"bg:red; fg:white"`
		}

		DevicesList string `themedef:"bg:from-theme; fg:from-theme"`

		ActionsList struct {
			Style string `themedef:"bg:from-theme; fg:from-theme"`
			Nodes string `themedef:"bg:from-theme; fg:from-theme"`
		}

		Status string `themedef:"bg:from-theme; fg:from-theme"`
	}

	TabsPane struct {
		ArrowStyle string `themedef:"bg:from-theme; fg:from-theme"`
		Tab        string `themedef:"bg:from-theme; fg:from-theme"`
		FocusedTab string `themedef:"bg:from-theme; fg:brightcyan; attr:bold,underline"`
	}

	Info struct {
		Title         string `themedef:"bg:purple; fg:white"`
		PropertyName  string `themedef:"bg:from-theme; fg:brightpurple; attr:bold"`
		PropertyValue string `themedef:"bg:from-theme; fg:white"`
	}

	Operations struct {
		TitleBarIndicator string `themedef:"bg:yellow; fg:black"`
		Selection         string `themedef:"bg:from-theme; fg:blue"`

		OperationNumber string `themedef:"bg:from-theme; fg:brightpurple; attr:bold,italic,underline"`
		Time            string `themedef:"bg:from-theme; fg:grey; attr:italic"`
		PropertyName    string `themedef:"bg:from-theme; fg:brightpurple; attr:bold,underline"`
		Error           string `themedef:"bg:from-theme; fg:red; attr:bold"`
		Completed       string `themedef:"bg:from-theme; fg:green"`
		InProgress      string `themedef:"bg:from-theme; fg:yellow"`
		Message         string `themedef:"bg:from-theme; fg:white; attr:bold"`
	}

	Log struct {
		FollowIndicator string `themedef:"bg:from-theme; fg:grey; attr:italic"`
		Selection       string `themedef:"bg:from-theme; fg:blue"`

		Time    string `themedef:"bg:from-theme; fg:grey"`
		Info    string `themedef:"bg:from-theme; fg:blue(lighten:0.2); attr:bold"`
		Debug   string `themedef:"bg:from-theme; fg:brightpurple; attr:bold"`
		Error   string `themedef:"bg:from-theme; fg:red; attr:bold,underline"`
		Message string `themedef:"bg:from-theme; fg:from-theme"`
		Attr    string `themedef:"bg:from-theme; fg:from-theme; attr:bold"`
	}

	StatusBar struct {
		Style string `themedef:"bg:purple; fg:white; attr:bold"`

		Info  string `themedef:"bg:blue(lighten:0.2); fg:black; attr:bold"`
		Debug string `themedef:"bg:brightpurple; fg:black; attr:bold"`
		Error string `themedef:"bg:red; fg:white; attr:bold"`
	}
}
