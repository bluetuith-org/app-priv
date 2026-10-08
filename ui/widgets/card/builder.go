package card

import "github.com/ayn2op/tview/richtext"

//revive:disable

type ContentBuilder interface {
	ContentProvider
	ContentPrefixer
}

type ContentPrefixer interface {
	Headers() (left richtext.Segment, right richtext.Segment)
	Message() (marker richtext.Segment, msg richtext.Segment)
}

type ContentProvider interface {
	Content() richtext.Text
}
