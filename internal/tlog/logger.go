package tlog

import (
	"context"
	"io"
	"iter"
	"log/slog"
	"time"

	"github.com/bluetuith-org/bluetuith/internal/ringbuffer"
)

// Logger provides a logger and a [slog.Handler].
type Logger struct {
	rb *ringbuffer.RingBuffer[LogData]

	logger *slog.Logger
}

// New returns a new logger.
func New(maxCapacity int) *Logger {
	l := &Logger{
		rb: ringbuffer.New[LogData](maxCapacity),
	}

	l.logger = slog.New(l)

	return l
}

// NewWithHandler returns a new logger after attaching a custom handler to [slog.Logger].
func NewWithHandler(maxCapacity int, handler slog.Handler) *Logger {
	l := &Logger{
		rb: ringbuffer.New[LogData](maxCapacity),
	}

	l.logger = slog.New(slog.NewMultiHandler(l, handler))

	return l
}

// Total returns the total number of entries in the log buffer.
func (l *Logger) Total() int {
	return l.rb.Len()
}

// TotalContentLen returns the complete content size.
func (l *Logger) TotalContentLen() int {
	return l.rb.TotalContentLen()
}

// Peek returns the latest log entry in the buffer.
func (l *Logger) Peek() (LogData, bool) {
	return l.rb.Peek()
}

// IterateAsc iterates over the current log buffer in ascending order.
func (l *Logger) IterateAsc() iter.Seq[LogData] {
	return l.rb.IterAsc()
}

// IterateDesc iterates over the current log buffer in descending order.
func (l *Logger) IterateDesc() iter.Seq[LogData] {
	return l.rb.IterDesc()
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
// It is called early, before any arguments are processed,
// to save effort if the log event should be discarded.
// If called from a Logger method, the first argument is the context
// passed to that method, or context.Background() if nil was passed
// or the method does not take a context.
// The context is passed so Enabled can use its values
// to make a decision.
func (l *Logger) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

// Handle handles the Record.
// It will only be called when Enabled returns true.
// The Context argument is as for Enabled.
// It is present solely to provide Handlers access to the context's values.
// Canceling the context should not affect record processing.
// (Among other things, log messages may be necessary to debug a
// cancellation-related problem.)
//
// Handle methods that produce output should observe the following rules:
//   - If r.Time is the zero time, ignore the time.
//   - If r.PC is zero, ignore it.
//   - Attr's values should be resolved.
//   - If an Attr's key and value are both the zero value, ignore the Attr.
//     This can be tested with attr.Equal(Attr{}).
//   - If a group's key is empty, inline the group's Attrs.
//   - If a group has no Attrs (even if it has a non-empty key),
//     ignore it.
//
// [Logger] discards any errors from Handle. Wrap the Handle method to
// process any errors from Handlers.
func (l *Logger) Handle(_ context.Context, r slog.Record) error {
	l.rb.Add(newLogData(r.Time, r.Level, getAttrBuf(r)))

	return nil
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
// The Handler owns the slice: it may retain, modify or discard it.
func (l *Logger) WithAttrs(_ []slog.Attr) slog.Handler {
	return l
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
// The keys of all subsequent attributes, whether added by With or in a
// Record, should be qualified by the sequence of group names.
//
// How this qualification happens is up to the Handler, so long as
// this Handler's attribute keys differ from those of another Handler
// with a different sequence of group names.
//
// A Handler should treat WithGroup as starting a Group of Attrs that ends
// at the end of the log event. That is,
//
//	logger.WithGroup("s").LogAttrs(ctx, level, msg, slog.Int("a", 1), slog.Int("b", 2))
//
// should behave like
//
//	logger.LogAttrs(ctx, level, msg, slog.Group("s", slog.Int("a", 1), slog.Int("b", 2)))
//
// If the name is empty, WithGroup returns the receiver.
func (l *Logger) WithGroup(_ string) slog.Handler {
	return l
}

// LogData defines the data to be logged
type LogData struct {
	Time  time.Time
	Level slog.Level

	b   *Buffer
	len int
}

func newLogData(t time.Time, level slog.Level, b *Buffer) LogData {
	l := LogData{}

	l.Time = t
	l.Level = level

	l.b = b
	l.len = b.Len()

	return l
}

// WriteBuffer writes to a writer
func (l *LogData) WriteBuffer(w io.Writer) {
	if l.len == 0 {
		return
	}

	l.b.WriteBuffer(w)
}

// Size returns the length of the data.
func (l LogData) Size() int {
	return l.len
}

func (l LogData) String() string {
	return l.b.String()
}

// Free frees the underlying byte slice.
func (l LogData) Free() {
	l.b.Free()
}
