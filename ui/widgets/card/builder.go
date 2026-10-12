package card

import "github.com/ayn2op/tview/richtext"

// ContentBuilder describes a provider which builds and returns
// headers, the text body and the message content for the card.
type ContentBuilder interface {
	ContentProvider
	ContentPrefixer
}

// ContentPrefixer describes a provider that returns styled segments
// for card headers and message content.
type ContentPrefixer interface {
	Headers() (left richtext.Segment, right richtext.Segment)
	Message() (marker richtext.Segment, msg richtext.Segment)
}

// ContentProvider describes a provider that returns styled text
// for the card's text body.
type ContentProvider interface {
	Content() richtext.Text
}
