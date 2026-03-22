package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Greite/unraid-tui/internal/api"
	"github.com/Greite/unraid-tui/internal/config"
	"github.com/Greite/unraid-tui/internal/i18n"
	"github.com/Greite/unraid-tui/internal/model"
	"github.com/Greite/unraid-tui/internal/tui/common"
	"github.com/Greite/unraid-tui/internal/tui/dashboard"
	"github.com/Greite/unraid-tui/internal/tui/docker"
	"github.com/Greite/unraid-tui/internal/tui/notifications"
	"github.com/Greite/unraid-tui/internal/tui/onboarding"
	"github.com/Greite/unraid-tui/internal/tui/shares"
	"github.com/Greite/unraid-tui/internal/tui/vms"
)

type notifTickMsg struct{}

// serverSwitchedMsg is sent after switching to a new server.
type serverSwitchedMsg struct {
	client api.UnraidClient
	err    error
}

type Model struct {
	activePage      common.Page
	dashboard       dashboard.Model
	docker          docker.Model
	vms             vms.Model
	notifications   notifications.Model
	shares          shares.Model
	client          api.UnraidClient
	notifOverview   *model.NotificationOverview
	width           int
	height          int
	// Server selector
	showServerPicker bool
	serverList       []config.ServerEntry
	serverCursor     int
	// Inline onboarding
	onboarding     *onboarding.Model
	showOnboarding bool
	// Language picker
	showLangPicker bool
	langCursor     int
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
	// Language picker mode
	if m.showLangPicker {
		return m.updateLangPicker(msg)
	}

	// Inline onboarding mode
	if m.showOnboarding && m.onboarding != nil {
		updated, cmd := m.onboarding.Update(msg)
		ob := updated.(onboarding.Model)
		m.onboarding = &ob
		if ob.Completed() {
			m.showOnboarding = false
			m.onboarding = nil
			// Reload server list and switch to new server
			servers := config.ListServers()
			if len(servers) > 0 {
				last := servers[len(servers)-1]
				return m, m.switchToServer(last.Name)
			}
		}
		if ob.Quitting() {
			m.showOnboarding = false
			m.onboarding = nil
		}
		return m, cmd
	}

	// Server picker mode
	if m.showServerPicker {
		return m.updateServerPicker(msg)
	}

	switch msg := msg.(type) {
	case serverSwitchedMsg:
		if msg.err != nil {
			return m, nil
		}
		m.client = msg.client
		m.dashboard = dashboard.New(msg.client)
		m.docker = docker.New(msg.client)
		m.vms = vms.New(msg.client)
		m.notifications = notifications.New(msg.client)
		m.shares = shares.New(msg.client)
		contentHeight := m.height - 4
		m.dashboard.SetSize(m.width, contentHeight)
		m.docker.SetSize(m.width, contentHeight)
		m.vms.SetSize(m.width, contentHeight)
		m.notifications.SetSize(m.width, contentHeight)
		m.shares.SetSize(m.width, contentHeight)
		m.activePage = common.PageDashboard
		return m, tea.Batch(m.dashboard.Init(), m.fetchNotifOverview)

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
		case msg.Code == 'l' && msg.Mod&tea.ModCtrl != 0:
			m.showLangPicker = true
			m.langCursor = 0
			if i18n.Lang() == "fr" {
				m.langCursor = 1
			}
			return m, nil
		case msg.Code == 's' && msg.Mod&tea.ModCtrl != 0:
			m.serverList = config.ListServers()
			m.serverCursor = 0
			m.showServerPicker = true
			return m, nil
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

	case common.NotifRefreshRequestMsg:
		return m, m.fetchNotifOverview

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

	if m.showLangPicker {
		content = m.renderLangPicker()
	} else if m.showOnboarding && m.onboarding != nil {
		obView := m.onboarding.View()
		content = obView.Content
	} else if m.showServerPicker {
		content = m.renderServerPicker()
	} else {
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

func (m Model) updateServerPicker(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc", "ctrl+s":
			m.showServerPicker = false
			return m, nil
		case "up", "k":
			if m.serverCursor > 0 {
				m.serverCursor--
			}
		case "down", "j":
			if m.serverCursor < len(m.serverList) {
				m.serverCursor++
			}
		case "enter":
			if m.serverCursor < len(m.serverList) {
				s := m.serverList[m.serverCursor]
				m.showServerPicker = false
				return m, m.switchToServer(s.Name)
			}
			if m.serverCursor == len(m.serverList) {
				m.showServerPicker = false
				ob := onboarding.New()
				m.onboarding = &ob
				m.showOnboarding = true
				return m, m.onboarding.Init()
			}
		case "d":
			if m.serverCursor < len(m.serverList) {
				s := m.serverList[m.serverCursor]
				config.SetDefault(s.Name)
			}
		case "x":
			if m.serverCursor < len(m.serverList) && len(m.serverList) > 1 {
				s := m.serverList[m.serverCursor]
				config.RemoveServer(s.Name)
				m.serverList = config.ListServers()
				if m.serverCursor >= len(m.serverList) {
					m.serverCursor = len(m.serverList) - 1
				}
			}
		}
	}
	return m, nil
}

func (m Model) switchToServer(name string) tea.Cmd {
	return func() tea.Msg {
		cfg, err := config.LoadServer(name)
		if err != nil {
			return serverSwitchedMsg{err: err}
		}
		client := api.NewClient(cfg.ServerURL, cfg.APIKey)
		return serverSwitchedMsg{client: client}
	}
}

func (m Model) renderServerPicker() string {
	var s strings.Builder
	s.WriteString("\n")

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(common.ColorPrimary)
	s.WriteString(titleStyle.Render("  " + i18n.T("server_picker_title")) + "\n\n")

	def := config.DefaultServer()
	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#CC6A1E"))

	for i, srv := range m.serverList {
		marker := "  "
		if srv.Name == def {
			marker = "* "
		}
		row := fmt.Sprintf("  %s%-15s  %s", marker, srv.Name, srv.ServerURL)
		if i == m.serverCursor {
			s.WriteString(selectedStyle.Render(row) + "\n")
		} else {
			s.WriteString(row + "\n")
		}
	}

	// Add new option
	addRow := "  " + i18n.T("add_server")
	if m.serverCursor == len(m.serverList) {
		s.WriteString(selectedStyle.Render(addRow) + "\n")
	} else {
		s.WriteString(common.StyleSubtle.Render(addRow) + "\n")
	}

	s.WriteString("\n" + common.StyleSubtle.Render("  enter: "+i18n.T("connect")+"  │  d: "+i18n.T("default")+"  │  x: "+i18n.T("delete")+"  │  esc: "+i18n.T("close")) + "\n")
	return s.String()
}

var langOptions = []struct {
	code string
	name string
}{
	{"en", "English"},
	{"fr", "Francais"},
}

func (m Model) updateLangPicker(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc", "ctrl+l":
			m.showLangPicker = false
		case "up", "k":
			if m.langCursor > 0 {
				m.langCursor--
			}
		case "down", "j":
			if m.langCursor < len(langOptions)-1 {
				m.langCursor++
			}
		case "enter":
			i18n.SetLang(langOptions[m.langCursor].code)
			m.showLangPicker = false
		}
	}
	return m, nil
}

func (m Model) renderLangPicker() string {
	var s strings.Builder
	s.WriteString("\n")

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(common.ColorPrimary)
	s.WriteString(titleStyle.Render("  Language / Langue") + "\n\n")

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#CC6A1E"))

	for idx, opt := range langOptions {
		marker := "  "
		if opt.code == i18n.Lang() {
			marker = "* "
		}
		row := fmt.Sprintf("  %s%s  (%s)", marker, opt.name, opt.code)
		if idx == m.langCursor {
			s.WriteString(selectedStyle.Render(row) + "\n")
		} else {
			s.WriteString(row + "\n")
		}
	}

	s.WriteString("\n" + common.StyleSubtle.Render("  enter: "+i18n.T("select")+"  │  esc: "+i18n.T("close")+"  │  * = "+i18n.T("default")) + "\n")
	return s.String()
}

// ActivePage returns the current page (for testing).
func (m Model) ActivePage() common.Page {
	return m.activePage
}
