package views

import (
	"github.com/ayn2op/tview"
	"github.com/bluetuith-org/bluetuith/ui/theme"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/rivo/uniseg"
)

type tabSection interface {
	view

	// Label returns the label for the tab.
	Label() string

	// Icon returns the icon for the tab.
	Icon() *theme.IconVariant
}

type tabItem struct {
	tabSection

	label string
	width int
}

func newTabItem(section tabSection) *tabItem {
	n := " " + section.Icon().String() + " " + section.Label() + " "
	width := uniseg.StringWidth(n)

	return &tabItem{
		tabSection: section,
		label:      n,
		width:      width,
	}
}

type tabsModel struct {
	*tview.Box

	tabs      []*tabItem
	activeTab int
	focused   bool

	labelStyle, activeLabelStyle tcell.Style
	labelAlignment               tview.Alignment

	arrowLeft, arrowRight           string
	arrowLeftWidth, arrowRightWidth int
	arrowStyle                      tcell.Style

	maxWidth int

	v rootView
}

// ViewID returns the view's ID.
func (t *tabsModel) ViewID() viewID {
	return viewIDTabs
}

// Initialize initializes a view.
func (t *tabsModel) Initialize() error {
	t.Box = tview.NewBox()

	t.labelAlignment = tview.AlignmentCenter

	t.arrowLeft = theme.Icons().ArrowLeft.String() + " "
	t.arrowRight = " " + theme.Icons().ArrowRight.String()
	t.arrowLeftWidth = uniseg.StringWidth(t.arrowLeft)
	t.arrowRightWidth = uniseg.StringWidth(t.arrowRight)

	t.UpdateStyles()

	return nil
}

// SetRootView sets the root view upon which the view is rendered.
// This will enable the view to access app-specific functions and send
// routed messages.
func (t *tabsModel) SetRootView(v rootView) {
	t.v = v
}

// AttachToTabView attaches this view to the tabbed view.
func (t *tabsModel) AttachToTabView() (tabSection, bool) {
	return nil, false
}

// HandleRouterMsg handles the routed message.
func (t *tabsModel) HandleRouterMsg(msg routerMsg) tview.Cmd {
	return handleRouterMsg(t, msg)
}

// SetFocus sets whether the view is currently focused.
func (t *tabsModel) SetFocus(focused bool) {
	t.focused = focused
}

// GetFocus gets whether the view is currently focused.
func (t *tabsModel) GetFocus() bool {
	return t.focused
}

// UpdateStyles updates the styles for the view.
func (t *tabsModel) UpdateStyles() {
	t.Box.SetBackgroundColor(color.Gray)

	t.labelStyle = tcell.StyleDefault.Background(color.Gray).Foreground(color.White)
	t.activeLabelStyle = t.labelStyle.
		Foreground(color.Green).
		Background(color.Gray).
		Underline(true).
		Bold(true)

	t.arrowStyle = tcell.StyleDefault.Background(color.Gray).Foreground(color.White)
}

// RefreshContent refreshes the content of the view.
func (t *tabsModel) RefreshContent() {
}

// AddTabSection adds a tab section to the tab view.
func (t *tabsModel) AddTabSection(tab tabSection) {
	item := newTabItem(tab)
	t.maxWidth += item.width

	t.tabs = append(t.tabs, item)
}

// Previous switches to the previous view.
func (t *tabsModel) Previous() {
	t.activeTab = ((t.activeTab - 1) + len(t.tabs)) % len(t.tabs)
}

// Next switches to the next view.
func (t *tabsModel) Next() {
	t.activeTab = (t.activeTab + 1) % len(t.tabs)
}

// Update receives messages when this model has focus.
func (t *tabsModel) Update(msg tview.Msg) tview.Cmd {
	if len(t.tabs) == 0 {
		return t.Box.Update(msg)
	}

	switch msg := msg.(type) {
	case tview.KeyMsg:

		switch msg.Key() {
		case tcell.KeyTAB:
			t.Next()
			return nil

		default:
		}

	case tview.MouseMsg:
		x, y := msg.Position()
		if !t.InRect(x, y) {
			return nil
		}

		if tab, ok := t.tabAt(x, y); ok {
			switch msg.Action {
			case tview.MouseLeftClick:
				if tab == t.activeTab {
					return nil
				}

				t.activeTab = tab
				return nil

			case tview.MouseScrollUp, tview.MouseScrollLeft:
				t.Previous()
				return nil

			case tview.MouseScrollDown, tview.MouseScrollRight:
				t.Next()
				return nil
			}
		}
	}

	return t.tabs[t.activeTab].Update(msg)
}

