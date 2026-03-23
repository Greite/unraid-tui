package common

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Greite/unraid-tui/internal/model"
)

// Page identifiers.
type Page int

const (
	PageDashboard Page = iota
	PageDocker
	PageVMs
	PageNotifications
	PageShares
	PageSyslog
	PageCount
)

func (p Page) Key() string {
	switch p {
	case PageDashboard:
		return "page_dashboard"
	case PageDocker:
		return "page_docker"
	case PageVMs:
		return "page_vms"
	case PageNotifications:
		return "page_notifications"
	case PageShares:
		return "page_shares"
	case PageSyslog:
		return "page_syslog"
	default:
		return "Unknown"
	}
}

// Messages sent between components.
type SystemInfoMsg struct {
	Info *model.SystemInfo
	Err  error
}

type SystemMetricsMsg struct {
	Metrics *model.SystemMetrics
	Err     error
}

type ContainersMsg struct {
	Containers []model.Container
	Err        error
}

type SharesMsg struct {
	Shares []model.Share
	Err    error
}

type ArrayInfoMsg struct {
	Info *model.ArrayInfo
	Err  error
}

type DisksMsg struct {
	Disks []model.Disk
	Err   error
}

type NetworkMsg struct {
	Network []model.NetworkAccess
	Err     error
}

type VMsMsg struct {
	VMs []model.VM
	Err error
}

type NotificationsListMsg struct {
	Notifications []model.Notification
	Err           error
}

type SharesListMsg struct {
	Shares []model.Share
	Err    error
}

// NotifRefreshRequestMsg asks the app to refresh the notification badge.
type NotifRefreshRequestMsg struct{}

type NotificationsOverviewMsg struct {
	Overview *model.NotificationOverview
	Err      error
}

type TickMsg time.Time

// Refresh interval for metrics polling.
const RefreshInterval = 3 * time.Second

// CPU temperature threshold for alerts (°C).
const CPUTempAlertThreshold = 90.0

// Bell emits a terminal bell sound.
func Bell() tea.Cmd {
	return tea.Printf("\a")
}

// Color palette.
var (
	ColorPrimary   = lipgloss.Color("#FF8C2F")
	ColorSecondary = lipgloss.Color("#6C6C6C")
	ColorSuccess   = lipgloss.Color("#73D216")
	ColorDanger    = lipgloss.Color("#FF5555")
	ColorWarning   = lipgloss.Color("#F4BF75")
	ColorMuted     = lipgloss.Color("#626262")
	ColorText      = lipgloss.Color("#FAFAFA")
	ColorBorder    = lipgloss.Color("#383838")
)

// Shared styles.
var (
	StyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	StyleSubtle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	StyleError = lipgloss.NewStyle().
			Foreground(ColorDanger).
			Bold(true)

	StyleSuccess = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	StylePanel = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(1, 2)
)

// ProgressBar renders a simple text progress bar.
func ProgressBar(percent float64, width int) string {
	if width < 4 {
		width = 4
	}
	filled := int(percent / 100 * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	bar := ""
	for i := range width {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return bar
}

// FormatBytes formats bytes to human-readable form.
func FormatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return formatUint(b) + " B"
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	suffixes := []string{"KB", "MB", "GB", "TB"}
	val := float64(b) / float64(div)
	return formatFloat(val, 1) + " " + suffixes[exp]
}

func formatFloat(f float64, decimals int) string {
	if decimals == 0 {
		return intToStr(int(f + 0.5))
	}
	whole := int(f)
	frac := int((f - float64(whole)) * 10)
	if frac < 0 {
		frac = -frac
	}
	return intToStr(whole) + "." + intToStr(frac)
}

func formatUint(n uint64) string {
	return intToStr(int(n))
}

func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + intToStr(-n)
	}
	digits := make([]byte, 0, 20)
	for n > 0 {
		digits = append(digits, byte('0'+n%10))
		n /= 10
	}
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	return string(digits)
}
