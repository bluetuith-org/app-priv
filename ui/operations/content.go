package operations

import (
	"time"

	"github.com/ayn2op/tview/richtext"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui/widgets/card"
)

type content struct {
	card.ContentProvider

	id     string
	header string

	stage operationStage
	err   error
	msg   string

	cc *card.Content
}

func newContent(p card.ContentProvider) *content {
	c := &content{ContentProvider: p}
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

func (c *content) Headers() (left richtext.Segment, right richtext.Segment) {
	hdr := richtext.NewSegment(c.header, theme.Current().Operations.OperationNumber)
	tm := richtext.NewSegment(time.Now().Format("02/01/2006 15:04:05"), theme.Current().Operations.Time)

	return hdr, tm
}

func (c *content) Message() (marker richtext.Segment, msg richtext.Segment) {
	var msgSeg richtext.Segment

	switch c.stage {
	case operationError:
		msgSeg = richtext.NewSegment(c.err.Error(), theme.Current().Operations.Error)

	default:
		msgSeg = richtext.NewSegment(c.msg, theme.Current().Operations.Message)
	}

	status := richtext.NewSegment(c.stage.Format())

	return status, msgSeg
}
