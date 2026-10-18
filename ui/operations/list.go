package operations

import (
	"github.com/ayn2op/tview"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui/widgets/card"
)

type listContent struct {
	cursor  int
	focused bool

	m *Model
}

func newListContent(m *Model, cursor int, focused bool) listContent {
	return listContent{cursor, focused, m}
}

func (l *listContent) Items(index int) tview.Widget {
	info := l.m.ord[index]
	style := theme.Current().Global

	borders := tview.BorderSetRound()
	borderStyle := style
	if l.cursor == index && l.focused {
		borderStyle = theme.Current().Operations.Selection
		borders = tview.BorderSetThick()
	}

	return card.New(info.cc).
		Borders(tview.BordersAll).
		BorderStyle(borderStyle).
		BorderSet(borders).
		Background(style.GetBackground())
}
