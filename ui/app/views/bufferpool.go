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

func useBuffer(fn func(b *strings.Builder)) string {
	buf := getBuffer()
	defer returnBuffer(buf)

	fn(buf)

	return buf.String()
}

func getBuffer() *strings.Builder {
	buf, ok := _bufferPool.Get().(*strings.Builder)
	if !ok {
		return nil
	}

	return buf
}

func returnBuffer(b *strings.Builder) {
	b.Reset()
	_bufferPool.Put(b)
}
