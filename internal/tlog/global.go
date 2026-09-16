package tlog

import (
	"sync/atomic"
)

// global stores the global logger.
var global atomic.Pointer[Log]

func init() {
	l := NewLog(nil, 1000)
	global.Store(l)
}

// G gets the logger.
func G() *Log {
	return global.Load()
}

// SetLogger sets the logger.
func SetLogger(l *Log) {
	global.Store(l)
}
