package views

import (
	"strings"
)

func useStringBuffer(growCapacity int, fn func(b *strings.Builder)) string {
	buf := strings.Builder{}

	if growCapacity > 0 {
		buf.Grow(growCapacity)
	}

	fn(&buf)

	return buf.String()
}
