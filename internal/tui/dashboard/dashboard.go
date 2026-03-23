package dashboard

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Greite/unraid-tui/internal/api"
	"github.com/Greite/unraid-tui/internal/i18n"
	"github.com/Greite/unraid-tui/internal/model"
	"github.com/Greite/unraid-tui/internal/tui/common"
)

type Model struct {
	client     api.UnraidClient
	systemInfo *model.SystemInfo
	metrics    *model.SystemMetrics
	arrayInfo  *model.ArrayInfo
	disks      []model.Disk
	network    []model.NetworkAccess
	arrayErr   error
	spinner    spinner.Model
	loading    bool
	err        error
	width      int
	height     int
	scroll     int
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
		m.fetchSystemInfo,
		m.fetchMetrics,
		m.fetchArrayInfo,
		m.fetchDisks,
		m.fetchNetwork,
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
			if m.scroll > 0 {
				m.scroll--
			}
		case "down", "j":
			m.scroll++
		case "g":
			m.scroll = 0
		}

	case common.SystemInfoMsg:
		m.loading = false
		if msg.Err != nil {
			slog.Error("system info fetch failed", "error", msg.Err)
			m.err = msg.Err
			return m, nil
		}
		m.systemInfo = msg.Info
		m.err = nil

	case common.SystemMetricsMsg:
		if msg.Err != nil {
			slog.Warn("metrics fetch failed", "error", msg.Err)
			m.err = msg.Err
			return m, m.scheduleRefresh()
		}
		m.metrics = msg.Metrics
		m.err = nil
		return m, m.scheduleRefresh()

	case common.ArrayInfoMsg:
		if msg.Err == nil {
			m.arrayInfo = msg.Info
		} else {
			slog.Warn("array info fetch failed", "error", msg.Err)
			m.arrayErr = msg.Err
		}

	case common.DisksMsg:
		if msg.Err == nil {
			m.disks = msg.Disks
		} else {
			slog.Warn("disk fetch failed", "error", msg.Err)
		}

	case common.NetworkMsg:
		if msg.Err == nil {
			m.network = msg.Network
		} else {
			slog.Warn("network fetch failed", "error", msg.Err)
		}

	case common.TickMsg:
		return m, m.fetchMetrics
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.loading {
		return "\n  " + m.spinner.View() + " " + i18n.T("loading")
	}

	var errLine string
	if m.err != nil {
		errLine = "\n  " + common.StyleError.Render("⚠ "+m.err.Error()) + "\n"
	}

	// Row 1: System + CPU summary + Memory + Network (4 cards)
	sysCard := m.renderSystemCard()
	cpuCard := m.renderCPUCard()
	memCard := m.renderMemoryCard()
	netCard := m.renderNetworkCard()

	row1 := lipgloss.JoinHorizontal(lipgloss.Top, sysCard, " ", cpuCard, " ", memCard, " ", netCard)

	// Row 2: CPU cores + Disks
	cpuCoresPanel := m.renderCPUCoresPanel()
	diskPanel := m.renderDiskPanel()

	row2 := lipgloss.JoinHorizontal(lipgloss.Top, cpuCoresPanel, " ", diskPanel)

	hwPanel := m.renderHardwarePanel()
	parityPanel := m.renderParityHistoryPanel()

	result := errLine + "\n" + row1 + "\n\n" + row2
	if hwPanel != "" || parityPanel != "" {
		row3Parts := []string{}
		if hwPanel != "" {
			row3Parts = append(row3Parts, hwPanel)
		}
		if parityPanel != "" {
			row3Parts = append(row3Parts, parityPanel)
		}
		result += "\n\n" + lipgloss.JoinHorizontal(lipgloss.Top, strings.Join(row3Parts, " "))
	}
	result += "\n"

	// Apply scroll
	lines := strings.Split(result, "\n")
	visible := m.height
	if visible < 5 {
		visible = 5
	}

	// Clamp scroll
	maxScroll := len(lines) - visible
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.scroll > maxScroll {
		m.scroll = maxScroll
	}

	start := m.scroll
	end := start + visible
	if end > len(lines) {
		end = len(lines)
	}

	return strings.Join(lines[start:end], "\n")
}

// --- Row 1: Compact cards ---

