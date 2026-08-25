package views

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
)

type infoView struct {
	vport viewport.Model

	width, height int
	focused       bool

	vp viewport.Model

	v rootView
}

func (i *infoView) Title() string {
	return "Info"
}

func (i *infoView) Icon() *iconVariant {
	return i.v.Icons().Info
}

// ViewID returns the view's ID.
func (i *infoView) ViewID() viewID {
	return viewIDInfo
}

// InitializeView initializes the view.
func (i *infoView) InitializeView(_ *appfeatures.FeatureSet) (inited bool, err error) {
	i.vp = viewport.New()

	return true, nil
}

// SetRootView sets the root view.
func (i *infoView) SetRootView(v rootView) {
	i.v = v
}

// AttachToTabView attaches this view to the tabbed view.
func (i *infoView) AttachToTabView() (tabSection, bool) {
	return i, true
}

// Resize resizes the view.
func (i *infoView) Resize(width int, height int) tea.WindowSizeMsg {
	i.width, i.height = width, height

	return tea.WindowSizeMsg{Width: i.width, Height: i.height}
}

// SetFocus sets whether the view is currently focused.
func (i *infoView) SetFocus(focused bool) {
	i.focused = focused
}

// GetFocus gets whether the view is currently focused.
func (i *infoView) GetFocus() bool {
	return i.focused
}

// UpdateStyles updates the styles for the view.
func (i *infoView) UpdateStyles() {
}

// HandleRouterMsg handles the routed message.
func (i *infoView) HandleRouterMsg(m routerMsg) tea.Cmd {
	return handleRouterMsg(i, m)
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (i *infoView) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (i *infoView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		msg = i.Resize(m.Width, m.Height)

	case tea.KeyMsg:
	}

	return i, nil
}

// View renders the program's UI, which can be a string or a [Layer]. The
// view is rendered after every Update.
func (i *infoView) View() tea.View {
	i.vp.Style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))

	i.vp.SetWidth(i.width)
	i.vp.SetHeight(i.height)
	i.vp.SetContent("Information")

	return tea.NewView(i.vp.View())
}
