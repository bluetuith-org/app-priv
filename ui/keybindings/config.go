package keybindings

//go:generate go run ../../cmd/structgen keybindings

// Configuration holds the configuration for the keybindings of the application.
type Configuration struct {
	SwitchPanes string `keydef:"tab" keyshorthelp:"SwitchPanes" keylonghelp:"SwitchPanes"`

	SelectItem  string `keydef:"enter" keyshorthelp:"SelectItem" keylonghelp:"SelectItem"`
	CloseItem   string `keydef:"esc" keyshorthelp:"CloseItem" keylonghelp:"CloseItem"`
	FilterItems string `keydef:"/" keyshorthelp:"FilterItems" keylonghelp:"FilterItems"`

	Help    string `keydef:"?" keyshorthelp:"Help" keylonghelp:"Help"`
	Suspend string `keydef:"ctrl+z" keyshorthelp:"Suspend" keylonghelp:"Suspend"`
	Quit    string `keydef:"q" keyshorthelp:"Quit" keylonghelp:"Quit"`

	NavigateUp     string `keydef:"up" keyshorthelp:"NavigateUp" keylonghelp:"NavigateUp"`
	NavigateDown   string `keydef:"down" keyshorthelp:"NavigateDown" keylonghelp:"NavigateDown"`
	NavigateLeft   string `keydef:"left" keyshorthelp:"NavigateLeft" keylonghelp:"NavigateLeft"`
	NavigateRight  string `keydef:"right" keyshorthelp:"NavigateRight" keylonghelp:"NavigateRight"`
	NavigateTop    string `keydef:"pgup" keyshorthelp:"NavigateTop" keylonghelp:"NavigateTop"`
	NavigateBottom string `keydef:"pgdown" keyshorthelp:"NavigateBottom" keylonghelp:"NavigateBottom"`

	Adapter struct {
		TogglePower        string `keydef:"o" keyshorthelp:"TogglePower" keylonghelp:"TogglePower"`
		ToggleDiscoverable string `keydef:"S" keyshorthelp:"ToggleDiscoverable" keylonghelp:"ToggleDiscoverable"`
		TogglePairable     string `keydef:"P" keyshorthelp:"TogglePairable" keylonghelp:"TogglePairable"`
		ToggleScan         string `keydef:"s" keyshorthelp:"ToggleScan" keylonghelp:"ToggleScan"`
	}

	Device struct {
		ToggleConnection  string `keydef:"c" keyshorthelp:"ToggleConnection" keylonghelp:"ToggleConnection"`
		TogglePairedState string `keydef:"p" keyshorthelp:"TogglePairedState" keylonghelp:"TogglePairedState"`
		Trust             string `keydef:"t" keyshorthelp:"Trust" keylonghelp:"Trust"`

		SendFiles      string `keydef:"f" keyshorthelp:"SendFiles" keylonghelp:"SendFiles"`
		NetworkOptions string `keydef:"n" keyshorthelp:"NetworkOptions" keylonghelp:"NetworkOptions"`
		AudioProfiles  string `keydef:"A" keyshorthelp:"AudioProfiles" keylonghelp:"AudioProfiles"`
		Block          string `keydef:"b" keyshorthelp:"Block" keylonghelp:"Block"`
	}

	Filebrowser struct {
		CdForward string `keydef:"right" keyshorthelp:"CdForward" keylonghelp:"CdForward"`
		CdBack    string `keydef:"left" keyshorthelp:"CdBack" keylonghelp:"CdBack"`

		SelectOne        string `keydef:"space" keyshorthelp:"SelectOne" keylonghelp:"SelectOne"`
		SelectAll        string `keydef:"A" keyshorthelp:"SelectAll" keylonghelp:"SelectAll"`
		InvertSelection  string `keydef:"a" keyshorthelp:"InvertSelection" keylonghelp:"InvertSelection"`
		ConfirmSelection string `keydef:"ctrl+s" keyshorthelp:"ConfirmSelection" keylonghelp:"ConfirmSelection"`

		Refresh           string `keydef:"ctrl+r" keyshorthelp:"Refresh" keylonghelp:"Refresh"`
		ToggleHiddenFiles string `keydef:"." keyshorthelp:"ToggleHiddenFiles" keylonghelp:"ToggleHiddenFiles"`
	}

	ADTree struct {
		ToggleNodes string `keydef:"right" keyshorthelp:"ToggleNodes" keylonghelp:"ToggleNodes"`
	}

	Operations struct {
		Cancel string `keydef:"x" keyshorthelp:"Cancel" keylonghelp:"Cancel"`
	}

	Transfers struct {
		Suspend string `keydef:"s" keyshorthelp:"Suspend" keylonghelp:"Suspend"`
		Resume  string `keydef:"r" keyshorthelp:"Resume" keylonghelp:"Resume"`
		Cancel  string `keydef:"x" keyshorthelp:"Cancel" keylonghelp:"Cancel"`
	}

	Player struct {
		ToggleDisplay      string `keydef:"M" keyshorthelp:"ToggleDisplay" keylonghelp:"ToggleDisplay"`
		ToggleMediaPlaying string `keydef:"space" keyshorthelp:"ToggleMediaPlaying" keylonghelp:"ToggleMediaPlaying"`

		Next         string `keydef:">" keyshorthelp:"Next" keylonghelp:"Next"`
		Previous     string `keydef:"<" keyshorthelp:"Previous" keylonghelp:"Previous"`
		SeekForward  string `keydef:"right" keyshorthelp:"SeekForward" keylonghelp:"SeekForward"`
		SeekBackward string `keydef:"left" keyshorthelp:"SeekBackward" keylonghelp:"SeekBackward"`
		Stop         string `keydef:"]" keyshorthelp:"Stop" keylonghelp:"Stop"`
	}
}