func (m Model) renderSystemCard() string {
	w := m.cardWidth()
	var content string
	if m.systemInfo != nil {
		info := m.systemInfo
		if info.OS.Hostname != "" {
			content += fmt.Sprintf("  %s: %s\n", i18n.T("hostname"), info.OS.Hostname)
		}
		content += fmt.Sprintf("  OS: %s %s\n", info.OS.Distro, info.OS.Release)
		if info.OS.Kernel != "" {
			content += fmt.Sprintf("  Kernel: %s\n", truncate(info.OS.Kernel, w-10))
		}
		if info.OS.Uptime > 0 {
			content += fmt.Sprintf("  %s: %s\n", i18n.T("uptime"), formatUptime(info.OS.Uptime))
		}
		if info.Versions.Unraid != "" {
			content += fmt.Sprintf("  Unraid: %s\n", info.Versions.Unraid)
		}
	} else {
		content = "  " + i18n.T("waiting") + "\n"
	}
	return common.StylePanel.Width(w).Render(
		common.StyleTitle.Render(i18n.T("system")) + "\n" + content,
	)
}

func (m Model) renderCPUCard() string {
	w := m.cardWidth()
	var content string
	if m.systemInfo != nil {
		content += fmt.Sprintf("  %s\n", truncate(m.systemInfo.CPU.Brand, w-6))
		specs := fmt.Sprintf("  %dC/%dT", m.systemInfo.CPU.Cores, m.systemInfo.CPU.Threads)
		if m.systemInfo.CPU.Temp > 0 {
			specs += fmt.Sprintf("  %.0f°C", m.systemInfo.CPU.Temp)
		}
		if m.systemInfo.CPU.Power > 0 {
			specs += fmt.Sprintf("  %.0fW", m.systemInfo.CPU.Power)
		}
		content += specs + "\n"
	}
	if m.metrics != nil {
		content += fmt.Sprintf("  %.1f%%\n", m.metrics.CPUUsage)
		barW := w - 8
		if barW < 5 {
			barW = 5
		}
		content += "  " + common.ProgressBar(m.metrics.CPUUsage, barW) + "\n"
	} else {
		content += "  " + i18n.T("waiting") + "\n"
	}
	return common.StylePanel.Width(w).Render(
		common.StyleTitle.Render(i18n.T("cpu")) + "\n" + content,
	)
}

func (m Model) renderMemoryCard() string {
	w := m.cardWidth()
	var content string
	if m.metrics != nil {
		used := common.FormatBytes(m.metrics.MemoryUsed)
		total := common.FormatBytes(m.metrics.MemoryTotal)
		content += fmt.Sprintf("  %s / %s\n", used, total)
		content += fmt.Sprintf("  %.1f%%\n", m.metrics.MemoryPct)
		barW := w - 8
		if barW < 5 {
			barW = 5
		}
		content += "  " + common.ProgressBar(m.metrics.MemoryPct, barW) + "\n"
	} else {
		content = "  " + i18n.T("waiting") + "\n"
	}
	return common.StylePanel.Width(w).Render(
		common.StyleTitle.Render(i18n.T("memory")) + "\n" + content,
	)
}

func (m Model) renderNetworkCard() string {
	w := m.cardWidth()
	var content string
	if len(m.network) > 0 {
		for _, n := range m.network {
			ip := n.IPv4
			if ip == "" {
				ip = n.IPv6
			}
			if ip == "" {
				continue
			}
			typeLabel := n.Type
			switch typeLabel {
			case "DEFAULT":
				typeLabel = "Default"
			case "LAN":
				typeLabel = "LAN"
			case "MDNS":
				typeLabel = "mDNS"
			}
			content += fmt.Sprintf("  %-8s %s\n", typeLabel, truncate(ip, w-14))
		}
	} else {
		content = "  " + i18n.T("waiting") + "\n"
	}
	return common.StylePanel.Width(w).Render(
		common.StyleTitle.Render(i18n.T("network")) + "\n" + content,
	)
}

// --- Row 2: Half-width panels ---

func (m Model) renderCPUCoresPanel() string {
	w := m.halfWidth()
	var content string
	if m.metrics != nil && len(m.metrics.CPUCores) > 0 {
		cols := 2
		coreBarW := (w-10)/cols - 14
		if coreBarW < 3 {
			coreBarW = 3
		}
		for i := 0; i < len(m.metrics.CPUCores); i += cols {
			line := "  "
			for j := 0; j < cols && i+j < len(m.metrics.CPUCores); j++ {
				core := m.metrics.CPUCores[i+j]
				bar := common.ProgressBar(core.Percent, coreBarW)
				line += fmt.Sprintf("cpu%-2d %s %3.0f%%  ", i+j, bar, core.Percent)
			}
			content += line + "\n"
		}
	} else {
		content = "  " + i18n.T("waiting") + "\n"
	}
	title := i18n.T("cpu_cores")
	if m.metrics != nil {
		title += fmt.Sprintf(" (%d)", len(m.metrics.CPUCores))
	}
	return common.StylePanel.Width(w).Render(
		common.StyleTitle.Render(title) + "\n" + content,
	)
}

