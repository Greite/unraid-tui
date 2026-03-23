package notifications

import (
	"context"
	"fmt"
	"strings"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Greite/unraid-tui/internal/api"
	"github.com/Greite/unraid-tui/internal/i18n"
	"github.com/Greite/unraid-tui/internal/model"
	"github.com/Greite/unraid-tui/internal/tui/common"
)

type Model struct {
	client        api.UnraidClient
	notifications []model.Notification
	spinner       spinner.Model
	loading       bool
	err           error
	cursor        int
	offset        int
	width         int
	height        int
}

func New(client api.UnraidClient) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(common.ColorPrimary)
	return Model{client: client, spinner: s, loading: true}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchNotifications)
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
			if m.cursor < len(m.notifications)-1 {
				m.cursor++
				visible := m.visibleRows()
				if m.cursor >= m.offset+visible {
					m.offset = m.cursor - visible + 1
				}
			}
		case "a":
			return m, m.archiveSelected()
		case "A":
			return m, m.archiveAll()
		}

	case notifActionMsg:
		if msg.Err == nil {
			return m, tea.Batch(m.fetchNotifications, func() tea.Msg {
				return common.NotifRefreshRequestMsg{}
			})
		}
		m.err = msg.Err

	case common.NotificationsListMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
			return m, nil
		}
		m.err = nil
		m.notifications = msg.Notifications
		m.cursor = 0
		m.offset = 0
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.loading {
		return "\n  " + m.spinner.View() + " " + i18n.T("loading_notifs")
	}

	var s strings.Builder

	if m.err != nil {
		s.WriteString("\n  " + common.StyleError.Render("⚠ "+m.err.Error()) + "\n")
	}

	title := common.StyleTitle.Render(fmt.Sprintf("  %s (%d)", i18n.T("notifications"), len(m.notifications)))
	s.WriteString("\n" + title + "\n\n")

	if len(m.notifications) == 0 && m.err == nil {
		s.WriteString("  " + i18n.T("no_notifs") + "\n")
		return s.String()
	}

	visible := m.visibleRows()
	end := m.offset + visible
	if end > len(m.notifications) {
		end = len(m.notifications)
	}

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#CC6A1E"))

	for idx := m.offset; idx < end; idx++ {
		n := m.notifications[idx]
		icon := importanceIcon(n.Importance)
		subject := n.Subject
		if subject == "" {
			subject = n.Title
		}
		ts := ""
		if len(n.Timestamp) >= 10 {
			ts = n.Timestamp[:10]
		}
		row := fmt.Sprintf("  %s  %-50s  %s", icon, truncate(subject, 50), ts)
		if idx == m.cursor {
			s.WriteString(selectedStyle.Render(row) + "\n")
		} else {
			s.WriteString(row + "\n")
		}

		// Show description for selected
		if idx == m.cursor && n.Description != "" {
			desc := common.StyleSubtle.Render("     " + truncate(n.Description, m.width-10))
			s.WriteString(desc + "\n")
		}
	}

	s.WriteString("\n" + common.StyleSubtle.Render("  ↑/↓: "+i18n.T("navigate")+"  │  a: "+i18n.T("archive")+"  │  A: "+i18n.T("archive_all")+"  │  r: "+i18n.T("refresh")) + "\n")
	return s.String()
}

func importanceIcon(importance string) string {
	switch strings.ToUpper(importance) {
	case "ALERT":
		return lipgloss.NewStyle().Foreground(common.ColorDanger).Render("✗")
	case "WARNING":
		return lipgloss.NewStyle().Foreground(common.ColorWarning).Render("⚠")
	default:
		return lipgloss.NewStyle().Foreground(common.ColorMuted).Render("●")
	}
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
	v := m.height - 8
	if v < 5 {
		v = 5
	}
	return v
}

type notifActionMsg struct{ Err error }

func (m Model) archiveSelected() tea.Cmd {
	if m.cursor >= len(m.notifications) {
		return nil
	}
	n := m.notifications[m.cursor]
	id, client := n.ID, m.client
	return func() tea.Msg {
		return notifActionMsg{client.ArchiveNotification(context.Background(), id)}
	}
}

func (m Model) archiveAll() tea.Cmd {
	client := m.client
	return func() tea.Msg {
		return notifActionMsg{client.ArchiveAllNotifications(context.Background())}
	}
}

func (m Model) fetchNotifications() tea.Msg {
	notifs, err := m.client.GetNotifications(context.Background())
	return common.NotificationsListMsg{Notifications: notifs, Err: err}
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m Model) Refresh() tea.Cmd {
	return m.fetchNotifications
}
