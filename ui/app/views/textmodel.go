package views

import (
	"github.com/ayn2op/tview"
	"github.com/gdamore/tcell/v3"
)

type textModel struct {
	*tview.Box

	text      string
	textStyle tcell.Style

	align tview.Alignment

	x, y, w, h int
}

func newTextModel() *textModel {
	m := &textModel{
		Box:   tview.NewBox(),
		align: tview.AlignmentRight,
	}

	return m
}

func (t *textModel) SetText(txt string) {
	t.text = txt
}

func (t *textModel) SetAlign(align tview.Alignment) {
	t.align = align
}

func (t *textModel) SetTextStyle(style tcell.Style) {
	t.textStyle = style
	t.Box.SetBackgroundColor(style.GetBackground())
}

func (t *textModel) Update(tview.Msg) tview.Cmd { return nil }

func (t *textModel) View(screen tview.Screen) {
	if t.w <= 0 {
		return
	}

	t.Box.View(screen)
	tview.PrintWithStyle(screen, t.text, t.x, t.y, t.w, t.align, t.textStyle)
}

func (t *textModel) Rect() (x int, y int, width int, height int) {
	return t.x, t.y, t.w, t.h
}

func (t *textModel) SetRect(x, y, w, h int) {
	t.Box.SetRect(x, y, w, h)
	t.x, t.y, t.w, t.h = x, y, w, h
}

func (t *textModel) Height(int) int {
	return 1
}
