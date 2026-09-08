package ringbuffer

import (
	"iter"
)

// ValueFreer describes a value which will free its internal state, once disposed.
type ValueFreer interface {
	Free()
	Size() int
}

// RingBuffer holds a simple ringbuffer implementation.
// This was adapted from: https://github.com/chronohq/ringslice.
type RingBuffer[T ValueFreer] struct {
	buffer []T

	count, index int
	closed       bool

	totalLen int
}

// New returns a new ring buffer.
func New[T ValueFreer](bufLength int) *RingBuffer[T] {
	lb := &RingBuffer[T]{
		buffer: make([]T, max(1, bufLength)),
	}

	return lb
}

// Len returns the total length of the buffer.
func (l *RingBuffer[T]) Len() int {
	return l.count
}

// Add adds a value to the buffer.
func (l *RingBuffer[T]) Add(val T) {
	if l.closed {
		return
	}

	buflen := len(l.buffer)

	if l.count >= buflen {
		v := l.buffer[l.index]
		l.totalLen -= v.Size()
		v.Free()
	}

	l.buffer[l.index] = val
	l.totalLen += val.Size()
	l.index = (l.index + 1) % buflen

	if l.count < buflen {
		l.count++
	}
}

// TotalContentLen returns the total size of all items.
func (l *RingBuffer[T]) TotalContentLen() int {
	return l.totalLen
}

// Peek returns the latest entry in the buffer.
func (l *RingBuffer[T]) Peek() (T, bool) {
	if l.count == 0 {
		var val T
		return val, false
	}

	buflen := len(l.buffer)
	index := ((l.index - 1) + buflen) % buflen

	return l.buffer[index], true
}

// IterAsc iterates over the buffer in ascending order.
func (l *RingBuffer[T]) IterAsc() iter.Seq[T] {
	return func(yield func(T) bool) {
		startIndex := l.startIndex()
		buflen := len(l.buffer)

		for i := 0; i < l.count; i++ {
			index := (startIndex + i) % buflen
			if !yield(l.buffer[index]) {
				return
			}
		}
	}
}

// IterDesc iterates over the buffer in descending order.
func (l *RingBuffer[T]) IterDesc() iter.Seq[T] {
	return func(yield func(T) bool) {
		startIndex := l.startIndex()
		buflen := len(l.buffer)

		for i := l.count - 1; i >= 0; i-- {
			index := (startIndex + i) % buflen
			if !yield(l.buffer[index]) {
				return
			}
		}
	}
}

func (l *RingBuffer[T]) startIndex() int {
	if l.count == len(l.buffer) {
		return l.index
	}

	return 0
}