func (m Model) renderDiskPanel() string {
	if len(m.disks) == 0 {
		return ""
	}

	w := m.halfWidth()
	var content string

	// Fixed text: "  ● DAT disk1   " (18) + "  53%  1005.6 GB / 1.8 TB  25°C" (~35) + border/padding (6)
	barWidth := w - 59
	if barWidth < 3 {
		barWidth = 3
	}

	for _, d := range m.disks {
		tempStr := ""
		if d.Temp > 0 {
			tempStr = fmt.Sprintf(" %d°C", d.Temp)
		}

		statusIcon := "●"
		if d.Status != "DISK_OK" {
			statusIcon = "⚠"
		}

		typeLabel := ""
		switch d.Type {
		case "PARITY":
			typeLabel = "PAR"
		case "CACHE":
			typeLabel = "NVM"
		default:
			typeLabel = "DAT"
		}

		if d.FsSize > 0 {
			pct := float64(d.FsUsed) / float64(d.FsSize) * 100
			used := common.FormatBytes(d.FsUsed)
			total := common.FormatBytes(d.FsSize)
			bar := common.ProgressBar(pct, barWidth)
			content += fmt.Sprintf("  %s %s %-7s %s  %4.0f%%  %9s / %-9s %4s\n",
				statusIcon, typeLabel, d.Name, bar, pct, used, total, tempStr)
		} else {
			size := common.FormatBytes(d.Size)
			// No filesystem — show size in the "total" column position, aligned with rows above
			// FS rows: bar(barWidth) + "  " + pct(4%) + "  " + used(9) + " / " + total(-9) + " " + temp(4)
			content += fmt.Sprintf("  %s %s %-7s %-*s  %5s  %9s / %-9s %4s\n",
				statusIcon, typeLabel, d.Name, barWidth, "", "--", "--", size, tempStr)
		}
	}

	// Parity check status
	if m.arrayInfo != nil {
		if m.arrayInfo.ParityRunning {
			parityBar := common.ProgressBar(m.arrayInfo.ParityProgress, barWidth)
			content += fmt.Sprintf("\n  %s: %s %s %.1f%%\n", i18n.T("parity"), m.arrayInfo.ParityStatus, parityBar, m.arrayInfo.ParityProgress)
		} else if m.arrayInfo.ParityStatus != "" {
			content += fmt.Sprintf("\n  %s: %s\n", i18n.T("parity"), m.arrayInfo.ParityStatus)
		}
	}

	return common.StylePanel.Width(w).Render(
		common.StyleTitle.Render(i18n.T("disks")) + fmt.Sprintf(" (%d)", len(m.disks)) + "\n" + content,
	)
}

func (m Model) renderHardwarePanel() string {
	if m.systemInfo == nil {
		return ""
	}
	hw := m.systemInfo.Hardware
	if len(hw.RAM) == 0 && len(hw.GPUs) == 0 && len(hw.USBs) == 0 {
		return ""
	}

	w := m.halfWidth()
	var content string

	// RAM
	if len(hw.RAM) > 0 {
		var totalRAM uint64
		for _, r := range hw.RAM {
			totalRAM += r.Size
		}
		first := hw.RAM[0]
		ramType := first.Type
		if ramType == "" {
			ramType = "DDR"
		}
		speed := ""
		if first.ClockSpeed > 0 {
			speed = fmt.Sprintf(" %dMHz", first.ClockSpeed)
		}
		content += fmt.Sprintf("  RAM    %dx %s %s%s  (%s %s)\n",
			len(hw.RAM), common.FormatBytes(first.Size), ramType, speed,
			common.FormatBytes(totalRAM), i18n.T("total"))
	}

	// GPU
	for _, g := range hw.GPUs {
		name := g.Name
		if name == "" {
			name = g.Model
		}
		if name == "" {
			continue
		}
		content += fmt.Sprintf("  GPU    %s\n", name)
	}

	// USB summary
	if len(hw.USBs) > 0 {
		content += fmt.Sprintf("  USB    %d %s\n", len(hw.USBs), i18n.T("devices"))
	}

	// PCI summary
	if len(hw.PCIs) > 0 {
		content += fmt.Sprintf("  PCI    %d %s\n", len(hw.PCIs), i18n.T("devices"))
	}

	if content == "" {
		return ""
	}

	return common.StylePanel.Width(w).Render(
		common.StyleTitle.Render(i18n.T("hardware")) + "\n" + content,
	)
}

