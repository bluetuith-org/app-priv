package views

import (
	"log/slog"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetuith/internal/tlog"
	"github.com/bluetuith-org/bluetuith/ui/theme"
)

type logView struct {
	focused bool

	width, height int
	vp            viewport.Model

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
	l.vp = viewport.New()
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

	default:
	}

	return l, nil
}

// View renders the program's UI, which can be a string or a [Layer]. The
// view is rendered after every Update.
func (l *logView) View() tea.View {
	text := useStringBuffer(1024, func(b *strings.Builder) {
		for v := range tlog.IterateAsc() {
			renderLogLine(l.width, b, &v)
		}
	})

	l.vp.Style = theme.Current().Log.Style.
		Align(lipgloss.Left).
		Border(lipgloss.RoundedBorder())

	l.vp.SetWidth(l.width)
	l.vp.SetHeight(l.height)
	l.vp.SetContent(text)

	return tea.NewView(l.vp.View())
}

func renderLogLine(w int, b *strings.Builder, l *tlog.LogData) {
	level := ""
	var levelStyle lipgloss.Style

	switch l.Level {
	case slog.LevelInfo:
		level = "INFO"
		levelStyle = theme.Current().Log.Info

	case slog.LevelDebug:
		level = "DEBUG"
		levelStyle = theme.Current().Log.Debug

	case slog.LevelError:
		level = "ERROR"
		levelStyle = theme.Current().Log.Error
	}

	currStyle := theme.Current().Log.Style

	levelStr := levelStyle.Render(level)
	logStr := currStyle.Render(" " + l.String())
	width := abs(w - (lipgloss.Width(levelStr) + lipgloss.Width(logStr) + 1))
	pad := currStyle.Render(strings.Repeat(" ", width))

	b.WriteString(levelStr)
	b.WriteString(logStr)
	b.WriteString(pad)
	b.WriteString("\n")
}

func abs[T ~int | ~uint](x T) T {
	var zero T
	if x < zero {
		return -x
	}

	return x
}
