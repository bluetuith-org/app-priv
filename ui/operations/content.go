package operations

import (
	"time"

	"github.com/ayn2op/tview/text"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui/widgets/card"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
)

// Content describes the text that is printed in a single operation status card.
type Content interface {
	// OperationDescription returns the current operation's static description.
	OperationDescription() []text.Span
}

type content struct {
	content Content

	id     string
	header string

	stage operationStage
	err   error
	msg   string

	cc *card.Content
}

func newContent(p Content) *content {
	c := &content{content: p}
	cc := card.NewCardContent(c)

	c.cc = cc

	return c
}

func (c *content) setID(id string) {
	c.id = id
}

func (c *content) setHeader(header string) {
	c.header = header
}

func (c *content) fromUpdateMsg(msg opUpdateMsg) {
	if c.id != msg.id || msg.message == "" && msg.err == nil {
		return
	}

	c.stage = msg.stage
	c.err = msg.err
	c.msg = msg.message

	c.cc.RefreshMessage()
}

func (c *content) Headers() (left text.Span, right text.Span) {
	hdr := stext.NewSpan(c.header, theme.Current().Operations.OperationNumber)
	tm := stext.NewSpan(time.Now().Format("02/01/2006 15:04:05"), theme.Current().Operations.Time)

	return hdr, tm
}

func (c *content) Content() []text.Span {
	return c.content.OperationDescription()
}

func (c *content) Message() (marker text.Span, msg text.Span) {
	var msgSeg text.Span

	switch c.stage {
	case operationError:
		msgSeg = stext.NewSpan(c.err.Error(), theme.Current().Operations.Error)

	default:
		msgSeg = stext.NewSpan(c.msg, theme.Current().Operations.Message)
	}

	status := stext.NewSpan(c.stage.Format())

	return status, msgSeg
}

var _ card.ContentBuilder = (*content)(nil)
