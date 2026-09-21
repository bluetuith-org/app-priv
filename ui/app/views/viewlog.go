package views

import (
	"github.com/ayn2op/tview"
	"github.com/bluetuith-org/bluetuith/ui/theme"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type logModel struct {
	*tview.TextView

	focused bool

	v rootView
}

// ViewID returns the view's ID.
func (l *logModel) ViewID() viewID {
	return viewIDLog
}

// Initialize initializes a view.
func (l *logModel) Initialize() error {
	l.TextView = tview.NewTextView()
	l.UpdateStyles()

	return nil
}

// SetRootView sets the root view upon which the view is rendered.
// This will enable the view to access app-specific functions and send
// routed messages.
func (l *logModel) SetRootView(v rootView) {
	l.v = v
}

// AttachToTabView attaches this view to the tabbed view.
func (l *logModel) AttachToTabView() (tabSection, bool) {
	return l, true
}

// HandleRouterMsg handles the routed message.
func (l *logModel) HandleRouterMsg(m routerMsg) tview.Cmd {
	return handleRouterMsg(l, m)
}

// SetFocus sets whether the view is currently focused.
func (l *logModel) SetFocus(focused bool) {
	l.focused = focused
}

// GetFocus gets whether the view is currently focused.
func (l *logModel) GetFocus() bool {
	return l.focused
}

// UpdateStyles updates the styles for the view.
func (l *logModel) UpdateStyles() {
	l.TextView.SetTextStyle(tcell.StyleDefault.Foreground(color.Black).Background(color.Gray))
}

// RefreshContent refreshes the content of the view.
func (l *logModel) RefreshContent() {
}

// Label returns the label for the tab.
func (l *logModel) Label() string {
	return "Log"
}

// Icon returns the icon for the tab.
func (l *logModel) Icon() *theme.IconVariant {
	return theme.Icons().Log
}

// Update receives messages when this model has focus.
func (l *logModel) Update(msg tview.Msg) tview.Cmd {
	return l.TextView.Update(msg)
}

// View draws this model onto the screen.
func (l *logModel) View(screen tview.Screen) {
	l.TextView.View(screen)
}

var _ tabSection = (*logModel)(nil)
