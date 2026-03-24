package plugins

import (
	"context"
	"log/slog"
	"strings"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Greite/unraid-tui/internal/api"
	"github.com/Greite/unraid-tui/internal/i18n"
	"github.com/Greite/unraid-tui/internal/tui/common"
)

type Model struct {
	client  api.UnraidClient
	plugins []string
	spinner spinner.Model
	loading bool
	err     error
	cursor  int
	offset  int
	width   int
	height  int
}

func New(client api.UnraidClient) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(common.ColorPrimary)
	return Model{
		client:  client,
		spinner: s,
		loading: true,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.fetchPlugins,
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.offset {
					m.offset = m.cursor
				}
			}
		case "down", "j":
			if m.cursor < len(m.plugins)-1 {
				m.cursor++
				visible := m.visibleRows()
				if m.cursor >= m.offset+visible {
					m.offset = m.cursor - visible + 1
				}
			}
		}

	case common.PluginsMsg:
		m.loading = false
		if msg.Err != nil {
			slog.Error("plugins fetch failed", "error", msg.Err)
			m.err = msg.Err
			return m, nil
		}
		m.err = nil
		m.plugins = msg.Plugins
		m.cursor = 0
		m.offset = 0
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.loading {
		return "\n  " + m.spinner.View() + " " + i18n.T("loading_plugins")
	}

	var s strings.Builder

	if m.err != nil {
		s.WriteString("\n  " + common.StyleError.Render("⚠ "+m.err.Error()) + "\n")
		return s.String()
	}

	if len(m.plugins) == 0 {
		s.WriteString("\n  " + common.StyleSubtle.Render(i18n.T("no_plugins")) + "\n")
		return s.String()
	}

	title := common.StyleTitle.Render("  " + i18n.T("plugins_installed"))
	count := common.StyleSubtle.Render(strings.Replace(" (X)", "X", intToStr(len(m.plugins)), 1))
	s.WriteString("\n" + title + count + "\n\n")

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#CC6A1E"))

	visible := m.visibleRows()
	end := m.offset + visible
	if end > len(m.plugins) {
		end = len(m.plugins)
	}

	for idx := m.offset; idx < end; idx++ {
		name := m.plugins[idx]
		// Strip .plg extension for cleaner display
		displayName := strings.TrimSuffix(name, ".plg")
		row := "  ● " + displayName
		if idx == m.cursor {
			s.WriteString(selectedStyle.Render(row) + "\n")
		} else {
			s.WriteString(row + "\n")
		}
	}

	s.WriteString("\n" + common.StyleSubtle.Render("  ↑/↓: "+i18n.T("navigate")+"  │  r: "+i18n.T("refresh")) + "\n")
	return s.String()
}

func (m Model) Refresh() tea.Cmd {
	return m.fetchPlugins
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m Model) visibleRows() int {
	v := m.height - 8
	if v < 5 {
		v = 5
	}
	return v
}

func (m Model) fetchPlugins() tea.Msg {
	plugins, err := m.client.GetInstalledPlugins(context.Background())
	return common.PluginsMsg{Plugins: plugins, Err: err}
}

func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	digits := make([]byte, 0, 10)
	for n > 0 {
		digits = append(digits, byte('0'+n%10))
		n /= 10
	}
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	return string(digits)
}
