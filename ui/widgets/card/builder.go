package card

import "github.com/ayn2op/tview/text"

// ContentBuilder describes a provider which builds and returns
// headers, the text body and the message content for the card.
type ContentBuilder interface {
	ContentProvider
	ContentPrefixer
}

// ContentPrefixer describes a provider that returns styled segments
// for card headers and message content.
type ContentPrefixer interface {
	Headers() (left text.Span, right text.Span)
	Message() (marker text.Span, msg text.Span)
}

// ContentProvider describes a provider that returns styled text
// for the card's text body.
type ContentProvider interface {
	Content() []text.Span
}
