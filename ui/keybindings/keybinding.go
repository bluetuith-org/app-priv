package keybindings

import (
	"iter"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// KeyID represents a keybinding's ID.
type KeyID int

// Keybinding represents a single Keybinding.
type Keybinding struct {
	ID                       KeyID
	Key, ShortHelp, LongHelp string
}

// NewKeybinding returns a new keybinding.
func NewKeybinding(id KeyID, keyCombo, shorthelp, longhelp string) Keybinding {
	return Keybinding{ID: id, Key: keyCombo, ShortHelp: shorthelp, LongHelp: longhelp}
}

// KeyBindingID describes the application keybinding type.
type KeyBindingID string

// The different application keybinding types.
const (
	KeyMenu    KeyBindingID = "Menu"
	KeySelect  KeyBindingID = "Select"
	KeySuspend KeyBindingID = "Suspend"
	KeyQuit    KeyBindingID = "Quit"
	KeySwitch  KeyBindingID = "Switch"
	KeyClose   KeyBindingID = "Close"
	KeyHelp    KeyBindingID = "Help"
	KeyFilter  KeyBindingID = "Filter"

	KeyAdapterChange             KeyBindingID = "AdapterChange"
	KeyAdapterTogglePower        KeyBindingID = "AdapterTogglePower"
	KeyAdapterToggleDiscoverable KeyBindingID = "AdapterToggleDiscoverable"
	KeyAdapterTogglePairable     KeyBindingID = "AdapterTogglePairable"
	KeyAdapterToggleScan         KeyBindingID = "AdapterToggleScan"

	KeyDeviceSendFiles     KeyBindingID = "DeviceSendFiles"
	KeyDeviceNetwork       KeyBindingID = "DeviceNetwork"
	KeyDeviceConnect       KeyBindingID = "DeviceConnect"
	KeyDevicePair          KeyBindingID = "DevicePair"
	KeyDeviceTrust         KeyBindingID = "DeviceTrust"
	KeyDeviceBlock         KeyBindingID = "DeviceBlock"
	KeyDeviceAudioProfiles KeyBindingID = "DeviceAudioProfiles"
	KeyDeviceRemove        KeyBindingID = "DeviceRemove"

	KeyFilebrowserDirForward       KeyBindingID = "FilebrowserDirForward"
	KeyFilebrowserDirBack          KeyBindingID = "FilebrowserDirBack"
	KeyFilebrowserSelect           KeyBindingID = "FilebrowserSelect"
	KeyFilebrowserInvertSelection  KeyBindingID = "FilebrowserInvertSelection"
	KeyFilebrowserSelectAll        KeyBindingID = "FilebrowserSelectAll"
	KeyFilebrowserRefresh          KeyBindingID = "FilebrowserRefresh"
	KeyFilebrowserToggleHidden     KeyBindingID = "FilebrowserToggleHidden"
	KeyFilebrowserConfirmSelection KeyBindingID = "FilebrowserConfirmSelection"

	KeyProgressView            KeyBindingID = "ProgressView"
	KeyProgressTransferSuspend KeyBindingID = "ProgressTransferSuspend"
	KeyProgressTransferResume  KeyBindingID = "ProgressTransferResume"
	KeyProgressTransferCancel  KeyBindingID = "ProgressTransferCancel"

	KeyInfo KeyBindingID = "Info"

	KeyTasks      KeyBindingID = "Tasks"
	KeyTaskCancel KeyBindingID = "TaskCancel"

	KeyPlayerShow         KeyBindingID = "PlayerShow"
	KeyPlayerHide         KeyBindingID = "PlayerHide"
	KeyPlayerTogglePlay   KeyBindingID = "PlayerTogglePlay"
	KeyPlayerNext         KeyBindingID = "PlayerNext"
	KeyPlayerPrevious     KeyBindingID = "PlayerPrevious"
	KeyPlayerSeekForward  KeyBindingID = "PlayerSeekForward"
	KeyPlayerSeekBackward KeyBindingID = "PlayerSeekBackward"
	KeyPlayerStop         KeyBindingID = "PlayerStop"

	KeyNavigateUp     KeyBindingID = "NavigateUp"
	KeyNavigateDown   KeyBindingID = "NavigateDown"
	KeyNavigateRight  KeyBindingID = "NavigateRight"
	KeyNavigateLeft   KeyBindingID = "NavigateLeft"
	KeyNavigateTop    KeyBindingID = "NavigateTop"
	KeyNavigateBottom KeyBindingID = "NavigateBottom"
)

// KeyContext describes the context where the keybinding is
// supposed to be applied in.
type KeyContext string

// The different context types for keybindings.
const (
	ContextApp      KeyContext = "App"
	ContextDevice   KeyContext = "Device"
	ContextFiles    KeyContext = "Files"
	ContextProgress KeyContext = "Progress"
	ContextTasks    KeyContext = "Tasks"
)

// IterKeyMatch describes an iterator which calls the provided function
// when a KeyID is found.
type IterKeyMatch[V any] iter.Seq2[KeyBindingID, V]

// KeybindingVar represents a custom key binding.
type KeybindingVar struct {
	ID      KeyBindingID
	Binding key.Binding
	Global  bool
}

// MatchesKey checks if a pressed key matches the stored binding setting.
func MatchesKey(k KeyBindingID, p tea.KeyPressMsg) bool {
	e, ok := _keybindings[k]
	if !ok {
		return false
	}

	return key.Matches(p, e.Binding)
}

// IterMatch iterates over a sequence of keys and finds a match.
func IterMatch[V any](p tea.KeyPressMsg, keys IterKeyMatch[V]) (KeybindingVar, V, bool) {
	var val V

	for k, v := range keys {
		e, ok := _keybindings[k]
		if !ok {
			continue
		}

		if key.Matches(p, e.Binding) {
			return e, v, true
		}
	}

	return KeybindingVar{}, val, false
}

// RawBinding returns the raw string-based keybinding setting.
func RawBinding(k KeyBindingID) []string {
	return _keybindings[k].Binding.Keys()
}

var _keybindings = map[KeyBindingID]KeybindingVar{
	KeySwitch: {
		ID:      KeySwitch,
		Global:  true,
		Binding: key.NewBinding(key.WithKeys("tab"), key.WithHelp("", "Switch")),
	},
	KeyClose: {
		ID:      KeyClose,
		Global:  true,
		Binding: key.NewBinding(key.WithKeys("esc"), key.WithHelp("", "Close")),
	},
	KeyQuit: {
		ID:      KeyQuit,
		Global:  true,
		Binding: key.NewBinding(key.WithKeys("Q"), key.WithHelp("", "Quit")),
	},
	KeyMenu: {
		ID:      KeyMenu,
		Global:  true,
		Binding: key.NewBinding(key.WithKeys("m"), key.WithHelp("", "Menu")),
	},
	KeySelect: {
		ID:      KeySelect,
		Global:  true,
		Binding: key.NewBinding(key.WithKeys("enter"), key.WithHelp("", "Select")),
	},
	KeySuspend: {
		ID:      KeySuspend,
		Global:  true,
		Binding: key.NewBinding(key.WithKeys("ctrl+z"), key.WithHelp("", "Suspend")),
	},
	KeyHelp: {
		ID:      KeyHelp,
		Global:  true,
		Binding: key.NewBinding(key.WithKeys("?"), key.WithHelp("", "Help")),
	},
	KeyFilter: {
		ID:      KeyFilter,
		Global:  true,
		Binding: key.NewBinding(key.WithKeys("/"), key.WithHelp("", "Filter")),
	},
	KeyNavigateUp: {
		ID:      KeyNavigateUp,
		Global:  true,
		Binding: key.NewBinding(key.WithKeys("up"), key.WithHelp("", "Navigate Up")),
	},
	KeyNavigateDown: {
		ID:      KeyNavigateDown,
		Global:  true,
		Binding: key.NewBinding(key.WithKeys("down"), key.WithHelp("", "Navigate Down")),
	},
	KeyNavigateRight: {
		ID:      KeyNavigateRight,
		Global:  true,
		Binding: key.NewBinding(key.WithKeys("right"), key.WithHelp("", "Navigate Right")),
	},
	KeyNavigateLeft: {
		ID:      KeyNavigateLeft,
		Global:  true,
		Binding: key.NewBinding(key.WithKeys("left"), key.WithHelp("", "Navigate Left")),
	},
	KeyNavigateTop: {
		ID:      KeyNavigateTop,
		Global:  true,
		Binding: key.NewBinding(key.WithKeys("pgup"), key.WithHelp("", "Navigate Top")),
	},
	KeyNavigateBottom: {
		ID:      KeyNavigateBottom,
		Global:  true,
		Binding: key.NewBinding(key.WithKeys("pgdown"), key.WithHelp("", "Navigate Bottom")),
	},
	KeyAdapterTogglePower: {
		ID:      KeyAdapterTogglePower,
		Binding: key.NewBinding(key.WithKeys("o"), key.WithHelp("", "Power")),
	},
	KeyAdapterToggleDiscoverable: {
		ID:      KeyAdapterToggleDiscoverable,
		Binding: key.NewBinding(key.WithKeys("S"), key.WithHelp("", "Discoverable")),
	},
	KeyAdapterTogglePairable: {
		ID:      KeyAdapterTogglePairable,
		Binding: key.NewBinding(key.WithKeys("P"), key.WithHelp("", "Pairable")),
	},
	KeyAdapterToggleScan: {
		ID:      KeyAdapterToggleScan,
		Binding: key.NewBinding(key.WithKeys("s"), key.WithHelp("", "Scan")),
	},
	KeyAdapterChange: {
		ID:      KeyAdapterChange,
		Binding: key.NewBinding(key.WithKeys("a"), key.WithHelp("", "Change")),
	},
	KeyDeviceConnect: {
		ID:      KeyDeviceConnect,
		Binding: key.NewBinding(key.WithKeys("c"), key.WithHelp("", "Connect")),
	},
	KeyDevicePair: {
		ID:      KeyDevicePair,
		Binding: key.NewBinding(key.WithKeys("p"), key.WithHelp("", "Pair")),
	},
	KeyDeviceTrust: {
		ID:      KeyDeviceTrust,
		Binding: key.NewBinding(key.WithKeys("t"), key.WithHelp("", "Trust")),
	},
	KeyDeviceBlock: {
		ID:      KeyDeviceBlock,
		Binding: key.NewBinding(key.WithKeys("b"), key.WithHelp("", "Block")),
	},
	KeyDeviceSendFiles: {
		ID:      KeyDeviceSendFiles,
		Binding: key.NewBinding(key.WithKeys("f"), key.WithHelp("", "Send")),
	},
	KeyDeviceNetwork: {
		ID:      KeyDeviceNetwork,
		Binding: key.NewBinding(key.WithKeys("n"), key.WithHelp("", "Network Options")),
	},
	KeyDeviceAudioProfiles: {
		ID:      KeyDeviceAudioProfiles,
		Binding: key.NewBinding(key.WithKeys("A"), key.WithHelp("", "Audio Profiles")),
	},
	KeyInfo: {
		ID:      KeyInfo,
		Binding: key.NewBinding(key.WithKeys("i"), key.WithHelp("", "Info")),
	},
	KeyDeviceRemove: {
		ID:      KeyDeviceRemove,
		Binding: key.NewBinding(key.WithKeys("d"), key.WithHelp("", "Remove")),
	},
	KeyPlayerShow: {
		ID:      KeyPlayerShow,
		Binding: key.NewBinding(key.WithKeys("m"), key.WithHelp("", "Show Media Player")),
	},
	KeyPlayerHide: {
		ID:      KeyPlayerHide,
		Binding: key.NewBinding(key.WithKeys("M"), key.WithHelp("", "Hide Media Player")),
	},
	KeyPlayerTogglePlay: {
		ID:      KeyPlayerTogglePlay,
		Binding: key.NewBinding(key.WithKeys(" "), key.WithHelp("", "Play/Pause")),
	},
	KeyPlayerNext: {
		ID:      KeyPlayerNext,
		Binding: key.NewBinding(key.WithKeys(">"), key.WithHelp("", "Next")),
	},
	KeyPlayerPrevious: {
		ID:      KeyPlayerPrevious,
		Binding: key.NewBinding(key.WithKeys("<"), key.WithHelp("", "Previous")),
	},
	KeyPlayerSeekForward: {
		ID:      KeyPlayerSeekForward,
		Binding: key.NewBinding(key.WithKeys("right"), key.WithHelp("", "Seek Forward")),
	},
	KeyPlayerSeekBackward: {
		ID:      KeyPlayerSeekBackward,
		Binding: key.NewBinding(key.WithKeys("left"), key.WithHelp("", "Seek Backward")),
	},
	KeyPlayerStop: {
		ID:      KeyPlayerStop,
		Binding: key.NewBinding(key.WithKeys("]"), key.WithHelp("", "Stop")),
	},
	KeyFilebrowserConfirmSelection: {
		ID:      KeyFilebrowserConfirmSelection,
		Binding: key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("", "Confirm Selection")),
	},
	KeyFilebrowserDirForward: {
		ID:      KeyFilebrowserDirForward,
		Binding: key.NewBinding(key.WithKeys("right"), key.WithHelp("", "Go Forward")),
	},
	KeyFilebrowserDirBack: {
		ID:      KeyFilebrowserDirBack,
		Binding: key.NewBinding(key.WithKeys("left"), key.WithHelp("", "Go Back")),
	},
	KeyFilebrowserSelect: {
		ID:      KeyFilebrowserSelect,
		Binding: key.NewBinding(key.WithKeys(" "), key.WithHelp("", "Select")),
	},
	KeyFilebrowserInvertSelection: {
		ID:      KeyFilebrowserInvertSelection,
		Binding: key.NewBinding(key.WithKeys("a"), key.WithHelp("", "Invert Selection")),
	},
	KeyFilebrowserSelectAll: {
		ID:      KeyFilebrowserSelectAll,
		Binding: key.NewBinding(key.WithKeys("A"), key.WithHelp("", "Select All")),
	},
	KeyFilebrowserRefresh: {
		ID:      KeyFilebrowserRefresh,
		Binding: key.NewBinding(key.WithKeys("ctrl+r"), key.WithHelp("", "Refresh")),
	},
	KeyFilebrowserToggleHidden: {
		ID:      KeyFilebrowserToggleHidden,
		Binding: key.NewBinding(key.WithKeys("h"), key.WithHelp("", "Hidden")),
	},
	KeyProgressTransferResume: {
		ID:      KeyProgressTransferResume,
		Binding: key.NewBinding(key.WithKeys("g"), key.WithHelp("", "Resume Transfer")),
	},
	KeyProgressTransferCancel: {
		ID:      KeyProgressTransferCancel,
		Binding: key.NewBinding(key.WithKeys("x"), key.WithHelp("", "Cancel Transfer")),
	},
	KeyProgressView: {
		ID:      KeyProgressView,
		Binding: key.NewBinding(key.WithKeys("v"), key.WithHelp("", "View Transfers")),
	},
	KeyProgressTransferSuspend: {
		ID:      KeyProgressTransferSuspend,
		Binding: key.NewBinding(key.WithKeys("z"), key.WithHelp("", "Suspend Transfer")),
	},
	KeyTasks: {
		ID:      KeyTasks,
		Binding: key.NewBinding(key.WithKeys("t"), key.WithHelp("", "Task")),
	},
	KeyTaskCancel: {
		ID:      KeyTaskCancel,
		Binding: key.NewBinding(key.WithKeys("ctrl+x"), key.WithHelp("", "Cancel Task")),
	},
}
