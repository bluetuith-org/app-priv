package tlog

import (
	"iter"
	"log/slog"
	"sync/atomic"
)

// global stores the global logger.
var global atomic.Pointer[Logger]

func init() {
	l := New(1000)
	global.Store(l)
	slog.SetDefault(l.logger)
}

func get() *Logger {
	return global.Load()
}

// SetLogger sets the logger.
func SetLogger(l *Logger) {
	global.Store(l)
}

// Total returns the total number of entries in the log buffer.
func Total() int {
	return get().Total()
}

// TotalContentLen returns the complete content size.
func TotalContentLen() int {
	return get().TotalContentLen()
}

// Peek returns the latest log entry in the buffer.
func Peek() (LogData, bool) {
	return get().Peek()
}

// IterateAsc iterates over the current log buffer in ascending order.
func IterateAsc() iter.Seq[LogData] {
	return get().IterateAsc()
}

// IterateDesc iterates over the current log buffer in descending order.
func IterateDesc() iter.Seq[LogData] {
	return get().IterateDesc()
}
