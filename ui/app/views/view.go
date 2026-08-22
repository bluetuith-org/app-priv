package views

import (
	"fmt"
	"iter"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/76creates/stickers/flexbox"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/ui/config"
	tint "github.com/lrstanley/bubbletint/v2"
)

// Views should be deleted.
// TODO: Delete.
type Views struct{}

type viewer interface {
	// ViewID returns the view's ID.
	ViewID() viewID

	// InitializeView initializes the view.
	InitializeView(features *appfeatures.FeatureSet) (inited bool, err error)

	// SetRootView sets the root view.
	SetRootView(v rootView)

	// AttachToTabView attaches this view to the tabbed view.
	AttachToTabView() (tabSection, bool)

	// Resize resizes the view.
	Resize(width, height int) tea.WindowSizeMsg

	// SetFocus sets whether the view is currently focused.
	SetFocus(focused bool)

	// GetFocus gets whether the view is currently focused.
	GetFocus() bool

	// UpdateStyles updates the styles for the view.
	UpdateStyles()

	// handleRouterMsg handles the routed message.
	handleRouterMsg(m routerMsg) tea.Cmd

	tea.Model
}

type rootView interface {
	// ViewWidth returns the total width of the view.
	ViewWidth() int

	// ViewHeight returns the total height of the view.
	ViewHeight() int

	// Icons returns the configured readonly icon set.
	Icons() *iconSet

	// FocusTreeView focuses the tree view.
	FocusTreeView()

	// FocusTabView focuses the tab view.
	FocusTabView()

	// SendRoutedUpdateMsg routes the message to the specified view.
	// Should only be called from the view's [Update] function.
	SendRoutedUpdateMsg(msg routerMsg) tea.Cmd

	AppBinder
}

// AppBinder represents the app interface.
type AppBinder interface {
	// Session returns the current bluetooth session.
	Session() bluetooth.Session

	// Features returns the features supported by the session.
	Features() *appfeatures.FeatureSet

	// Configuration returns the application's configuration.
	Configuration() *config.Config

	// SendMsg sends a [tea.Msg] to the program.
	SendMsg(msg tea.Msg)
}

// ViewModel represents a complete view.
type ViewModel struct {
	width, height         int
	prevWidth, prevHeight int

	tabsView       *tabsView
	adTreeView     *adTree
	infoView       *infoView
	operationsView *operationsView

	initedViews map[viewID]viewer

	icons *iconSet

	flexVertical   *flexbox.FlexBox
	flexHorizontal *flexbox.HorizontalFlexBox
	cells          []*flexbox.Cell

	AppBinder
}

// NewViewModel returns the main view.
func NewViewModel(appBinder AppBinder) (*ViewModel, error) {
	v := &ViewModel{
		width:  120,
		height: 30,

		tabsView:    &tabsView{},
		initedViews: make(map[viewID]viewer),

		adTreeView:     &adTree{},
		infoView:       &infoView{},
		operationsView: &operationsView{},

		// TODO: Check for ascii icons.
		icons: newIconSet(false),

		flexVertical:   flexbox.New(0, 0),
		flexHorizontal: flexbox.NewHorizontal(0, 0),

		AppBinder: appBinder,
	}

	v.cells = []*flexbox.Cell{
		flexbox.NewCell(2, 1).SetContentGenerator(func(_, _ int) string {
			return v.renderHeader()
		}),
		flexbox.NewCell(2, -1).SetContentGenerator(func(_, _ int) string {
			return v.flexVertical.Render()
		}),
		flexbox.NewCell(1, 1).SetContentGenerator(func(_, _ int) string {
			return v.adTreeView.View().Content
		}),
		flexbox.NewCell(1, 1).SetContentGenerator(func(_, _ int) string {
			return v.tabsView.View().Content
		}),
	}

	v.flexHorizontal.AddColumns([]*flexbox.Column{
		v.flexHorizontal.NewColumn().AddCells(v.cells[0], v.cells[1]),
	})

	v.flexVertical.AddRows([]*flexbox.Row{
		v.flexVertical.NewRow().AddCells(v.cells[2], v.cells[3]),
	})

	return v, v.initAllViews()
}

// ViewWidth returns the total width of the view.
func (v *ViewModel) ViewWidth() int {
	return v.width
}

// ViewHeight returns the total height of the view.
func (v *ViewModel) ViewHeight() int {
	return v.height
}

