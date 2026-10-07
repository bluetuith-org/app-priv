package info

import (
	"fmt"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/richtext"
	"github.com/ayn2op/tview/textview"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/widgets/tabs"
)

// Model represents an information section.
type Model struct {
	state       textview.ScrollState
	text        richtext.Text
	placeholder richtext.Text
}

// New returns a new information section.
func New() Model {
	return Model{
		state: textview.ScrollState{},
		text:  make(richtext.Text, 0),
		placeholder: richtext.New(richtext.Line{
			richtext.NewSegment(
				"Settle on an adapter, device or the top of the tree to view information.",
				theme.Current().App.PlaceHolder,
			),
		}),
	}
}

// Title returns the title of the tab section.
func (m Model) Title() string {
	return fmt.Sprintf("Info (%s)", ui.Kb().Information.Show.KeyName)
}

// Icon returns the icon associated with the tab section.
func (m Model) Icon() string {
	return theme.Icons().Info.String()
}

// ViewID returns the view's ID.
func (m Model) ViewID() ui.ViewID {
	return ui.ViewIDInfo
}

// Init returns a command to run when the model starts, or nil.
func (m Model) Init() tview.Cmd {
	return nil
}

// Update receives messages when this model has focus.
func (m Model) Update(focused bool, msg tview.Msg) (Model, tview.Cmd) {
	cmd := m.updateModel(focused, msg)
	return m, cmd
}

// View draws this model onto the screen.
func (m Model) View(focused bool) tview.Widget {
	return textview.New(m.text).
		Focused(focused).
		ScrollState(&m.state).
		Style(theme.Current().Global).
		Wrap(true).
		WordWrap(true).
		OnChange(func(c textview.Change) tview.Msg {
			return scrollMsg(c)
		})
}

func (m *Model) updateModel(_ bool, msg tview.Msg) tview.Cmd {
	switch ms := msg.(type) {
	case scrollMsg:
		m.state.Apply(textview.Change(ms))

	case infoMsg:
		m.text = ms.content

	case resetMsg:
		m.text = m.placeholder

	default:
	}

	return nil
}

type scrollMsg textview.Change

var (
	_ ui.Model[Model] = Model{}
	_ tabs.TabSection = Model{}
	_ ui.View         = (*Model)(nil)
)
