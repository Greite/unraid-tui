package tui

import (
	"context"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Greite/unraid-tui/internal/api"
	"github.com/Greite/unraid-tui/internal/model"
	"github.com/Greite/unraid-tui/internal/tui/common"
	"github.com/Greite/unraid-tui/internal/tui/dashboard"
	"github.com/Greite/unraid-tui/internal/tui/docker"
	"github.com/Greite/unraid-tui/internal/tui/notifications"
	"github.com/Greite/unraid-tui/internal/tui/shares"
	"github.com/Greite/unraid-tui/internal/tui/vms"
)

type notifTickMsg struct{}

type Model struct {
	activePage    common.Page
	dashboard     dashboard.Model
	docker        docker.Model
	vms           vms.Model
	notifications notifications.Model
	shares        shares.Model
	client        api.UnraidClient
	notifOverview *model.NotificationOverview
	width         int
	height        int
}

func NewModel(client api.UnraidClient) Model {
	return Model{
		activePage:    common.PageDashboard,
		dashboard:     dashboard.New(client),
		docker:        docker.New(client),
		vms:           vms.New(client),
		notifications: notifications.New(client),
		shares:        shares.New(client),
		client:        client,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.dashboard.Init(),
		m.fetchNotifOverview,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		contentHeight := msg.Height - 4
		m.dashboard.SetSize(msg.Width, contentHeight)
		m.docker.SetSize(msg.Width, contentHeight)
		m.vms.SetSize(msg.Width, contentHeight)
		m.notifications.SetSize(msg.Width, contentHeight)
		m.shares.SetSize(msg.Width, contentHeight)

	case tea.KeyPressMsg:
		switch {
		case msg.Code == 'c' && msg.Mod&tea.ModCtrl != 0:
			return m, tea.Quit
		case msg.Code == 'q' && !(m.activePage == common.PageDocker && m.docker.InSubView()):
			return m, tea.Quit
		case msg.Code == tea.KeyTab && msg.Mod&tea.ModShift != 0:
			return m, m.switchPage((m.activePage - 1 + common.PageCount) % common.PageCount)
		case msg.Code == tea.KeyTab:
			return m, m.switchPage((m.activePage + 1) % common.PageCount)
		case msg.Code == tea.KeyF1:
			return m, m.switchPage(common.PageDashboard)
		case msg.Code == tea.KeyF2:
			return m, m.switchPage(common.PageDocker)
		case msg.Code == tea.KeyF3:
			return m, m.switchPage(common.PageVMs)
		case msg.Code == tea.KeyF4:
			return m, m.switchPage(common.PageNotifications)
		case msg.Code == tea.KeyF5:
			return m, m.switchPage(common.PageShares)
		case msg.Code == 'r':
			return m, m.refreshActivePage()
		}

	case tea.MouseClickMsg:
		mouse := msg.Mouse()
		if mouse.Y == 0 {
			for _, zone := range TabZones {
				if mouse.X >= zone.Start && mouse.X < zone.End {
					return m, m.switchPage(zone.Page)
				}
			}
		}

	case common.NotificationsOverviewMsg:
		if msg.Err == nil {
			m.notifOverview = msg.Overview
		}
		return m, m.scheduleNotifRefresh()

	case notifTickMsg:
		return m, m.fetchNotifOverview
	}

	var cmd tea.Cmd
	switch m.activePage {
	case common.PageDashboard:
		m.dashboard, cmd = m.dashboard.Update(msg)
	case common.PageDocker:
		m.docker, cmd = m.docker.Update(msg)
	case common.PageVMs:
		m.vms, cmd = m.vms.Update(msg)
	case common.PageNotifications:
		m.notifications, cmd = m.notifications.Update(msg)
	case common.PageShares:
		m.shares, cmd = m.shares.Update(msg)
	}
	return m, cmd
}

func (m *Model) switchPage(page common.Page) tea.Cmd {
	prev := m.activePage
	m.activePage = page
	if page == prev {
		return nil
	}
	switch page {
	case common.PageDashboard:
		return m.dashboard.Init()
	case common.PageDocker:
		return m.docker.Init()
	case common.PageVMs:
		return m.vms.Init()
	case common.PageNotifications:
		return m.notifications.Init()
	case common.PageShares:
		return m.shares.Init()
	}
	return nil
}

func (m Model) refreshActivePage() tea.Cmd {
	switch m.activePage {
	case common.PageDocker:
		return m.docker.Refresh()
	case common.PageVMs:
		return m.vms.Refresh()
	case common.PageNotifications:
		return m.notifications.Refresh()
	case common.PageShares:
		return m.shares.Refresh()
	}
	return nil
}

func (m Model) View() tea.View {
	header := RenderHeader(m.activePage, m.width, m.notifOverview)

	var content string
	switch m.activePage {
	case common.PageDashboard:
		content = m.dashboard.View()
	case common.PageDocker:
		content = m.docker.View()
	case common.PageVMs:
		content = m.vms.View()
	case common.PageNotifications:
		content = m.notifications.View()
	case common.PageShares:
		content = m.shares.View()
	}

	footer := RenderFooter(m.width)

	contentHeight := m.height - 3
	if contentHeight < 5 {
		contentHeight = 5
	}
	contentBox := lipgloss.NewStyle().
		Height(contentHeight).
		Render(content)

	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Left, header, contentBox, footer))
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

func (m Model) fetchNotifOverview() tea.Msg {
	overview, err := m.client.GetNotificationsOverview(context.Background())
	return common.NotificationsOverviewMsg{Overview: overview, Err: err}
}

func (m Model) scheduleNotifRefresh() tea.Cmd {
	return tea.Tick(30*time.Second, func(_ time.Time) tea.Msg {
		return notifTickMsg{}
	})
}

// ActivePage returns the current page (for testing).
func (m Model) ActivePage() common.Page {
	return m.activePage
}
