package views

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/tree"
	"github.com/gdamore/tcell/v3/color"
)

// adTree represents a tree of adapters and its associated devices.
type adTree struct {
	*tree.Model

	focused bool
	rv      rootView
}

// ViewID returns the view's ID.
func (a *adTree) ViewID() viewID {
	return viewIDAdTree
}

// Initialize initializes a view.
func (a *adTree) Initialize() error {
	treeModel := tree.NewModel()
	treeModel.SetRoot(tree.NewNode("Root"))
	treeModel.SetCurrentNode(treeModel.Root())
	treeModel.Root().AddChild(tree.NewNode(" value"))
	treeModel.SetBackgroundColor(color.Gray)

	a.Model = treeModel

	return nil
}

// SetRootView sets the root view upon which the view is rendered.
// This will enable the view to access app-specific functions and send
// routed messages.
func (a *adTree) SetRootView(v rootView) {
	a.rv = v
}

// HandleRouterMsg handles the routed message.
func (a *adTree) HandleRouterMsg(m routerMsg) tview.Cmd {
	return handleRouterMsg(a, m)
}

// SetFocus sets whether the view is currently focused.
func (a *adTree) SetFocus(focused bool) {
	a.focused = focused
}

// GetFocus gets whether the view is currently focused.
func (a *adTree) GetFocus() bool {
	return a.focused
}

// UpdateStyles updates the styles for the view.
func (a *adTree) UpdateStyles() {
}

// RefreshContent refreshes the content of the view.
func (a *adTree) RefreshContent() {
}

// Update receives messages when this model has focus.
func (a *adTree) Update(msg tview.Msg) tview.Cmd {
	return a.Model.Update(msg)
}

// View draws this model onto the screen.
func (a *adTree) View(screen tview.Screen) {
	a.Model.View(screen)
}

// Rect returns the current position of the model, x, y, width, and
// height.
func (a *adTree) Rect() (x int, y int, width int, height int) {
	return a.Model.Rect()
}

// SetRect sets a new position of the model.
func (a *adTree) SetRect(x int, y int, width int, height int) {
	a.Model.SetRect(x, y, width, height)
}

var _ view = (*adTree)(nil)
