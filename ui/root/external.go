package root

import (
	"sync/atomic"

	"github.com/gdamore/tcell/v3"
)

type evState struct {
	ch     chan tcell.Event
	closed chan struct{}

	done atomic.Bool
}

func newEvState(screen tcell.Screen) *evState {
	return &evState{
		ch:     screen.EventQ(),
		closed: make(chan struct{}),
	}
}

func (e *evState) sendMsg(msg tcell.Event) {
	if e.done.Load() {
		return
	}

	select {
	case <-e.closed:
		e.done.Store(true)

	case e.ch <- msg:
	}
}

func (e *evState) close() {
	close(e.closed)
}
