package tabs

import (
	"github.com/ayn2op/tview"
	"github.com/gdamore/tcell/v3"
	"github.com/rivo/uniseg"
)

// TabSection defines a section of the tabs widget where
// the selected view is shown.
type TabSection interface {
	// Title returns the title of the tab section.
	Title() string

	// Icon returns the icon associated with the tab section.
	Icon() string
}

// Provider provides an interface to configure tab labels and arrows.
type Provider interface {
	// LabelStyles provides the styles for the inactive and active labels.
	LabelStyle(active bool) tcell.Style

	// ArrowConfig provides the left/right arrow icons and the style to be applied.
	ArrowConfig() (arrowLeft string, arrowRight string, arrowStyle tcell.Style, arrowWidth int)

	// HandleMsg handles the widget's message.
	HandleMsg(msg tview.Msg) tview.Msg
}

// Tabs returns a tab widget.
type Tabs[T Provider] struct {
	*State

	provider T

	focused bool
	content tview.Element
}

// New returns a new tab widget.
func New[T Provider](provider T, state *State) Tabs[T] {
	return Tabs[T]{State: state, provider: provider}
}

// Focused sets the widget to the focused state.
func (t Tabs[T]) Focused(focused bool) Tabs[T] {
	t.focused = focused
	return t
}

// Content sets the active content for the tabbed pane.
func (t Tabs[T]) Content(content tview.Element) Tabs[T] {
	t.content = content
	return t
}

// Draw draws the element onto the screen within the given area.
func (t Tabs[T]) Draw(screen tview.Screen, area tview.Rectangle) {
	if area.Width <= 0 || area.Height <= 0 {
		return
	}

	area.Y++

	startX, y, maxWidth := area.X, area.Y, area.Width
	defer t.setContent(screen, area)

	// All the tabs' labels fit within the width, draw them and exit.
	if t.maxWidth <= maxWidth {
		x := startX + t.stripOffset(maxWidth)
		for active, ti := range t.tabs {
			style := t.provider.LabelStyle(active == t.activeTab)

			screen.PutStrStyled(x, y, ti.label, style)
			x += ti.width
		}

		return
	}

	labelStyle, activeLabelStyle := t.provider.LabelStyle(false), t.provider.LabelStyle(true)
	arrowLeft, arrowRight, arrowStyle, arrowWidth := t.provider.ArrowConfig()

	leftBoundaryX := startX + arrowWidth               // reserve space for left arrow
	rightBoundaryX := (startX + maxWidth) - arrowWidth // reserve space for right arrow
	usableWidth := rightBoundaryX - leftBoundaryX

	activeTab := t.tabs[t.activeTab]

	// Only the active tab label will fit in the current width,
	// draw and truncate from the right as necessary, then exit.
	if activeTab.width > usableWidth {
		t.drawTruncatedRight(screen, leftBoundaryX, y, usableWidth, activeLabelStyle, activeTab.label)
		screen.PutStrStyled(startX, y, arrowLeft, arrowStyle)
		screen.PutStrStyled(rightBoundaryX, y, arrowRight, arrowStyle)
		return
	}

	// prevWidth is the combined widths of all tab labels before the active
	// tab label. If the prevWidth exceeds the usableWidth, just right-align the
	// active tab's label as much as visually possible.
	xPos := min(activeTab.prevWidth, usableWidth-activeTab.width)

	activeStartX := leftBoundaryX + xPos
	activeEndX := activeStartX + activeTab.width
	activeIndex := t.activeTab

	hasHiddenLeft := false
	hasHiddenRight := false
	gap := 1 // This could be a constant.

	// Draw the active tab's label.
	screen.PutStrStyled(activeStartX, y, activeTab.label, activeLabelStyle)

	// Walk back from the active tab's position, finding labels that we can
	// add to it's left side (after truncation, if necessary).
	currentLeftX := activeStartX - gap
	for i := activeIndex - 1; i >= 0; i-- {
		if currentLeftX <= leftBoundaryX {
			hasHiddenLeft = true
			break
		}

		tab := t.tabs[i]
		availableWidth := currentLeftX - leftBoundaryX

		if tab.width > availableWidth {
			t.drawTruncatedLeft(screen, leftBoundaryX, y, availableWidth, tab.width, labelStyle, tab.label)
			hasHiddenLeft = true
			break
		}

		tabX := currentLeftX - tab.width
		currentLeftX = tabX - gap

		screen.PutStrStyled(tabX, y, tab.label, labelStyle)
	}

	// Walk forward from the active tab's position, finding labels that we can
	// add to it's right side (after truncation, if necessary).
	currentRightX := activeEndX + gap
	for i := activeIndex + 1; i < len(t.tabs); i++ {
		if currentRightX >= rightBoundaryX {
			hasHiddenRight = true
			break
		}

		tab := t.tabs[i]
		availableWidth := rightBoundaryX - currentRightX

		if tab.width > availableWidth {
			t.drawTruncatedRight(screen, currentRightX, y, availableWidth, labelStyle, tab.label)
			hasHiddenRight = true
			break
		}

		screen.PutStrStyled(currentRightX, y, tab.label, labelStyle)
		currentRightX += tab.width + gap
	}

	if hasHiddenLeft || activeIndex > 0 && currentLeftX <= leftBoundaryX {
		screen.PutStrStyled(startX, y, arrowLeft, arrowStyle)
	}
	if hasHiddenRight || activeIndex < len(t.tabs)-1 && currentRightX >= rightBoundaryX {
		screen.PutStrStyled(rightBoundaryX, y, arrowRight, arrowStyle)
	}
}

// Handle translates an input message received within the given area into the message passed to Update. It returns nil to drop the message.
func (t Tabs[T]) Handle(msg tview.Msg, _ tview.Rectangle) tview.Msg {
	if !t.focused {
		return msg
	}

	return t.provider.HandleMsg(msg)
}

func (t Tabs[T]) drawTruncatedRight(s tcell.Screen, x, y, maxWidth int, style tcell.Style, str string) {
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

func (t Tabs[T]) drawTruncatedLeft(s tcell.Screen, x, y, maxWidth, stringWidth int, style tcell.Style, str string) {
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

func (t Tabs[T]) setContent(screen tcell.Screen, area tview.Rectangle) {
	content := t.content
	if content != nil {
		area.Y += 2
		area.Height -= 3

		content.Draw(screen, area)
	}
}

// stripOffset returns the horizontal offset at which the tab strip starts so
// that the labels, laid out left to right and separated by a single space, are
// aligned as a group within the given width.
//
// Taken from: https://github.com/ayn2op/tview
func (t Tabs[T]) stripOffset(width int) int {
	stripWidth := t.maxWidth - 1

	return max((width-stripWidth)/2, 0)
}
