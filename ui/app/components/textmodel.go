package components

import (
	"github.com/ayn2op/tview"
	"github.com/gdamore/tcell/v3"
)

// TextModel holds a single line of text.
type TextModel struct {
	*tview.Box

	text      string
	textStyle tcell.Style

	align tview.Alignment

	x, y, w, h int
}

// NewTextModel returns a new textModel.
func NewTextModel() *TextModel {
	m := &TextModel{
		Box:   tview.NewBox(),
		align: tview.AlignmentRight,
	}

	return m
}

// SetText sets the text for the model.
func (t *TextModel) SetText(txt string) {
	t.text = txt
}

// SetAlign sets the text alignment.
func (t *TextModel) SetAlign(align tview.Alignment) {
	t.align = align
}

// SetTextStyle sets the text style for the model.
func (t *TextModel) SetTextStyle(style tcell.Style) {
	t.textStyle = style
	t.Box.SetBackgroundColor(style.GetBackground())
}

// Update receives messages when this model has focus.
func (t *TextModel) Update(tview.Msg) tview.Cmd { return nil }

// View draws this model onto the screen.
func (t *TextModel) View(screen tview.Screen) {
	if t.w <= 0 {
		return
	}

	t.Box.View(screen)
	tview.PrintWithStyle(screen, t.text, t.x, t.y, t.w, t.align, t.textStyle)
}

// Rect returns the current position of the model, x, y, width, and
// height.
func (t *TextModel) Rect() (x int, y int, width int, height int) {
	return t.x, t.y, t.w, t.h
}

// SetRect sets a new position of the model.
func (t *TextModel) SetRect(x, y, w, h int) {
	t.Box.SetRect(x, y, w, h)
	t.x, t.y, t.w, t.h = x, y, w, h
}

// Height returns the height of the textmodel.
func (t *TextModel) Height(int) int {
	return 1
}
