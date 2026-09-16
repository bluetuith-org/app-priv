package views

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetuith/ui/theme"
)

type statusBarView struct {
	v rootView

	currMsg string

	width, height int
}

// ViewID returns the view's ID.
func (s *statusBarView) ViewID() viewID {
	return viewIDStatusBar
}

// InitializeView initializes the view.
func (s *statusBarView) InitializeView(_ *appfeatures.FeatureSet) (inited bool, err error) {
	return true, nil
}

// SetRootView sets the root view.
func (s *statusBarView) SetRootView(v rootView) {
	s.v = v
}

// AttachToTabView attaches this view to the tabbed view.
func (s *statusBarView) AttachToTabView() (tabSection, bool) {
	return nil, false
}

// Resize resizes the view.
func (s *statusBarView) Resize(width int, _ int) tea.WindowSizeMsg {
	s.width = width
	s.height = 1

	return tea.WindowSizeMsg{Width: s.width, Height: s.height}
}

// SetFocus sets whether the view is currently focused.
func (s *statusBarView) SetFocus(_ bool) {
}

// GetFocus gets whether the view is currently focused.
func (s *statusBarView) GetFocus() bool {
	return false
}

// HandleRouterMsg handles the routed message.
func (s *statusBarView) HandleRouterMsg(m routerMsg) tea.Cmd {
	return handleRouterMsg(s, m)
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (s *statusBarView) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (s *statusBarView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		msg = s.Resize(m.Width, m.Height)

	default:
	}

	return s, nil
}

// View renders the program's UI, which can be a string or a [Layer]. The
// view is rendered after every Update.
func (s *statusBarView) View() tea.View {
	style := theme.Current().StatusBar.Style.
		Width(s.width).
		Height(s.height).
		MaxWidth(s.width).
		MaxHeight(s.height).
		Align(lipgloss.Left)

	text := ""

	return tea.NewView(style.Render(" " + text))
}
