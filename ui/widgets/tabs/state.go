package tabs

// State represents the current tab and selection state.
type State struct {
	tabs      []*Item
	activeTab int

	maxWidth int
}

// NewState returns a new state.
func NewState() *State {
	return &State{}
}

// AddSection adds a tab section.
func (t *State) AddSection(icon, title string) {
	item := NewItem(icon, title, t.maxWidth)
	t.maxWidth += item.width

	t.tabs = append(t.tabs, item)
}

// Next selects the next tab item.
func (t *State) Next() {
	t.activeTab = (t.activeTab + 1) % len(t.tabs)
}

// ActiveIndex returns the content of the selected tab.
func (t *State) ActiveIndex() int {
	return t.activeTab
}
