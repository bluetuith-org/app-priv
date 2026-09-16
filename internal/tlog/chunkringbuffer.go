package tlog

import (
	"slices"
	"unsafe"
)

//revive:disable

const maxBufSize = (24 * 1024)

type ChunkRingRbuffer struct {
	buf []byte
	p   []int

	writeIndex int

	count, index int
	bufSize      int
	length       int
}

func NewChunkRingBuffer(ringSize, bufSize int) *ChunkRingRbuffer {
	bufSize = nextPowerOf2(bufSize)
	return &ChunkRingRbuffer{
		buf:     make([]byte, 0, bufSize),
		p:       make([]int, ringSize),
		bufSize: bufSize,
		length:  ringSize,
	}
}

func (c *ChunkRingRbuffer) Append(buf []byte) {
	total := c.length

	if c.count >= total {
		if c.index == 0 {
			c.writeIndex = 0
		}

		currLength := c.p[c.index]

		c.p[c.index] = len(buf)
		c.buf = slices.Replace(c.buf, c.writeIndex, currLength+c.writeIndex, buf...)

		c.writeIndex += len(buf)
	} else {
		c.buf = slices.Insert(c.buf, len(c.buf), buf...)
		c.p[c.index] = len(buf)
	}

	if cap(c.buf) >= maxBufSize {
		c.buf = slices.Clip(c.buf)
	}

	c.index = (c.index + 1) % total
	if c.count < total {
		c.count++
	}
}

func (c *ChunkRingRbuffer) Len() int {
	return len(c.buf)
}

func (c *ChunkRingRbuffer) Cap() int {
	return cap(c.buf)
}

func (c *ChunkRingRbuffer) Strings() (string, string) {
	return ustr(c.buf[c.writeIndex:]), ustr(c.buf[:c.writeIndex])
}

func ustr(buf []byte) string {
	return unsafe.String(unsafe.SliceData(buf), len(buf))
}

// nextPowerOf2 returns the next power of 2 greater than or equal to n
//
//go:inline
func nextPowerOf2(n int) int {
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	n++
	return n
}
