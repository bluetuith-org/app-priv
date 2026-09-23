package views

import (
	"context"
	"fmt"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/flex"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/ui/config"
	"github.com/bluetuith-org/bluetuith/ui/keybindings"
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

	// AttachToTabView attaches this view to the tabbed view.
	AttachToTabView() (tabSection, bool)

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

	// SendMsg sends a message to the program.
	SendMsg(msg tview.Msg)

	// SendRoutedUpdateMsg routes the message to the specified view.
	// Should only be called from the view's [view.Update] function.
	SendRoutedUpdateMsg(msg routerMsg) tview.Cmd
}

// ViewModel holds and renders a set of views for the application.
type ViewModel struct {
	AppBinder

	adTree   *adTreeModel
	tabsView *tabsModel

	infoView *infoModel
	logView  *logModel

	header *textModel

	vflex, hflex *flex.Model
	layout       *flex.Model

	initedViews map[viewID]view

	ctx     context.Context
	cancel  context.CancelFunc
	msgChan chan tview.Msg
}

// NewViews returns a new set of views with arranged layouts.
func NewViews(appBinder AppBinder) (*ViewModel, error) {
	ctx, cancel := context.WithCancel(context.Background())

	v := &ViewModel{
		AppBinder: appBinder,

		adTree:   &adTreeModel{},
		tabsView: &tabsModel{},

		infoView: &infoModel{},
		logView:  &logModel{},

		header: newTextModel(),
		vflex:  flex.NewModel(),
		hflex:  flex.NewModel(),
		layout: flex.NewModel(),

		initedViews: make(map[viewID]view, _viewIDMax),

		ctx:     ctx,
		cancel:  cancel,
		msgChan: make(chan tview.Msg, 1),
	}

	return v, v.initAllViews()
}

// Update receives messages when this model has focus.
func (v *ViewModel) Update(msg tview.Msg) tview.Cmd {
	switch m := msg.(type) {
	case tview.InitMsg:

	case tview.KeyMsg:
		switch {
		case kb().Quit.Matches(m):
			return v.quitCmd()

		default:
		}

		switch m.Str() {
		case "q":
			return v.quitCmd()

		default:
		}

	case externalMsg:
		return tview.Batch(v.parseViewMessages(m.msg), v.listenForMsg())

	default:
	}

	return v.layout.Update(msg)
}

func (v *ViewModel) parseViewMessages(msg tview.Msg) tview.Cmd {
	switch m := msg.(type) {
	case *tcell.CellBuffer:
		_ = m

	default:
	}

	return nil
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

type externalMsg struct {
	msg tview.Msg
}

func (v *ViewModel) listenForMsg() tview.Cmd {
	return func() tview.Msg {
		return externalMsg{<-v.msgChan}
	}
}

func (v *ViewModel) quitCmd() tview.Cmd {
	return tview.Batch(func() tview.Msg {
		v.cancel()
		return nil
	}, tview.Quit())
}

func (v *ViewModel) initAllViews() error {
	for _, view := range []view{
		v.tabsView,
		v.adTree,
		v.infoView,
		v.logView,
	} {
		view.SetRootView(v)

		if err := view.Initialize(); err != nil {
			return err
		}

		v.initedViews[view.ViewID()] = view
		if tab, ok := view.AttachToTabView(); ok {
			v.tabsView.AddTabSection(tab)
		}
	}

	v.layout = v.arrangeViews()

	return nil
}

func (v *ViewModel) arrangeViews() *flex.Model {
	v.UpdateStyles(true)

	v.header.SetText(fmt.Sprintf(" bluetuith %s (%s)", v.Configuration().Version, v.Configuration().Revision))

	status := tview.NewTextView()
	status.SetText(" status")

	treeModel := v.adTree
	treeModel.SetBorderPadding(2, 1, 1, 2)

	tabsBox := v.tabsView

	v.hflex.SetDirection(flex.DirectionColumn)
	v.hflex.AddItem(treeModel, 0, 1, true)
	v.hflex.AddItem(tabsBox, 0, 1, false)

	v.vflex.SetDirection(flex.DirectionRow)
	v.vflex.AddItem(v.header, 1, 0, false)
	v.vflex.AddItem(v.hflex, 0, 4, true)
	v.vflex.AddItem(status, 1, 0, false)

	return v.vflex
}

// SendMsg sends a message to the program.
func (v *ViewModel) SendMsg(msg tview.Msg) {
	select {
	case <-v.ctx.Done():
	case v.msgChan <- msg:
	}
}

var _ rootView = (*ViewModel)(nil)

func kb() *keybindings.Keybindings {
	return keybindings.Current
}
