package info

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/richtext"
	"github.com/ayn2op/tview/textview"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
)

// Model represents an information section.
type Model struct {
	state textview.ScrollState
	text  richtext.Text
}

// New returns a new information section.
func New() Model {
	return Model{
		state: textview.ScrollState{},
		text:  make(richtext.Text, 0),
	}
}

// Title returns the title of the tab section.
func (m Model) Title() string {
	return "Info"
}

// Icon returns the icon associated with the tab section.
func (m Model) Icon() string {
	return theme.Icons().Info.String()
}

// ViewID returns the view's ID.
func (m Model) ViewID() ui.ViewID {
	return ui.ViewIDInfo
}

// HandleRouterMsg handles the routed message.
func (m *Model) HandleRouterMsg(msg ui.RouterMsg) tview.Cmd {
	return m.updateModel(msg)
}

// Init returns a command to run when the model starts, or nil.
func (m Model) Init() tview.Cmd {
	return nil
}

// Update receives messages when this model has focus.
func (m Model) Update(msg tview.Msg) (Model, tview.Cmd) {
	cmd := m.updateModel(msg)
	return m, cmd
}

// View draws this model onto the screen.
func (m Model) View(focused bool) tview.Element {
	return textview.New(m.text).
		Focused(focused).
		ScrollState(&m.state).
		Style(theme.Current().Info.Style).
		OnChange(func(c textview.Change) tview.Msg {
			return scrollMsg(c)
		})
}

func (m *Model) updateModel(msg tview.Msg) tview.Cmd {
	switch ms := msg.(type) {
	case scrollMsg:
		m.state.Apply(textview.Change(ms))

	default:
	}

	return nil
}

type scrollMsg textview.Change

var (
	_ ui.Model[Model] = Model{}
	_ ui.View         = (*Model)(nil)
)
