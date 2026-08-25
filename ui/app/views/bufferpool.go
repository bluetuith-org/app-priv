package views

import (
	"strings"
	"sync"
)

var _bufferPool = sync.Pool{
	New: func() any {
		return new(strings.Builder)
	},
}

func useStringBuffer(growCapacity int, fn func(b *strings.Builder)) string {
	buf := strings.Builder{}

	if growCapacity > 0 {
		buf.Grow(growCapacity)
	}

	fn(&buf)

	return buf.String()
}
