package shares

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
	shares  []model.Share
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
	return tea.Batch(m.spinner.Tick, m.fetchShares)
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
			if m.cursor < len(m.shares)-1 {
				m.cursor++
				visible := m.visibleRows()
				if m.cursor >= m.offset+visible {
					m.offset = m.cursor - visible + 1
				}
			}
		}

	case common.SharesListMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
			return m, nil
		}
		m.err = nil
		m.shares = msg.Shares
		m.cursor = 0
		m.offset = 0
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.loading {
		return "\n  " + m.spinner.View() + " Chargement des shares..."
	}

	var s strings.Builder

	if m.err != nil {
		s.WriteString("\n  " + common.StyleError.Render("⚠ "+m.err.Error()) + "\n")
	}

	title := common.StyleTitle.Render(fmt.Sprintf("  Shares (%d)", len(m.shares)))
	s.WriteString("\n" + title + "\n\n")

	if len(m.shares) == 0 && m.err == nil {
		s.WriteString("  Aucun share configure\n")
		return s.String()
	}

	// Header
	colName := 20
	colBar := 20
	if m.width > 80 {
		colBar = m.width - 65
	}
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(common.ColorPrimary)
	header := fmt.Sprintf("  %-*s %-*s %6s %10s %10s %8s",
		colName, "NAME", colBar, "USAGE", "%", "USED", "SIZE", "CACHE")
	s.WriteString(headerStyle.Render(header) + "\n")
	sep := "  " + strings.Repeat("─", colName+colBar+50)
	s.WriteString(common.StyleSubtle.Render(sep) + "\n")

	// Rows
	visible := m.visibleRows()
	end := m.offset + visible
	if end > len(m.shares) {
		end = len(m.shares)
	}

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#CC6A1E"))

	for idx := m.offset; idx < end; idx++ {
		sh := m.shares[idx]
		var pct float64
		if sh.Size > 0 {
			pct = float64(sh.Used) / float64(sh.Size) * 100
		}
		bar := common.ProgressBar(pct, colBar)
		used := common.FormatBytes(sh.Used)
		size := common.FormatBytes(sh.Size)
		cache := sh.Cache
		if cache == "" {
			cache = "-"
		}

		row := fmt.Sprintf("  %-*s %s %5.1f%% %10s %10s %8s",
			colName, truncate(sh.Name, colName), bar, pct, used, size, cache)

		if idx == m.cursor {
			s.WriteString(selectedStyle.Render(row) + "\n")
		} else {
			s.WriteString(row + "\n")
		}

		// Comment for selected
		if idx == m.cursor && sh.Comment != "" {
			s.WriteString(common.StyleSubtle.Render("     "+sh.Comment) + "\n")
		}
	}

	s.WriteString("\n" + common.StyleSubtle.Render("  ↑/↓: naviguer  │  r: rafraîchir") + "\n")
	return s.String()
}

func truncate(s string, max int) string {
	if max < 4 {
		max = 4
	}
	if len(s) <= max {
		return s
	}
	return s[:max-2] + ".."
}

func (m Model) visibleRows() int {
	v := m.height - 10
	if v < 5 {
		v = 5
	}
	return v
}

func (m Model) fetchShares() tea.Msg {
	shares, err := m.client.GetShares(context.Background())
	return common.SharesListMsg{Shares: shares, Err: err}
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m Model) Refresh() tea.Cmd {
	return m.fetchShares
}
