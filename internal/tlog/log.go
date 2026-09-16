package tlog

import (
	"context"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/bluetuith-org/bluetuith/internal/buffers"
	"github.com/spf13/cast"
)

//revive:disable

type Log struct {
	rb *ChunkRingRbuffer
	w  LogWriter

	currBuf int

	mu sync.Mutex
}

func NewLog(w LogWriter, maxLines int) *Log {
	return &Log{
		rb: NewChunkRingBuffer(maxLines, 1024),
		w:  w,
	}
}

func (l *Log) LogAttrs(ctx context.Context, t time.Time, level slog.Level, msg string, attrs ...slog.Attr) {
	if l.w == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	b := buffers.NewBuffer()
	defer b.Free()

	l.w.Reset()
	defer l.w.Finish()

	l.w.AppendTime(b, time.Now())
	l.w.AppendLevel(b, level)
	l.w.AppendMsg(b, msg)

	l.w.StartKVWrite(b)
	for _, a := range attrs {
		l.w.AppendKeyValue(b, a)
	}
	l.w.EndKVWrite(b)

	*b = append(*b, "\n"...)
	l.rb.Append(*b)
}

func (l *Log) Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	l.LogAttrs(ctx, time.Now(), slog.LevelInfo, msg, attrs...)
}

func (l *Log) Error(ctx context.Context, msg string, attrs ...slog.Attr) {
	l.LogAttrs(ctx, time.Now(), slog.LevelError, msg, attrs...)
}

func (l *Log) Debug(ctx context.Context, msg string, attrs ...slog.Attr) {
	l.LogAttrs(ctx, time.Now(), slog.LevelDebug, msg, attrs...)
}

func (l *Log) GetContent(wrap func(string) string) string {
	l.mu.Lock()
	defer l.mu.Unlock()

	start, end := l.rb.Strings()

	return start + end
}

type LogWriter interface {
	Reset()
	Finish()

	AppendTime(buf *buffers.Buffer, n time.Time)
	AppendLevel(b *buffers.Buffer, level slog.Level)
	AppendKeyValue(b *buffers.Buffer, attr slog.Attr)
	AppendMsg(by *buffers.Buffer, m string)

	StartKVWrite(b *buffers.Buffer)
	EndKVWrite(b *buffers.Buffer)
}

func AppendAttrToBuffer(by *buffers.Buffer, a slog.Attr) {
	a.Value = a.Value.Resolve()

	*by = append(*by, a.Key...)
	*by = append(*by, kvSep...)

	switch a.Value.Kind() {
	case slog.KindBool:
		*by = strconv.AppendBool(*by, a.Value.Bool())

	case slog.KindInt64:
		*by = strconv.AppendInt(*by, a.Value.Int64(), 10)

	case slog.KindUint64:
		*by = strconv.AppendUint(*by, a.Value.Uint64(), 10)

	case slog.KindDuration:
		var arr [32]byte
		n := formatDuration(a.Value.Duration(), &arr)
		*by = append(*by, arr[n:]...)

	case slog.KindTime:
		*by = a.Value.Time().AppendFormat(*by, timeFmt)

	case slog.KindFloat64:
		*by = strconv.AppendFloat(*by, a.Value.Float64(), 'f', -1, 64)

	case slog.KindString:
		*by = append(*by, a.Value.String()...)

	case slog.KindAny:
		*by = append(*by, cast.ToString(a.Value.Any())...)
	}
}
