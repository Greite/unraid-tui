package syslog

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/url"
	"os/exec"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Greite/unraid-tui/internal/api"
	"github.com/Greite/unraid-tui/internal/i18n"
	"github.com/Greite/unraid-tui/internal/tui/common"
)

type syslogMsg struct {
	logs string
	err  error
}

type tickMsg struct{}

type Model struct {
	client  api.UnraidClient
	logs    string
	offset  int
	follow  bool
	err     error
	width   int
	height  int
	loading bool
}

func New(client api.UnraidClient) Model {
	return Model{
		client:  client,
		follow:  true,
		loading: true,
	}
}

func (m Model) Init() tea.Cmd {
	host := extractHost(m.client.ServerURL())
	return doFetchSyslog(host)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyPressMsg:
		maxOffset := m.maxOffset()
		switch msg.String() {
		case "up", "k":
			m.follow = false
			if m.offset > 0 {
				m.offset--
			}
		case "down", "j":
			if m.offset < maxOffset {
				m.offset++
			}
			if m.offset >= maxOffset {
				m.follow = true
			}
		case "g":
			m.follow = false
			m.offset = 0
		case "G":
			m.follow = true
			m.offset = maxOffset
		case "f":
			m.follow = !m.follow
			if m.follow {
				m.offset = maxOffset
			}
		}
		return m, nil

	case syslogMsg:
		m.loading = false
		if msg.err != nil {
			slog.Warn("syslog fetch failed", "error", msg.err)
			m.err = msg.err
		} else {
			m.logs = msg.logs
			m.err = nil
			if m.follow {
				m.offset = m.maxOffset()
			}
		}
		return m, scheduleRefresh()

	case tickMsg:
		host := extractHost(m.client.ServerURL())
		return m, doFetchSyslog(host)
	}
	return m, nil
}

func (m Model) View() string {
	if m.loading && m.logs == "" {
		return "\n  " + i18n.T("loading")
	}

	var s strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(common.ColorPrimary)
	followIndicator := ""
	if m.follow {
		followIndicator = common.StyleSuccess.Render("  ● " + i18n.T("follow_on"))
	} else {
		followIndicator = common.StyleSubtle.Render("  ○ " + i18n.T("follow_off"))
	}
	s.WriteString("\n" + titleStyle.Render("  Syslog") + followIndicator + "\n")
	s.WriteString(common.StyleSubtle.Render("  "+strings.Repeat("─", 40)) + "\n")

	if m.err != nil && m.logs == "" {
		s.WriteString("\n  " + common.StyleError.Render(m.err.Error()) + "\n")
		return s.String()
	}

	lines := strings.Split(m.logs, "\n")
	visible := m.visibleLines()
	start := m.offset
	if start > len(lines) {
		start = len(lines)
	}
	end := start + visible
	if end > len(lines) {
		end = len(lines)
	}

	lineNoStyle := lipgloss.NewStyle().Foreground(common.ColorMuted)
	for i, line := range lines[start:end] {
		lineNo := lineNoStyle.Render(fmt.Sprintf(" %4d ", start+i+1))
		s.WriteString(lineNo + line + "\n")
	}

	total := len(lines)
	pos := ""
	if total > 0 {
		pct := 0
		if maxOff := total - visible; maxOff > 0 {
			pct = m.offset * 100 / maxOff
		}
		pos = fmt.Sprintf("  %s %d-%d / %d (%d%%)", i18n.T("line"), start+1, end, total, pct)
	}

	s.WriteString("\n" + common.StyleSubtle.Render(pos+"  │  ↑/↓: "+i18n.T("scroll")+"  │  f: "+i18n.T("follow")+"  │  g/G: "+i18n.T("start_end")) + "\n")
	return s.String()
}

func (m Model) Refresh() tea.Cmd {
	host := extractHost(m.client.ServerURL())
	return doFetchSyslog(host)
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m Model) visibleLines() int {
	v := m.height - 6
	if v < 5 {
		v = 5
	}
	return v
}

func (m Model) maxOffset() int {
	lines := strings.Split(m.logs, "\n")
	max := len(lines) - m.visibleLines()
	if max < 0 {
		max = 0
	}
	return max
}

func doFetchSyslog(host string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("ssh",
			"-o", "StrictHostKeyChecking=no",
			"-o", "ConnectTimeout=5",
			"root@"+host,
			"tail -500 /var/log/syslog",
		)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			errMsg := strings.TrimSpace(stderr.String())
			if errMsg == "" {
				errMsg = err.Error()
			}
			return syslogMsg{err: fmt.Errorf("ssh: %s", errMsg)}
		}
		return syslogMsg{logs: stdout.String()}
	}
}

func scheduleRefresh() tea.Cmd {
	return tea.Tick(3*time.Second, func(_ time.Time) tea.Msg {
		return tickMsg{}
	})
}

func extractHost(serverURL string) string {
	parsed, err := url.Parse(serverURL)
	if err != nil {
		return serverURL
	}
	return parsed.Hostname()
}
