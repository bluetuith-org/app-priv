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
	text        richtext.Text
	placeholder richtext.Text

	rv ui.RootView
}

// New returns a new information section.
func New(rv ui.RootView) Model {
	return Model{
		rv:   rv,
		text: make(richtext.Text, 0),
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
	return fmt.Sprintf("%s (%s)", m.ViewID().String(), ui.Kb().Information.Show.KeyName)
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
func (m Model) Update(msg tview.Msg) (Model, tview.Cmd) {
	cmd := m.updateModel(msg)
	return m, cmd
}

// View draws this model onto the screen.
func (m Model) View(_ bool) tview.Widget {
	return textview.New(m.text).
		Style(theme.Current().Global).
		Wrap(true).
		WordWrap(true)
}

func (m *Model) updateModel(msg tview.Msg) tview.Cmd {
	switch ms := msg.(type) {
	case scrollMsg:

	case infoMsg:
		m.text = ms.content

	case resetMsg:
		m.text = m.placeholder

	default:
	}

	return nil
}

type scrollMsg struct{}

var (
	_ ui.Model[Model] = Model{}
	_ tabs.TabSection = Model{}
	_ ui.View         = (*Model)(nil)
)