// View draws this model onto the screen.
func (t *tabsModel) View(screen tview.Screen) {
	t.Box.View(screen)

	if len(t.tabs) == 0 {
		return
	}

	startX, y, maxWidth, maxHeight := t.InnerRect()
	y++

	defer t.setContent(screen, startX, y, maxWidth, maxHeight)

	if t.maxWidth <= maxWidth {
		x := startX + t.stripOffset(maxWidth)
		for active, ti := range t.tabs {
			style := t.labelStyle
			if active == t.activeTab {
				style = t.activeLabelStyle
			}

			screen.PutStrStyled(x, y, ti.label, style)
			x += ti.width
		}

		return
	}

	leftBoundaryX := startX + t.arrowLeftWidth
	rightBoundaryX := (startX + maxWidth) - t.arrowRightWidth
	usableWidth := rightBoundaryX - leftBoundaryX

	activeTab := t.tabs[t.activeTab]

	if activeTab.width > usableWidth {
		t.drawTruncatedRight(screen, leftBoundaryX, y, usableWidth, t.activeLabelStyle, activeTab.label)
		screen.PutStrStyled(startX, y, t.arrowLeft, t.arrowStyle)
		screen.PutStrStyled(rightBoundaryX, y, t.arrowRight, t.arrowStyle)
		return
	}

	activeStartX := leftBoundaryX + (usableWidth-activeTab.width)/2
	activeEndX := activeStartX + activeTab.width
	activeIndex := t.activeTab

	hasHiddenLeft := false
	hasHiddenRight := false
	gap := 1

	screen.PutStrStyled(activeStartX, y, activeTab.label, t.activeLabelStyle)

	currentLeftX := activeStartX - gap
	for i := activeIndex - 1; i >= 0; i-- {
		if currentLeftX <= leftBoundaryX {
			hasHiddenLeft = true
			break
		}

		tab := t.tabs[i]
		availableWidth := currentLeftX - leftBoundaryX

		if tab.width > availableWidth {
			t.drawTruncatedLeft(screen, leftBoundaryX, y, availableWidth, tab.width, t.labelStyle, tab.label)
			hasHiddenLeft = true
			break
		}

		tabX := currentLeftX - tab.width
		currentLeftX = tabX - gap

		screen.PutStrStyled(tabX, y, tab.label, t.labelStyle)
	}

	currentRightX := activeEndX + gap
	for i := activeIndex + 1; i < len(t.tabs); i++ {
		if currentRightX >= rightBoundaryX {
			hasHiddenRight = true
			break
		}

		tab := t.tabs[i]
		availableWidth := rightBoundaryX - currentRightX

		if tab.width > availableWidth {
			t.drawTruncatedRight(screen, currentRightX, y, availableWidth, t.labelStyle, tab.label)
			hasHiddenRight = true
			break
		}

		screen.PutStrStyled(currentRightX, y, tab.label, t.labelStyle)
		currentRightX += tab.width + gap
	}

	if hasHiddenLeft || activeIndex > 0 && currentLeftX <= leftBoundaryX {
		screen.PutStrStyled(startX, y, t.arrowLeft, t.arrowStyle)
	}
	if hasHiddenRight || activeIndex < len(t.tabs)-1 && currentRightX >= rightBoundaryX {
		screen.PutStrStyled(rightBoundaryX, y, t.arrowRight, t.arrowStyle)
	}
}

func (t *tabsModel) drawTruncatedRight(s tcell.Screen, x, y, maxWidth int, style tcell.Style, str string) {
	parsedWidth := 0
	state := -1

	for len(str) > 0 {
		if parsedWidth > maxWidth {
			return
		}

		cl, rest, boundaries, next := uniseg.StepString(str, state)
		str = rest
		state = next

		s.PutStrStyled(x+parsedWidth, y, cl, style)

		bw := boundaries >> uniseg.ShiftWidth
		parsedWidth += bw
	}
}

func (t *tabsModel) drawTruncatedLeft(s tcell.Screen, x, y, maxWidth, stringWidth int, style tcell.Style, str string) {
	parsedWidth := 0
	state := -1
	skipWidth := maxWidth - stringWidth

	if skipWidth < 0 {
		skipWidth = -skipWidth
	}

	for len(str) > 0 {
		cl, rest, boundaries, next := uniseg.StepString(str, state)
		str = rest
		state = next

		bw := boundaries >> uniseg.ShiftWidth
		if skipWidth > 0 {
			skipWidth -= bw
			continue
		}

		s.PutStrStyled(x+parsedWidth, y, cl, style)
		parsedWidth += bw
	}
}

func (t *tabsModel) setContent(screen tcell.Screen, startX, startY, width, height int) {
	content := t.tabs[t.activeTab]
	if content != nil {
		startY += 2
		height -= 3

		content.SetRect(startX, startY, width, height)
		content.View(screen)
	}
}

// Taken from: https://github.com/ayn2op/tview
func (t *tabsModel) tabAt(x, y int) (int, bool) {
	innerX, innerY, width, _ := t.InnerRect()
	if y != innerY {
		return 0, false
	}

	tmpX := innerX + t.stripOffset(width)
	for i, tab := range t.tabs {
		labelWidth := tab.width
		if x >= tmpX && x < tmpX+labelWidth {
			return i, true
		}
		tmpX += labelWidth + 1
	}

	return 0, false
}

// stripOffset returns the horizontal offset at which the tab strip starts so
// that the labels, laid out left to right and separated by a single space, are
// aligned as a group within the given width.
//
// Taken from: https://github.com/ayn2op/tview
func (t *tabsModel) stripOffset(width int) int {
	stripWidth := t.maxWidth - 1

	switch t.labelAlignment {
	case tview.AlignmentCenter:
		return max((width-stripWidth)/2, 0)

	case tview.AlignmentRight:
		return max(width-stripWidth, 0)

	default:
		return 0
	}
}

var _ view = (*tabsModel)(nil)
