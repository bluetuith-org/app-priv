package tabpane

import (
	"iter"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/box"
	"github.com/ayn2op/tview/column"
	"github.com/ayn2op/tview/tabs"
	"github.com/ayn2op/tview/text"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
)

// Model holds a tabbed pane.
type Model struct {
	handler *viewHandler

	rv ui.RootView
}

// New returns a new tabbed pane.
func New(rv ui.RootView) *Model {
	return &Model{
		rv:      rv,
		handler: newHandler(rv),
	}
}

// ViewID returns the view's ID.
func (m *Model) ViewID() ui.ViewID {
	return ui.ViewIDTabs
}

// Header returns the content to print in the title bar.
func (m *Model) Header() iter.Seq[[]text.Span] {
	return func(yield func([]text.Span) bool) {
		for _, v := range m.handler.views {
			if content := v.Header(); content != nil {
				if !yield(content) {
					return
				}
			}
		}
	}
}

// Status returns the content to print in the status bar.
func (m *Model) Status() iter.Seq[[]text.Span] {
	return func(yield func([]text.Span) bool) {
		for _, v := range m.handler.views {
			if content := v.Status(); content != nil {
				if !yield(content) {
					return
				}
			}
		}
	}
}

// GetActiveID returns the current active tab view.
func (m *Model) GetActiveID() ui.ViewID {
	return m.handler.activeID()
}

// HandleRouterMsg handles the routed message.
func (m *Model) HandleRouterMsg(id ui.ViewID, routedMsg tview.Msg, rsp *ui.Response) {
	if id == ui.ViewIDTabs {
		m.Update(routedMsg, rsp)
		return
	}

	m.handler.updateRouted(id, routedMsg, rsp)
}

// Keybind converts a [tview.KeyMsg] to an [Action].
func (m *Model) Keybind(msg tview.KeyMsg) tabs.Action {
	if !m.rv.PanelFocused() {
		return tabs.ActionNone
	}

	switch {
	case ui.Kb().TabPane.Next.Matches(msg):
		return tabs.ActionNext

	case ui.Kb().TabPane.Previous.Matches(msg):
		return tabs.ActionPrevious
	}

	return tabs.ActionNone
}

// OnSelect is called when a new tab item is selected.
func (m *Model) OnSelect(index int) tview.Msg {
	return ui.ViewIDTabs.RouterMessage(selectedMsg{index})
}

// Init returns a command to run when the model starts, or nil.
func (m *Model) Init() tview.Cmd {
	return m.handler.init()
}

// Update returns the model changed in response to a message and a command to run, or nil.
func (m *Model) Update(msg tview.Msg, rsp *ui.Response) {
	switch ms := msg.(type) {
	case ui.R[selectedMsg]:
		m.handler.setActiveIndex(ms.Value.index)
		m.rv.SetFocused(m.handler.activeID())

	default:
		m.handler.updateAll(msg, rsp)
	}
}

// View returns the element that draws the model.
func (m *Model) View(focused bool) tview.Widget {
	bgStyle := theme.Current().Global

	t := tabs.New(m.handler.labels...).
		Active(m.handler.activeIndex).
		Style(theme.Current().TabsPane.Tab).
		ActiveStyle(theme.Current().TabsPane.FocusedTab).
		Alignment(tview.AlignmentCenter).
		Separator(" ").
		Padding(" ", " ").
		Arrows(theme.Icons().ArrowLeft.String(), theme.Icons().ArrowRight.String()).
		ArrowStyle(theme.Current().TabsPane.ArrowStyle).
		ClickableArrows(true).
		Wrap(true).
		Keybind(m.Keybind).
		OnSelect(m.OnSelect)

	content := box.New(m.handler.activeContent(focused)).
		Background(bgStyle.GetBackground()).
		Padding(1, 0, 0, 0)

	return box.New(column.New(t, content)).
		Background(bgStyle.GetBackground()).
		Padding(1, 0, 0, 0)
}

// SetViewKeybind focuses the view if the keybind matches.
func (m *Model) SetViewKeybind(msg tview.KeyMsg) bool {
	return m.handler.setActive(msg)
}

type selectedMsg struct {
	index int
}

var (
	_ ui.Model = (*Model)(nil)
	_ ui.View  = (*Model)(nil)
)