func (m Model) renderParityHistoryPanel() string {
	w := m.halfWidth()
	var content string

	if m.arrayInfo == nil || m.arrayInfo.ParityDate == "" {
		if m.arrayErr != nil {
			content = "  " + common.StyleError.Render(m.arrayErr.Error()) + "\n"
		} else {
			content = "  " + common.StyleSubtle.Render(i18n.T("no_parity_history")) + "\n"
		}
		return common.StylePanel.Width(w).Render(
			common.StyleTitle.Render(i18n.T("parity_history")) + "\n" + content,
		)
	}

	a := m.arrayInfo

	// Date
	dateStr := a.ParityDate
	if len(dateStr) >= 10 {
		dateStr = dateStr[:10] // Keep YYYY-MM-DD
	}
	content += fmt.Sprintf("  %s:    %s\n", i18n.T("parity_date"), dateStr)

	// Status
	statusStyle := lipgloss.NewStyle().Foreground(common.ColorSuccess)
	if a.ParityStatus != "COMPLETED" {
		statusStyle = lipgloss.NewStyle().Foreground(common.ColorWarning)
	}
	content += fmt.Sprintf("  %s:  %s\n", i18n.T("parity_status"), statusStyle.Render(a.ParityStatus))

	// Duration
	if a.ParityDuration > 0 {
		hours := a.ParityDuration / 3600
		mins := (a.ParityDuration % 3600) / 60
		content += fmt.Sprintf("  %s: %dh %dm\n", i18n.T("parity_duration"), hours, mins)
	}

	// Speed
	if a.ParitySpeed != "" && a.ParitySpeed != "0" {
		content += fmt.Sprintf("  %s:   %s MB/s\n", i18n.T("parity_speed"), a.ParitySpeed)
	}

	// Errors
	errStr := fmt.Sprintf("%d", a.ParityErrors)
	errStyle := lipgloss.NewStyle()
	if a.ParityErrors > 0 {
		errStyle = lipgloss.NewStyle().Foreground(common.ColorDanger).Bold(true)
	}
	content += fmt.Sprintf("  %s:  %s\n", i18n.T("parity_errors"), errStyle.Render(errStr))

	return common.StylePanel.Width(w).Render(
		common.StyleTitle.Render(i18n.T("parity_history")) + "\n" + content,
	)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-2] + ".."
}

func (m Model) cardWidth() int {
	w := (m.width - 6) / 4
	if w < 20 {
		w = 20
	}
	return w
}

func (m Model) halfWidth() int {
	w := (m.width - 4) / 2
	if w < 30 {
		w = 30
	}
	return w
}

func (m Model) fetchSystemInfo() tea.Msg {
	info, err := m.client.GetSystemInfo(context.Background())
	if err == nil && info != nil {
		// Best-effort fetch of extra info (versions, hardware)
		m.client.GetSystemInfoExtra(context.Background(), info)
	}
	return common.SystemInfoMsg{Info: info, Err: err}
}

func (m Model) fetchArrayInfo() tea.Msg {
	info, err := m.client.GetArrayInfo(context.Background())
	return common.ArrayInfoMsg{Info: info, Err: err}
}

func (m Model) fetchNetwork() tea.Msg {
	network, err := m.client.GetNetwork(context.Background())
	return common.NetworkMsg{Network: network, Err: err}
}

func (m Model) fetchDisks() tea.Msg {
	disks, err := m.client.GetDisks(context.Background())
	return common.DisksMsg{Disks: disks, Err: err}
}

func (m Model) fetchMetrics() tea.Msg {
	metrics, err := m.client.GetSystemMetrics(context.Background())
	return common.SystemMetricsMsg{Metrics: metrics, Err: err}
}

func (m Model) scheduleRefresh() tea.Cmd {
	return tea.Tick(common.RefreshInterval, func(t time.Time) tea.Msg {
		return common.TickMsg(t)
	})
}

// SetSize updates the dimensions for layout.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func formatUptime(seconds int64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	mins := (seconds % 3600) / 60
	if days > 0 {
		return fmt.Sprintf("%dj %dh %dm", days, hours, mins)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, mins)
	}
	return fmt.Sprintf("%dm", mins)
}
