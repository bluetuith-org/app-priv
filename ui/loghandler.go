package ui

import (
	"log/slog"
	"slices"
	"sync"
)

var log *LogHandler

func init() {
	log = &LogHandler{}
	slog.SetDefault(slog.New(slog.DiscardHandler))
}

// LogHandler describes a logging handler register for slog.
type LogHandler struct {
	handlers []slog.Handler
	started  bool

	sync.Mutex
}

// GetLogger gets the logging handler register.
func GetLogger() *LogHandler {
	return log
}

// AddHandler adds a slog handler.
func (l *LogHandler) AddHandler(handler slog.Handler) {
	l.Lock()
	defer l.Unlock()

	if l.started {
		return
	}

	l.handlers = append(l.handlers, handler)
}

// Start sets the slog handlers. After Start is called, no
// handlers can be added.
func (l *LogHandler) Start() {
	l.Lock()
	defer l.Unlock()

	slog.SetDefault(slog.New(slog.NewMultiHandler(
		slices.Clone(l.handlers)...,
	)))

	l.started = true
	l.handlers = nil
}
