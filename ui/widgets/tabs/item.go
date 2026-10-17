package tabs

import "github.com/rockorager/go-uucode"

// Item represents a tab item.
type Item struct {
	label            string
	width, prevWidth int
}

// NewItem creates a new tab item.
func NewItem(icon, title string, prevWidth int) Item {
	n := " " + icon + " " + title + " "
	width := uucode.StringWidth(n)

	return Item{
		label:     n,
		width:     width,
		prevWidth: prevWidth,
	}
}
