package info

import (
	"fmt"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/text"
	"github.com/bluetuith-org/bluetuith/keybindings"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
)

// Model represents an information section.
type Model struct {
	text        []text.Span
	placeholder []text.Span

	rv ui.RootView
}

// New returns a new information section.
func New(rv ui.RootView) *Model {
	return &Model{
		rv:   rv,
		text: make([]text.Span, 0),
		placeholder: []text.Span{
			stext.NewSpan(
				"Settle on an adapter, device or the top of the tree to view information.",
				theme.Current().App.PlaceHolder,
			),
		},
	}
}

// Title returns the title of the tab section.
func (m *Model) Title() string {
	return fmt.Sprintf("%s (%s)", m.ViewID().String(), ui.Kb().Information.Show.KeyName)
}

// Icon returns the icon associated with the tab section.
func (m *Model) Icon() string {
	return theme.Icons().Info.String()
}

// ViewID returns the view's ID.
func (m *Model) ViewID() ui.ViewID {
	return ui.ViewIDInfo
}

// Header returns the content to print in the title bar area.
func (m *Model) Header() []text.Span {
	return nil
}

// Status returns the content to print in the status bar area.
func (m *Model) Status() []text.Span {
	return nil
}

// ViewKeybind returns the keybinding used to focus the view.
func (m *Model) ViewKeybind() keybindings.Keybinding {
	return ui.Kb().Information.Show
}

// Init returns a command to run when the model starts, or nil.
func (m *Model) Init() tview.Cmd {
	return nil
}

// Update receives messages when this model has focus.
func (m *Model) Update(msg tview.Msg, _ *ui.Response) {
	switch ms := msg.(type) {
	case ui.R[scrollMsg]:

	case ui.R[infoMsg]:
		m.text = ms.Value.content

	case ui.R[resetMsg]:
		m.text = m.placeholder

	default:
	}
}

// View draws this model onto the screen.
func (m *Model) View(_ bool) tview.Widget {
	return text.Rich(m.text...).
		Style(theme.Current().Global)
}

type scrollMsg struct{}

var (
	_ ui.Model      = (*Model)(nil)
	_ ui.TabSection = (*Model)(nil)
	_ ui.View       = (*Model)(nil)
)
