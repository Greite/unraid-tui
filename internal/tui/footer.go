package tui

import (
	"charm.land/lipgloss/v2"
	"github.com/Greite/unraid-tui/internal/i18n"
	"github.com/Greite/unraid-tui/internal/tui/common"
)

var (
	footerKeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(common.ColorText)

	footerDescStyle = lipgloss.NewStyle().
			Foreground(common.ColorMuted)

	footerStyle = lipgloss.NewStyle().
			Foreground(common.ColorMuted).
			Padding(0, 1)
)

func RenderFooter(width int) string {
	keys := []struct{ key, desc string }{
		{"F1-F7", i18n.T("pages")},
		{"tab", i18n.T("next")},
		{"^S", i18n.T("servers")},
		{"^L", i18n.T("lang")},
		{"q", i18n.T("quit")},
	}

	content := ""
	for i, k := range keys {
		if i > 0 {
			content += footerDescStyle.Render("  │  ")
		}
		content += footerKeyStyle.Render(k.key) + " " + footerDescStyle.Render(k.desc)
	}

	return footerStyle.Width(width).Render(content)
}
