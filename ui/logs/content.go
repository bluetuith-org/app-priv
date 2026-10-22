package logs

import (
	"log/slog"
	"strconv"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/layout"
	"github.com/ayn2op/tview/text"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
	"github.com/gdamore/tcell/v3"
	"github.com/spf13/cast"
)

type logContent struct {
	cached [][]text.Span
	record slog.Record

	width    int
	rendered bool
}

func newLogContent(orig slog.Record) *logContent {
	return &logContent{record: orig}
}

func (c *logContent) Refresh() *logContent {
	c.rendered = false
	return c
}

// Draw implements [tview.Widget].
func (c *logContent) Draw(screen tview.Screen, area tview.Rectangle) {
	area.Height = c.height()

	for _, line := range c.cached {
		stext.DrawLine(screen, line, area.X, area.Y, area)
		area.Y++
	}
}

// Handle implements [tview.Widget].
func (c *logContent) Handle(msg tview.Msg, _ tview.Rectangle) tview.Msg {
	return msg
}

// Layout implements [tview.Widget].
func (c *logContent) Layout(limits layout.Limits) layout.Size {
	return layout.Sized(limits, layout.Fill, layout.Shrink, func(l layout.Limits) layout.Size {
		h := c.ensureSize(l.Max.Width)
		return layout.Size{Height: h + 1}
	})
}

// Size implements [tview.Widget].
func (c *logContent) Size() (width layout.Length, height layout.Length) {
	return layout.Fill, layout.Shrink
}

func (c *logContent) ensureSize(width int) int {
	if c.width == width && c.rendered {
		return c.height()
	}

	if c.width != width {
		c.rendered = false
	}

	c.render(width)
	c.width = width

	return c.height()
}

func (c *logContent) height() int {
	return len(c.cached)
}

const timeFmt = "2006-01-02 15:04:05"

func (c *logContent) render(width int) {
	c.cached = c.cached[:0]
	c.rendered = true

	var levelStyle tcell.Style
	var levelStr string

	timeStyle := theme.Current().Log.Time
	bgStyle := theme.Current().Global
	attrStyle := theme.Current().Log.Attr
	msgStyle := theme.Current().Log.Message

	switch c.record.Level {
	case slog.LevelInfo:
		levelStr = "INFO"
		levelStyle = theme.Current().Log.Info

	case slog.LevelDebug:
		levelStr = "DEBUG"
		levelStyle = theme.Current().Log.Debug

	case slog.LevelError:
		levelStr = "ERROR"
		levelStyle = theme.Current().Log.Error
	}

	builder := stext.NewBuilder(bgStyle)

	builder.Append(c.record.Time.Format(timeFmt), timeStyle)
	builder.Space()

	builder.Append(levelStr, levelStyle)
	builder.Space()

	builder.Append(c.record.Message, msgStyle)
	builder.Newline()

	c.record.Attrs(func(a slog.Attr) bool {
		k, v := c.getKV(a)

		builder.Space()
		builder.Append(theme.Icons().BottomLeft.String(), attrStyle)
		builder.Space()

		builder.Append(k, attrStyle)
		builder.Append(" = ", attrStyle)
		builder.Append(v, attrStyle)
		builder.Newline()

		return true
	})

	c.cached = text.Lines(builder.Text(), width, text.WrappingGlyph)
}

func (c *logContent) getKV(a slog.Attr) (key string, value string) {
	a.Value = a.Value.Resolve()

	val := ""

	switch a.Value.Kind() {
	case slog.KindBool:
		val = strconv.FormatBool(a.Value.Bool())

	case slog.KindInt64:
		val = strconv.FormatInt(a.Value.Int64(), 10)

	case slog.KindUint64:
		val = strconv.FormatUint(a.Value.Uint64(), 10)

	case slog.KindDuration:
		val = a.Value.Duration().String()

	case slog.KindTime:
		val = a.Value.Time().Format(timeFmt)

	case slog.KindFloat64:
		val = strconv.FormatFloat(a.Value.Float64(), 'f', -1, 64)

	case slog.KindString:
		val = a.Value.String()

	case slog.KindAny:
		val = cast.ToString(a.Value.Any())
	}

	return a.Key, val
}
