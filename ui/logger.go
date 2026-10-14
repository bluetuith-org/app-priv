package ui

import (
	"log/slog"
	"sync/atomic"
)

//revive:disable

var log atomic.Pointer[Logger]

func init() {
	log.Store(&Logger{})
	slog.SetDefault(slog.New(slog.DiscardHandler))
}

// Logger describes a logging utility for the UI only.
type Logger struct {
	handlers []slog.Handler
}

func GetLogger() *Logger {
	return log.Load()
}

func (l *Logger) AddHandler(handler slog.Handler) {
	l.handlers = append(l.handlers, handler)
}

func (l *Logger) Start() {
	slog.SetDefault(slog.New(slog.NewMultiHandler(l.handlers...)))
}
