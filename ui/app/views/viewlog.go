package views

import (
	"bytes"
	"log/slog"
	"time"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetuith/internal/buffers"
	"github.com/bluetuith-org/bluetuith/internal/tlog"
	"github.com/bluetuith-org/bluetuith/ui/theme"
)

type logView struct {
	focused bool

	width, height int
	vp            viewport.Model

	w  *styledLogWriter
	ch chan string

	v rootView
}

func (l *logView) Title() string {
	return "Log"
}

func (l *logView) Icon() *theme.IconVariant {
	return theme.Icons().Log
}

// ViewID returns the view's ID.
func (l *logView) ViewID() viewID {
	return viewIDLog
}

// InitializeView initializes the view.
func (l *logView) InitializeView(_ *appfeatures.FeatureSet) (inited bool, err error) {
	l.w = newStyledLogWriter(l.v)
	tlog.SetLogger(tlog.NewLog(l.w, 1000))

	l.vp = viewport.New()
	l.vp.SoftWrap = false

	return true, nil
}

// SetRootView sets the root view.
func (l *logView) SetRootView(v rootView) {
	l.v = v
}

// AttachToTabView attaches this view to the tabbed view.
func (l *logView) AttachToTabView() (tabSection, bool) {
	return l, true
}

// Resize resizes the view.
func (l *logView) Resize(width int, height int) tea.WindowSizeMsg {
	l.width = width
	l.height = height

	return tea.WindowSizeMsg{Width: l.width, Height: l.height}
}

// SetFocus sets whether the view is currently focused.
func (l *logView) SetFocus(focused bool) {
	l.focused = focused
}

// GetFocus gets whether the view is currently focused.
func (l *logView) GetFocus() bool {
	return l.focused
}

// HandleRouterMsg handles the routed message.
func (l *logView) HandleRouterMsg(m routerMsg) tea.Cmd {
	return handleRouterMsg(l, m)
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (l *logView) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (l *logView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		msg = l.Resize(m.Width, m.Height)

		l.vp.SetWidth(l.width)
		l.vp.SetHeight(l.height)

	default:
	}

	return l, nil
}

// View renders the program's UI, which can be a string or a [Layer]. The
// view is rendered after every Update.
func (l *logView) View() tea.View {
	currViewStyle := theme.Current().Log.Style
	borderStyle := theme.Current().Border

	l.vp.Style = currViewStyle.
		Height(l.vp.Height()).
		Width(l.vp.Width()).
		Align(lipgloss.Left).
		Border(lipgloss.RoundedBorder()).
		BorderBackground(borderStyle.GetBackground()).
		BorderForeground(borderStyle.GetForeground())

	content := tlog.G().GetContent(func(s string) string {
		return lipgloss.PlaceHorizontal(
			l.vp.Width(), lipgloss.Top,
			s,
			lipgloss.WithWhitespaceStyle(theme.Current().Log.Style),
		)
	})

	l.vp.SetContent(content)
	l.vp.GotoBottom()

	return tea.NewView(currViewStyle.Render(l.vp.View()))
}

type escPos struct {
	start, end []byte
}

type styledLogWriter struct {
	v rootView

	Log struct {
		Time escPos

		Level struct {
			Info  escPos
			Debug escPos
			Error escPos
		}

		Msg   escPos
		Style escPos
	}

	Status struct{}
}

func newStyledLogWriter(v rootView) *styledLogWriter {
	s := &styledLogWriter{
		v: v,
	}

	s.Log.Time.start, s.Log.Time.end, _ = s.splitStyle(theme.Current().Log.Time)

	s.Log.Level.Info.start, s.Log.Level.Info.end, _ = s.splitStyle(theme.Current().Log.Info)
	s.Log.Level.Debug.start, s.Log.Level.Debug.end, _ = s.splitStyle(theme.Current().Log.Debug)
	s.Log.Level.Error.start, s.Log.Level.Error.end, _ = s.splitStyle(theme.Current().Log.Error)

	s.Log.Msg.start, s.Log.Msg.end, _ = s.splitStyle(theme.Current().Log.Style)
	s.Log.Style.start, s.Log.Style.end, _ = s.splitStyle(theme.Current().Log.Style)

	return s
}

func (s *styledLogWriter) Reset() {
}

func (s *styledLogWriter) Finish() {
}

func (s *styledLogWriter) AppendTime(b *buffers.Buffer, n time.Time) {
	start, end := s.Log.Time.start, s.Log.Time.end
	*b = append(*b, start...)
	*b = n.AppendFormat(*b, "15:04:05")
	*b = append(*b, " "...)
	*b = append(*b, end...)
}

func (s *styledLogWriter) AppendLevel(by *buffers.Buffer, l slog.Level) {
	var level string
	var start, end []byte

	switch l.Level() {
	case slog.LevelInfo:
		level = "INFO"
		start, end = s.Log.Level.Info.start, s.Log.Level.Info.end

	case slog.LevelDebug:
		level = "DEBUG"
		start, end = s.Log.Level.Debug.start, s.Log.Level.Debug.end

	case slog.LevelError:
		level = "ERROR"
		start, end = s.Log.Level.Error.start, s.Log.Level.Error.end
	}

	*by = append(*by, start...)
	*by = append(*by, level...)
	*by = append(*by, " "...)
	*by = append(*by, end...)
}

func (s *styledLogWriter) AppendMsg(by *buffers.Buffer, m string) {
	start, end := s.Log.Msg.start, s.Log.Msg.end
	*by = append(*by, start...)
	*by = append(*by, m...)
	*by = append(*by, " "...)
	*by = append(*by, end...)
}

func (s *styledLogWriter) AppendKeyValue(by *buffers.Buffer, a slog.Attr) {
	tlog.AppendAttrToBuffer(by, a)
}

func (s *styledLogWriter) StartKVWrite(by *buffers.Buffer) {
	by.Write(s.Log.Style.start)
}

func (s *styledLogWriter) EndKVWrite(by *buffers.Buffer) {
	by.Write(s.Log.Style.end)
}

func (s *styledLogWriter) splitStyle(style lipgloss.Style) (start []byte, end []byte, ok bool) {
	const _styleSep = "@"
	_styleByteSep := []byte("@")

	str := style.Render(_styleSep)
	start, end, ok = bytes.Cut([]byte(str), _styleByteSep)

	return
}
