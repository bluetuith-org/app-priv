package keybindings

//go:generate go run ../cmd/structgen keybindings

// Configuration holds the configuration for the keybindings of the application.
type Configuration struct {
	SelectItem  string `keydef:"enter"`
	CloseItem   string `keydef:"esc"`
	FilterItems string `keydef:"/"`

	Help    string `keydef:"?"`
	Suspend string `keydef:"ctrl+z"`
	Quit    string `keydef:"q"`

	NavigateUp     string `keydef:"up"`
	NavigateDown   string `keydef:"down"`
	NavigateLeft   string `keydef:"left"`
	NavigateRight  string `keydef:"right"`
	NavigateTop    string `keydef:"pgup"`
	NavigateBottom string `keydef:"pgdn"`

	Adapter struct {
		TogglePower        string `keydef:"o"`
		ToggleDiscoverable string `keydef:"S"`
		TogglePairable     string `keydef:"P"`
		ToggleScan         string `keydef:"s"`
	}

	Device struct {
		ToggleConnection  string `keydef:"c"`
		TogglePairedState string `keydef:"p"`
		Trust             string `keydef:"t"`

		SendFiles      string `keydef:"f"`
		NetworkOptions string `keydef:"n"`
		AudioProfiles  string `keydef:"A"`
		Block          string `keydef:"b"`
	}

	Filebrowser struct {
		CdForward string `keydef:"right"`
		CdBack    string `keydef:"left"`

		SelectOne        string `keydef:"space"`
		SelectAll        string `keydef:"A"`
		InvertSelection  string `keydef:"a"`
		ConfirmSelection string `keydef:"ctrl+s"`

		Refresh           string `keydef:"ctrl+r"`
		ToggleHiddenFiles string `keydef:"."`
	}

	ADTree struct {
		ExpandOrSelect string `keydef:"right"`
		TabPane        string `keydef:"tab"`
	}

	TabPane struct {
		Next     string `keydef:"tab"`
		Previous string `keydef:"shift+tab"`
		ADTree   string `keydef:"esc"`
	}

	Information struct {
		Show string `keydef:"I"`
	}

	Operations struct {
		Show     string `keydef:"O"`
		Cancel   string `keydef:"x"`
		ClearAll string `keydef:"C"`
	}

	Transfers struct {
		Show    string `keydef:"T"`
		Suspend string `keydef:"s"`
		Resume  string `keydef:"r"`
		Cancel  string `keydef:"x"`
	}

	Log struct {
		Show   string `keydef:"L"`
		Follow string `keydef:"f"`
	}

	Player struct {
		ToggleDisplay      string `keydef:"M"`
		ToggleMediaPlaying string `keydef:"space"`

		Next         string `keydef:">"`
		Previous     string `keydef:"<"`
		SeekForward  string `keydef:"right"`
		SeekBackward string `keydef:"left"`
		Stop         string `keydef:"]"`
	}
}
