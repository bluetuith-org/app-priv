package views

import (
	"fmt"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/flex"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/ui/config"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

// AppBinder represents an interface to interact with the root application itself.
type AppBinder interface {
	// Session returns the current session.
	Session() bluetooth.Session

	// Features returns the current features of the session.
	Features() *appfeatures.FeatureSet

	// Configuration returns the app's configuration.
	Configuration() *config.Config
}

// view represents a single view.
// All views must implement this interface.
type view interface {
	tview.Model

	// ViewID returns the view's ID.
	ViewID() viewID

	// Initialize initializes a view.
	Initialize() error

	// SetRootView sets the root view upon which the view is rendered.
	// This will enable the view to access app-specific functions and send
	// routed messages.
	SetRootView(v rootView)

	// HandleRouterMsg handles the routed message.
	HandleRouterMsg(m routerMsg) tview.Cmd

	// SetFocus sets whether the view is currently focused.
	SetFocus(focused bool)

	// GetFocus gets whether the view is currently focused.
	GetFocus() bool

	// UpdateStyles updates the styles for the view.
	UpdateStyles()

	// RefreshContent refreshes the content of the view.
	RefreshContent()
}

// rootView represents the root view. All views must inherit this interface,
// to call app-related functions and route messages to components.
type rootView interface {
	AppBinder

	// SendRoutedUpdateMsg routes the message to the specified view.
	// Should only be called from the view's [view.Update] function.
	SendRoutedUpdateMsg(msg routerMsg) tview.Cmd
}

// ViewModel holds and renders a set of views for the application.
type ViewModel struct {
	AppBinder

	adTreeModel *adTree

	header *tview.TextView

	vflex, hflex *flex.Model
	layout       *flex.Model

	initedViews map[viewID]view
}

// NewViews returns a new set of views with arranged layouts.
func NewViews(appBinder AppBinder) (*ViewModel, error) {
	v := &ViewModel{
		AppBinder: appBinder,

		adTreeModel: &adTree{},

		header: tview.NewTextView(),
		vflex:  flex.NewModel(),
		hflex:  flex.NewModel(),
		layout: flex.NewModel(),

		initedViews: make(map[viewID]view, _viewIDMax),
	}

	return v, v.initAllViews()
}

// Update receives messages when this model has focus.
func (v *ViewModel) Update(msg tview.Msg) tview.Cmd {
	switch m := msg.(type) {
	case tview.KeyMsg:
		switch m.Str() {
		case "q":
			return tview.Quit()

		default:
		}

	default:
	}

	return v.layout.Update(msg)
}

// View draws this model onto the screen.
func (v *ViewModel) View(screen tview.Screen) {
	v.layout.View(screen)
}

// Rect returns the current position of the model, x, y, width, and
// height.
func (v *ViewModel) Rect() (x int, y int, width int, height int) {
	return v.layout.Rect()
}

// SetRect sets a new position of the model.
func (v *ViewModel) SetRect(x int, y int, width int, height int) {
	v.layout.SetRect(x, y, width, height)
}

// SendRoutedUpdateMsg routes the message to the specified view.
// Should only be called from the view's [Update] function.
func (v *ViewModel) SendRoutedUpdateMsg(msg routerMsg) tview.Cmd {
	if !msg.isValid() {
		return nil
	}

	return v.Update(msg)
}

// UpdateStyles updates all styles for all views.
func (v *ViewModel) UpdateStyles(init bool) {
	if !init {
		for _, v := range v.initedViews {
			v.UpdateStyles()
		}
	}

	v.vflex.SetBackgroundColor(color.LightGray)
	v.hflex.SetBackgroundColor(color.LightGray)

	v.header.SetTextStyle(
		tcell.StyleDefault.Background(color.Purple).Foreground(color.Black).Bold(true),
	)
}

func (v *ViewModel) initAllViews() error {
	for _, view := range []view{
		v.adTreeModel,
	} {
		view.SetRootView(v)

		if err := view.Initialize(); err != nil {
			return err
		}

		v.initedViews[view.ViewID()] = view
		/*if tab, ok := view.AttachToTabView(); ok {
			v.tabsView.AddTabSection(tab)
		}*/
	}

	v.layout = v.arrangeViews()

	return nil
}

func (v *ViewModel) arrangeViews() *flex.Model {
	v.UpdateStyles(true)

	v.header.SetText(fmt.Sprintf(" bluetuith %s (%s)", v.Configuration().Version, v.Configuration().Revision))

	status := tview.NewTextView()
	status.SetText(" status")

	treeModel := v.adTreeModel
	treeModel.SetBorderPadding(2, 1, 1, 2)

	tabsBox := tview.NewBox()
	tabsBox.SetBackgroundColor(color.Gray)
	tabsBox.SetBorders(tview.BordersAll)
	tabsBox.SetBorderSet(tview.BorderSetRound())
	tabsBox.SetBorderPadding(2, 1, 1, 2)

	v.hflex.SetDirection(flex.DirectionColumn)
	v.hflex.AddItem(treeModel, 0, 1, true)
	v.hflex.AddItem(tabsBox, 0, 1, false)

	v.vflex.SetDirection(flex.DirectionRow)
	v.vflex.AddItem(v.header, 1, 0, false)
	v.vflex.AddItem(v.hflex, 0, 4, true)
	v.vflex.AddItem(status, 1, 0, false)

	return v.vflex
}

var _ rootView = (*ViewModel)(nil)
