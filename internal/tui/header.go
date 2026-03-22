package tui

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/Greite/unraid-tui/internal/model"
	"github.com/Greite/unraid-tui/internal/tui/common"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(common.ColorText).
			Background(common.ColorPrimary).
			Padding(0, 2)

	tabKeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#AAAAAA"))

	tabKeyActiveStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#000000")).
				Background(common.ColorPrimary)

	tabLabelStyle = lipgloss.NewStyle().
			Foreground(common.ColorText).
			Background(lipgloss.Color("#2A2A2A"))

	tabLabelActiveStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(common.ColorText).
				Background(lipgloss.Color("#CC6A1E"))
)

// TabZones stores the X positions of each tab for mouse click detection.
var TabZones []TabZone

type TabZone struct {
	Page  common.Page
	Start int
	End   int
}

func RenderHeader(activePage common.Page, width int, notifs *model.NotificationOverview) string {
	title := titleStyle.Render(" UNRAID CLI ")

	tabs := ""
	TabZones = nil
	cursor := lipgloss.Width(title) + 2

	type tabDef struct {
		key  string
		page common.Page
	}
	tabDefs := []tabDef{
		{"F1", common.PageDashboard},
		{"F2", common.PageDocker},
		{"F3", common.PageVMs},
		{"F4", common.PageNotifications},
		{"F5", common.PageShares},
	}

	for i, td := range tabDefs {
		var key, label string
		if td.page == activePage {
			key = tabKeyActiveStyle.Render(td.key)
			label = tabLabelActiveStyle.Render(" " + td.page.String() + " ")
		} else {
			key = tabKeyStyle.Render(td.key)
			label = tabLabelStyle.Render(" " + td.page.String() + " ")
		}
		tab := key + label
		tabWidth := lipgloss.Width(tab)

		TabZones = append(TabZones, TabZone{
			Page:  td.page,
			Start: cursor,
			End:   cursor + tabWidth,
		})

		tabs += tab
		cursor += tabWidth

		// Space between tabs
		if i < len(tabDefs)-1 {
			tabs += " "
			cursor++
		}
	}

	// Notification badges
	notifBadge := ""
	if notifs != nil && notifs.Total > 0 {
		if notifs.Alert > 0 {
			notifBadge += lipgloss.NewStyle().Bold(true).Foreground(common.ColorDanger).Render(fmt.Sprintf(" ✗%d", notifs.Alert))
		}
		if notifs.Warning > 0 {
			notifBadge += lipgloss.NewStyle().Foreground(common.ColorWarning).Render(fmt.Sprintf(" ⚠%d", notifs.Warning))
		}
		if notifs.Info > 0 {
			notifBadge += lipgloss.NewStyle().Foreground(common.ColorMuted).Render(fmt.Sprintf(" ●%d", notifs.Info))
		}
	}

	// Right-aligned page indicator
	indicator := lipgloss.NewStyle().
		Foreground(common.ColorMuted).
		Render(fmt.Sprintf("%d/%d", int(activePage)+1, int(common.PageCount)))

	header := lipgloss.JoinHorizontal(lipgloss.Center, title, " ", tabs, notifBadge)

	// Fill remaining width
	headerWidth := lipgloss.Width(header)
	indicatorWidth := lipgloss.Width(indicator)
	gap := width - headerWidth - indicatorWidth - 1
	if gap < 0 {
		gap = 0
	}
	padding := fmt.Sprintf("%*s", gap, "")

	return lipgloss.NewStyle().
		Width(width).
		Render(header + padding + indicator)
}
