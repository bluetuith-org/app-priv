package logs

import (
	"iter"
)

// ringBuffer holds a simple ringbuffer implementation.
// This was adapted from: https://github.com/chronohq/ringslice.
type ringBuffer[T any] struct {
	buffer []T

	count, index int
	totalLen     int
}

// newRingBuffer returns a new ring buffer.
func newRingBuffer[T any](bufLength int) *ringBuffer[T] {
	lb := &ringBuffer[T]{
		buffer: make([]T, max(1, bufLength)),
	}

	return lb
}

// Len returns the total length of the buffer.
func (l *ringBuffer[T]) Len() int {
	return l.count
}

// Add adds a value to the buffer.
func (l *ringBuffer[T]) Add(val T) {
	buflen := len(l.buffer)

	l.buffer[l.index] = val
	l.index = (l.index + 1) % buflen

	if l.count < buflen {
		l.count++
	}
}

// TotalContentLen returns the total size of all items.
func (l *ringBuffer[T]) TotalContentLen() int {
	return l.totalLen
}

// Peek returns the latest entry in the buffer.
func (l *ringBuffer[T]) Peek() (T, bool) {
	if l.count == 0 {
		var val T
		return val, false
	}

	buflen := len(l.buffer)
	index := ((l.index - 1) + buflen) % buflen

	return l.buffer[index], true
}

// Get gets an item from the buffer.
func (l *ringBuffer[T]) Get(index int) T {
	idx := (l.startIndex() + index) % len(l.buffer)
	return l.buffer[idx]
}

// IterAsc iterates over the buffer in ascending order.
func (l *ringBuffer[T]) IterAsc() iter.Seq[T] {
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
func (l *ringBuffer[T]) IterDesc() iter.Seq[T] {
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

func (l *ringBuffer[T]) startIndex() int {
	if l.count == len(l.buffer) {
		return l.index
	}

	return 0
}
