package tlog

import (
	"io"
	"slices"
	"sync"
	"unsafe"
)

// Ref: https://cs.opensource.google/go/go/+/refs/tags/go1.27.1:src/log/slog/internal/buffer/buffer.go
//
//revive:disable

type Buffer []byte

// Having an initial size gives a dramatic speedup.
var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 1024)
		return (*Buffer)(&b)
	},
}

func NewBuffer() *Buffer {
	return bufPool.Get().(*Buffer)
}

func (b *Buffer) Free() {
	// To reduce peak allocation, return only smaller buffers to the pool.
	const maxBufferSize = 16 << 10
	if cap(*b) <= maxBufferSize {
		*b = (*b)[:0]
		bufPool.Put(b)
	}
}

func (b *Buffer) Reset() {
	b.SetLen(0)
}

func (b *Buffer) Buf() []byte {
	return *b
}

func (b *Buffer) Clip() *Buffer {
	*b = slices.Clip(*b)
	return b
}

func (b *Buffer) WriteBuffer(w io.Writer) (int, error) {
	return w.Write([]byte(*b))
}

func (b *Buffer) Write(p []byte) (int, error) {
	*b = append(*b, p...)
	return len(p), nil
}

func (b *Buffer) WriteString(s string) (int, error) {
	*b = append(*b, s...)
	return len(s), nil
}

func (b *Buffer) WriteByte(c byte) error {
	*b = append(*b, c)
	return nil
}

func (b *Buffer) String() string {
	return unsafe.String(unsafe.SliceData(*b), len(*b))
}

func (b *Buffer) Len() int {
	return len(*b)
}

func (b *Buffer) SetLen(n int) {
	*b = (*b)[:n]
}
