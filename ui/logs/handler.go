package logs

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/ayn2op/tview/text"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
	"github.com/gdamore/tcell/v3"
)

type handler struct {
	rb *ringBuffer[*logContent]
	rv ui.RootView

	lcontent []text.Span
	latest   time.Time

	sync.Mutex
}

func newHandler(rv ui.RootView, maxSize int) *handler {
	return &handler{
		rb: newRingBuffer[*logContent](maxSize),
		rv: rv,
	}
}

// Enabled implements [slog.Handler].
func (h *handler) Enabled(context.Context, slog.Level) bool {
	return true
}

// Handle implements [slog.Handler].
func (h *handler) Handle(_ context.Context, record slog.Record) error {
	var idx int

	{
		h.Lock()

		idx = h.rb.Len()
		h.rb.Add(newLogContent(record))

		h.Unlock()
	}

	h.rv.SendMsg(ui.ViewIDLog.R(newLogMsg{idx}))

	return nil
}

// WithAttrs implements [slog.Handler].
func (h *handler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

// WithGroup implements [slog.Handler].
func (h *handler) WithGroup(string) slog.Handler {
	return h
}

func (h *handler) count() int {
	h.Lock()
	defer h.Unlock()

	return h.rb.Len()
}

func (h *handler) getLatest() []text.Span {
	h.Lock()
	defer h.Unlock()

	c, ok := h.rb.Peek()
	if !ok || c.record.Time.Equal(h.latest) {
		return h.lcontent
	}

	h.latest = c.record.Time

	if h.lcontent == nil {
		h.lcontent = make([]text.Span, 3)
	}

	var levelStyle tcell.Style
	var levelStr string

	switch c.record.Level {
	case slog.LevelInfo:
		levelStr = " INFO "
		levelStyle = theme.Current().StatusBar.Info

	case slog.LevelDebug:
		levelStr = " DEBUG "
		levelStyle = theme.Current().StatusBar.Debug

	case slog.LevelError:
		levelStr = " ERROR "
		levelStyle = theme.Current().StatusBar.Error
	}

	style := theme.Current().StatusBar.Style

	h.lcontent[0] = stext.NewSpan(levelStr, levelStyle)
	h.lcontent[1] = stext.NewSpan(" ", style)
	h.lcontent[2] = stext.NewSpan(c.record.Message, style)

	return h.lcontent
}

var _ slog.Handler = (*handler)(nil)
