package views

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/tree"
	"github.com/gdamore/tcell/v3/color"
)

// adTreeModel represents a tree of adapters and its associated devices.
type adTreeModel struct {
	*tree.Model

	focused bool
	rv      rootView
}

// ViewID returns the view's ID.
func (a *adTreeModel) ViewID() viewID {
	return viewIDAdTree
}

// Initialize initializes a view.
func (a *adTreeModel) Initialize() error {
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
func (a *adTreeModel) SetRootView(v rootView) {
	a.rv = v
}

// AttachToTabView attaches this view to the tabbed view.
func (a *adTreeModel) AttachToTabView() (tabSection, bool) {
	return nil, false
}

// HandleRouterMsg handles the routed message.
func (a *adTreeModel) HandleRouterMsg(m routerMsg) tview.Cmd {
	return handleRouterMsg(a, m)
}

// SetFocus sets whether the view is currently focused.
func (a *adTreeModel) SetFocus(focused bool) {
	a.focused = focused
}

// GetFocus gets whether the view is currently focused.
func (a *adTreeModel) GetFocus() bool {
	return a.focused
}

// UpdateStyles updates the styles for the view.
func (a *adTreeModel) UpdateStyles() {
}

// RefreshContent refreshes the content of the view.
func (a *adTreeModel) RefreshContent() {
}

// Update receives messages when this model has focus.
func (a *adTreeModel) Update(msg tview.Msg) tview.Cmd {
	return a.Model.Update(msg)
}

// View draws this model onto the screen.
func (a *adTreeModel) View(screen tview.Screen) {
	a.Model.View(screen)
}

var _ view = (*adTreeModel)(nil)
