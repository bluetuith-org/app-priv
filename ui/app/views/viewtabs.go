package views

import (
	"iter"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetuith/ui/keybindings"
	"github.com/bluetuith-org/bluetuith/ui/theme"
	tint "github.com/lrstanley/bubbletint/v2"
)

type tabSection interface {
	Title() string
	Icon() *theme.IconVariant

	viewer
}

type tabStyles struct {
	focused     lipgloss.Style
	blurred     lipgloss.Style
	selectedTab lipgloss.Style
	normalTab   lipgloss.Style
	arrowStyle  lipgloss.Style
}

type tabsView struct {
	tabs      []tabSection
	activeTab int

	width, height int
	focused       bool

	vport viewport.Model

	styles tabStyles

	v rootView
}

// ViewID returns the view's ID.
func (t *tabsView) ViewID() viewID {
	return viewIDTabs
}

// InitializeView initializes the view.
func (t *tabsView) InitializeView(_ *appfeatures.FeatureSet) (inited bool, err error) {
	t.tabs = make([]tabSection, 0, 5)
	t.activeTab = 0
	t.vport = viewport.New()

	t.setStyles()

	return true, nil
}

// SetRootView sets the root view.
func (t *tabsView) SetRootView(v rootView) {
	t.v = v
}

// AttachToTabView attaches this view to the tabbed view.
func (t *tabsView) AttachToTabView() (tabSection, bool) {
	return nil, false
}

// Resize resizes the view.
func (t *tabsView) Resize(width int, height int) tea.WindowSizeMsg {
	const (
		cHPaddingSize = 8
		cMinHeight    = 2
		cMinWidth     = 4
	)

	w, h := (width / 2), height

	t.vport.SetWidth(clamp(w-cHPaddingSize, cMinWidth, w))
	t.vport.SetHeight(clamp(cMinHeight, cMinHeight, h))

	t.width = w
	t.height = clamp((h - t.vport.Height()), 2, h)

	return tea.WindowSizeMsg{Width: t.width, Height: t.height}
}

// SetFocus sets whether the view is currently focused.
func (t *tabsView) SetFocus(focused bool) {
	t.focused = focused

	t.tabs[t.activeTab].SetFocus(focused)
}

// GetFocus gets whether the view is currently focused.
func (t *tabsView) GetFocus() bool {
	return t.focused
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (t *tabsView) Init() tea.Cmd {
	return collectCmds(false, nil, t.viewsIterator(false))
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (t *tabsView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		msg = t.Resize(m.Width, m.Height)

	case tea.KeyPressMsg:
		if !t.focused {
			return t, nil
		}

		switch {
		case keybindings.MatchesKey(keybindings.KeySwitch, m):
			t.move(true)
			return t, nil

		case keybindings.MatchesKey(keybindings.KeyClose, m):
			t.v.FocusTreeView()
			return t, nil

		default:
			return t, collectCmds(true, msg, t.viewsIterator(true))
		}
	}

	t.vport.Update(msg)

	return t, collectCmds(true, msg, t.viewsIterator(false))
}

// View renders the program's UI, which can be a string or a [Layer]. The
// view is rendered after every Update.
func (t *tabsView) View() tea.View {
	tabView := t.renderTabs()

	vertical := lipgloss.JoinVertical(lipgloss.Top, tabView, " ", t.tabs[t.activeTab].View().Content)

	return tea.NewView(vertical)
}

func (t *tabsView) AddTabSection(section tabSection) {
	t.tabs = append(t.tabs, section)
}

func (t *tabsView) move(fwd bool) {
	numTabs := len(t.tabs)
	prevPos := t.activeTab

	if numTabs == 0 {
		return
	}

	t.tabs[prevPos].SetFocus(false)

	switch fwd {
	case true:
		t.activeTab = (t.activeTab + 1) % numTabs

	case false:
		t.activeTab = ((t.activeTab - 1) + numTabs) % numTabs
	}

	t.tabs[t.activeTab].SetFocus(t.focused)
}

// HandleRouterMsg handles the routed message.
func (t *tabsView) HandleRouterMsg(m routerMsg) tea.Cmd {
	return handleRouterMsg(t, m)
}

func (t *tabsView) renderTabs() string {
	const (
		cTabPadding  = 2
		cTabsPadding = 4
	)

	totalTabs := len(t.tabs)
	if totalTabs == 0 {
		return ""
	}

	tabColStart, tabColEnd := 0, 0
	padStyle := lipgloss.NewStyle().Padding(0, cTabPadding)

	sections := make([]string, 0, len(t.tabs))
	for index, section := range t.tabs {
		widthBefore, widthAfter := 0, 0
		activeTabFound := false

		style := t.styles.normalTab
		if t.focused && index == t.activeTab {
			style = t.styles.selectedTab
		}

		renderedTitle := style.Render(section.Icon().Render(padStyle, section.Title()))
		switch {
		case index != t.activeTab && !activeTabFound:
			widthBefore += lipgloss.Width(renderedTitle)

		case index == t.activeTab:
			widthAfter = widthBefore + lipgloss.Width(renderedTitle)

			tabColStart = widthBefore
			tabColEnd = widthAfter

			activeTabFound = true
		}

		sections = append(sections, renderedTitle)
	}

	centered := lipgloss.Place(
		t.vport.Width(), t.vport.Height(),
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Center, sections...),
	)

	t.vport.Style = lipgloss.NewStyle().Align(lipgloss.Center)
	t.vport.SetContent(centered)
	t.vport.EnsureVisible(0, tabColStart+cTabsPadding, tabColEnd+cTabsPadding)

	arrowStyle := t.styles.arrowStyle
	arrowLeft, arrowRight := "  ", "  "
	switch {
	case t.activeTab == 0 && totalTabs > 1:
		arrowRight = theme.Icons().ArrowRight.Render(arrowStyle, "")

	case t.activeTab == totalTabs-1 && totalTabs > 1:
		arrowLeft = theme.Icons().ArrowLeft.Render(arrowStyle, "")

	case t.activeTab > 0 && t.activeTab < totalTabs:
		arrowLeft = theme.Icons().ArrowLeft.Render(arrowStyle, "")
		arrowRight = theme.Icons().ArrowRight.Render(arrowStyle, "")
	}

	s := lipgloss.JoinHorizontal(lipgloss.Center, arrowLeft, t.vport.View(), arrowRight)

	return s
}

func (t *tabsView) viewsIterator(focusedOnly bool) iter.Seq[viewer] {
	return func(yield func(viewer) bool) {
		if focusedOnly {
			yield(t.tabs[t.activeTab])
			return
		}

		for _, v := range t.tabs {
			if !yield(v) {
				return
			}
		}
	}
}

func (t *tabsView) setStyles() {
	s := tabStyles{
		blurred: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).Padding(1),

		normalTab: lipgloss.NewStyle().
			MarginTop(1).
			Foreground(tint.Current().Fg),

		arrowStyle: lipgloss.NewStyle().MarginTop(1).Padding(0, 1),
	}

	s.focused = s.blurred.BorderForeground(lipgloss.Color("62"))
	s.selectedTab = s.normalTab.
		Foreground(lipgloss.Darken(tint.Current().BrightPurple, 0.25)).
		Underline(true)

	t.styles = s
}

func clamp(v, low, high int) int {
	return min(max(v, low), high)
}
