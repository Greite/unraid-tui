package vms

import (
	"context"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/spinner"
	"charm.land/lipgloss/v2"
	"github.com/Greite/unraid-tui/internal/api"
	"github.com/Greite/unraid-tui/internal/model"
	"github.com/Greite/unraid-tui/internal/tui/common"
)

type Model struct {
	client  api.UnraidClient
	vms     []model.VM
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
	return Model{client: client, spinner: s, loading: true}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchVMs)
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
			if m.cursor < len(m.vms)-1 {
				m.cursor++
				visible := m.visibleRows()
				if m.cursor >= m.offset+visible {
					m.offset = m.cursor - visible + 1
				}
			}
		}

	case common.VMsMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
			return m, nil
		}
		m.err = nil
		m.vms = msg.VMs
		m.cursor = 0
		m.offset = 0
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.loading {
		return "\n  " + m.spinner.View() + " Chargement des VMs..."
	}

	var s strings.Builder

	if m.err != nil {
		if strings.Contains(m.err.Error(), "not available") {
			s.WriteString("\n  " + common.StyleSubtle.Render("Les VMs ne sont pas activees sur ce serveur.") + "\n")
			s.WriteString("  " + common.StyleSubtle.Render("Activez-les dans Settings > VM Manager.") + "\n")
			return s.String()
		}
		s.WriteString("\n  " + common.StyleError.Render("⚠ "+m.err.Error()) + "\n")
	}

	running := 0
	for _, v := range m.vms {
		if v.State == "running" || v.State == "RUNNING" {
			running++
		}
	}
	title := common.StyleTitle.Render(fmt.Sprintf("  VMs (%d)", len(m.vms)))
	status := common.StyleSubtle.Render(fmt.Sprintf("  %d running", running))
	s.WriteString("\n" + title + status + "\n\n")

	// Header
	colName := 30
	colState := 15
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(common.ColorPrimary)
	header := fmt.Sprintf("  %-*s %-*s", colName, "NAME", colState, "STATE")
	s.WriteString(headerStyle.Render(header) + "\n")

	sep := "  " + strings.Repeat("─", colName+colState+1)
	s.WriteString(common.StyleSubtle.Render(sep) + "\n")

	// Rows
	visible := m.visibleRows()
	end := m.offset + visible
	if end > len(m.vms) {
		end = len(m.vms)
	}

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#CC6A1E"))

	for idx := m.offset; idx < end; idx++ {
		v := m.vms[idx]
		icon := stateIcon(v.State)
		row := fmt.Sprintf("  %-*s %-*s", colName, v.Name, colState, icon)
		if idx == m.cursor {
			s.WriteString(selectedStyle.Render(row) + "\n")
		} else {
			s.WriteString(row + "\n")
		}
	}

	if len(m.vms) == 0 && m.err == nil {
		s.WriteString("  Aucune VM configuree\n")
	}

	s.WriteString("\n" + common.StyleSubtle.Render("  ↑/↓: naviguer  │  r: rafraîchir") + "\n")
	return s.String()
}

func (m Model) visibleRows() int {
	v := m.height - 8
	if v < 5 {
		v = 5
	}
	return v
}

func stateIcon(state string) string {
	lower := strings.ToLower(state)
	switch lower {
	case "running":
		return "● running"
	case "shutoff", "shut off":
		return "○ shutoff"
	case "paused":
		return "◑ paused"
	default:
		return state
	}
}

func (m Model) fetchVMs() tea.Msg {
	vms, err := m.client.GetVMs(context.Background())
	return common.VMsMsg{VMs: vms, Err: err}
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m Model) Refresh() tea.Cmd {
	return m.fetchVMs
}
