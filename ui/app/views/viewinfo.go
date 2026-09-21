package views

import (
	"github.com/ayn2op/tview"
	"github.com/bluetuith-org/bluetuith/ui/theme"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type infoModel struct {
	*tview.TextView

	focused bool

	v rootView
}

// ViewID returns the view's ID.
func (i *infoModel) ViewID() viewID {
	return viewIDInfo
}

// Initialize initializes a view.
func (i *infoModel) Initialize() error {
	i.TextView = tview.NewTextView()
	i.UpdateStyles()
	i.TextView.SetBorders(tview.BordersAll)

	i.TextView.SetText("Information")

	return nil
}

// SetRootView sets the root view upon which the view is rendered.
// This will enable the view to access app-specific functions and send
// routed messages.
func (i *infoModel) SetRootView(v rootView) {
	i.v = v
}

// AttachToTabView attaches this view to the tabbed view.
func (i *infoModel) AttachToTabView() (tabSection, bool) {
	return i, true
}

// HandleRouterMsg handles the routed message.
func (i *infoModel) HandleRouterMsg(m routerMsg) tview.Cmd {
	return handleRouterMsg(i, m)
}

// SetFocus sets whether the view is currently focused.
func (i *infoModel) SetFocus(focused bool) {
	i.focused = focused
}

// GetFocus gets whether the view is currently focused.
func (i *infoModel) GetFocus() bool {
	return i.focused
}

// UpdateStyles updates the styles for the view.
func (i *infoModel) UpdateStyles() {
	i.TextView.SetTextStyle(tcell.StyleDefault.Foreground(color.Black).Background(color.Gray))
	i.TextView.SetBorderStyle(tcell.StyleDefault.Background(color.Gray).Foreground(color.White))
}

// RefreshContent refreshes the content of the view.
func (i *infoModel) RefreshContent() {
}

// Label returns the label for the tab.
func (i *infoModel) Label() string {
	return "Info"
}

// Icon returns the icon for the tab.
func (i *infoModel) Icon() *theme.IconVariant {
	return theme.Icons().Info
}

// Update receives messages when this model has focus.
func (i *infoModel) Update(msg tview.Msg) tview.Cmd {
	return i.TextView.Update(msg)
}

// View draws this model onto the screen.
func (i *infoModel) View(screen tview.Screen) {
	i.TextView.View(screen)
}

var _ tabSection = (*infoModel)(nil)