// Icons returns the configured readonly icon set.
func (v *ViewModel) Icons() *iconSet {
	return v.icons
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (v *ViewModel) Init() tea.Cmd {
	return collectCmds(false, nil, v.modelIterator())
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
//
//	TODO: Implement focus
func (v *ViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		v.width = m.Width
		v.height = m.Height

	case tea.KeyPressMsg:
		switch m.Code {
		case 'q':
			return v, tea.Quit

		default:
			vid := viewIDNone

			for view := range v.modelIterator() {
				if view.GetFocus() {
					vid = view.ViewID()
					break
				}
			}

			return v, v.updateViewByID(vid, m)
		}

	case routerMsg:
		view, ok := v.getViewByID(m.id)
		if !ok {
			return v, nil
		}

		return v, view.handleRouterMsg(m)
	}

	return v, collectCmds(true, msg, v.modelIterator())
}

// View renders the program's UI, which can be a string or a [Layer]. The
// view is rendered after every Update.
func (v *ViewModel) View() tea.View {
	var view tea.View

	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion

	fbrow := v.flexVertical
	fbcol := v.flexHorizontal
	if v.prevWidth != v.width || v.prevHeight != v.height {
		fbrow.SetHeight(v.height)
		fbrow.SetWidth(v.width)

		fbcol.SetHeight(v.height)
		fbcol.SetWidth(v.width)
	}

	v.prevWidth = v.width
	v.prevHeight = v.height

	view.SetContent(fbcol.Render())

	return view
}

// UpdateStyles updates the styles for the view.
func (v *ViewModel) UpdateStyles() {
	for _, view := range v.initedViews {
		view.UpdateStyles()
	}
}

// FocusTreeView focuses the tree view.
func (v *ViewModel) FocusTreeView() {
	v.tabsView.SetFocus(false)
	v.adTreeView.SetFocus(true)
}

// FocusTabView focuses the tab view.
func (v *ViewModel) FocusTabView() {
	v.adTreeView.SetFocus(false)
	v.tabsView.SetFocus(true)
}

// SendRoutedUpdateMsg routes the message to the specified view.
// Should only be called from the view's [Update] function.
func (v *ViewModel) SendRoutedUpdateMsg(msg routerMsg) tea.Cmd {
	if !msg.isValid() {
		return nil
	}

	_, cmd := v.Update(msg)
	return cmd
}

func (v *ViewModel) initAllViews() error {
	tint.NewDefaultRegistry()
	tint.SetTint(tint.TintDraculaPlus)

	for _, view := range [maxViews]viewer{
		v.tabsView,
		v.adTreeView,
		v.infoView,
		v.operationsView,
	} {
		view.SetRootView(v)

		inited, err := view.InitializeView(v.Features())
		if err != nil {
			return err
		}

		if inited {
			v.initedViews[view.ViewID()] = view
			if tab, ok := view.AttachToTabView(); ok {
				v.tabsView.AddTabSection(tab)
			}
		}
	}

	return nil
}

func (v *ViewModel) modelIterator() iter.Seq[viewer] {
	return func(yield func(viewer) bool) {
		for _, view := range []viewer{v.adTreeView, v.tabsView} {
			if !yield(view) {
				return
			}
		}
	}
}

func (v *ViewModel) getViewByID(id viewID) (viewer, bool) {
	if !id.isValid() {
		return nil, false
	}

	view, ok := v.initedViews[id]
	return view, ok
}

func (v *ViewModel) updateViewByID(id viewID, msg tea.Msg) tea.Cmd {
	view, ok := v.getViewByID(id)
	if !ok {
		return nil
	}

	_, cmd := view.Update(msg)
	return cmd
}

func (v *ViewModel) renderHeader() string {
	style := lipgloss.NewStyle().
		Bold(true).
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("15")).
		Width(v.width).
		Align(lipgloss.Left)

	cfg := v.Configuration()
	title := useBuffer(len(cfg.Version)+len(cfg.Revision)+10, func(b *strings.Builder) {
		fmt.Fprintf(b, " bluetuith %s (%s)", cfg.Version, cfg.Revision)
	})

	return style.Render(title)
}

func collectCmds(update bool, msg tea.Msg, views iter.Seq[viewer]) tea.Cmd {
	cmds := make([]tea.Cmd, 0, 5)

	for view := range views {
		var cmd tea.Cmd

		if update {
			_, cmd = view.Update(msg)
		} else {
			cmd = view.Init()
		}

		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if len(cmds) == 0 {
		return nil
	}

	switch len(cmds) {
	case 0:
		return nil

	case 1:
		return cmds[0]
	}

	return tea.Batch(cmds...)
}
