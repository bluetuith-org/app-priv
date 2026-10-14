package ui

import (
	"log/slog"
	"sync"
)

var log *logHandler

func init() {
	log = &logHandler{}
	slog.SetDefault(slog.New(slog.DiscardHandler))
}

// logHandler describes a logging handler register for slog.
type logHandler struct {
	handlers []slog.Handler
	started  bool

	sync.Mutex
}

// AddLogHandler adds a slog handler.
func AddLogHandler(handler slog.Handler) {
	l := log

	l.Lock()
	defer l.Unlock()

	if l.started {
		return
	}

	l.handlers = append(l.handlers, handler)
}

// StartLogHandler sets the slog handlers. After Start is called, no
// handlers can be added.
func StartLogHandler() {
	l := log

	l.Lock()
	defer l.Unlock()

	slog.SetDefault(slog.New(slog.NewMultiHandler(
		l.handlers...,
	)))

	l.started = true
	l.handlers = nil
}
