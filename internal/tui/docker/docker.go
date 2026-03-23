package docker

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os/exec"
	"sort"
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

type sortColumn int

const (
	sortName sortColumn = iota
	sortImage
	sortState
	sortStatus
	sortAuto
	sortPorts
	sortColumnCount
)

func (c sortColumn) label() string {
	switch c {
	case sortName:
		return "NAME"
	case sortImage:
		return "IMAGE"
	case sortState:
		return "STATE"
	case sortStatus:
		return "STATUS"
	case sortAuto:
		return "AUTO"
	case sortPorts:
		return "PORTS"
	default:
		return ""
	}
}

type viewMode int

const (
	viewList viewMode = iota
	viewLogs
)

type colZone struct {
	col   sortColumn
	start int
	end   int
}

var (
	columnZones   []colZone
	columnHeaderY int
)

// Messages
type LogsMsg struct {
	Name string
	Logs string
	Err  error
}

type ConsoleOutputMsg struct {
	Output string
	Err    error
}

type ContainerActionMsg struct {
	Action string
	Name   string
	Err    error
}

type logsTickMsg struct{}

type Model struct {
	client     api.UnraidClient
	containers []model.Container
	sorted     []model.Container
	spinner    spinner.Model
	loading    bool
	err        error
	cursor     int
	offset     int
	width      int
	height     int
	sortCol    sortColumn
	sortAsc    bool
	mode       viewMode
	logs       string
	logsName   string
	logsOffset int
	logsFollow bool
	statusMsg  string
}

func New(client api.UnraidClient) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(common.ColorPrimary)

	return Model{
		client:  client,
		spinner: s,
		loading: true,
		sortCol: sortName,
		sortAsc: true,
		mode:    viewList,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.fetchContainers,
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.mode == viewLogs {
		return m.updateLogs(msg)
	}

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
			if m.cursor < len(m.sorted)-1 {
				m.cursor++
				visible := m.visibleRows()
				if m.cursor >= m.offset+visible {
					m.offset = m.cursor - visible + 1
				}
			}
		case "n":
			m.toggleSort(sortName)
		case "i":
			m.toggleSort(sortImage)
		case "s":
			m.toggleSort(sortState)
		case "t":
			m.toggleSort(sortStatus)
		case "p":
			m.toggleSort(sortPorts)
		case "o":
			m.toggleSort(sortAuto)
		case "l":
			return m, m.fetchLogs()
		case "w":
			m.openWebUI()
		case "c":
			return m, m.execConsole()
		case "S":
			return m, m.toggleStartStop()
		case "P":
			return m, m.togglePause()
		case "u":
			return m, m.updateContainer()
		case "a":
			return m, m.toggleAutostart()
		case "U":
			return m, m.updateAllContainers()
		}

	case tea.MouseClickMsg:
		mouse := msg.Mouse()
		if mouse.Y == columnHeaderY {
			for _, z := range columnZones {
				if mouse.X >= z.start && mouse.X < z.end {
					m.toggleSort(z.col)
					return m, nil
				}
			}
		}

	case LogsMsg:
		if msg.Err != nil {
			m.statusMsg = fmt.Sprintf(i18n.T("logs_error"), msg.Err.Error())
			return m, nil
		}
		m.mode = viewLogs
		m.logsName = msg.Name
		m.logs = msg.Logs
		m.logsFollow = true
		// Start at bottom
		lines := strings.Split(m.logs, "\n")
		maxOffset := len(lines) - m.logsVisible()
		if maxOffset < 0 {
			maxOffset = 0
		}
		m.logsOffset = maxOffset
		return m, m.scheduleLogsRefresh()

	case ContainerActionMsg:
		if msg.Err != nil {
			m.statusMsg = fmt.Sprintf(i18n.T("action_error"), msg.Action, msg.Name, msg.Err)
		} else {
			m.statusMsg = fmt.Sprintf(i18n.T("action_ok"), msg.Action, msg.Name)
		}
		return m, m.fetchContainers

	case ConsoleOutputMsg:
		if msg.Err != nil {
			m.statusMsg = i18n.T("console_error")
		} else {
			m.statusMsg = i18n.T("console_done")
		}

	case common.ContainersMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
			return m, nil
		}
		m.err = nil
		m.containers = msg.Containers
		m.applySort()
		m.cursor = 0
		m.offset = 0
		m.statusMsg = ""
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m Model) updateLogs(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyPressMsg:
		maxOffset := m.logsMaxOffset()
		switch msg.String() {
		case "esc", "q":
			m.mode = viewList
			m.logs = ""
			m.logsFollow = false
			return m, nil
		case "up", "k":
			m.logsFollow = false
			if m.logsOffset > 0 {
				m.logsOffset--
			}
		case "down", "j":
			if m.logsOffset < maxOffset {
				m.logsOffset++
			}
			// Re-enable follow if we scroll back to the bottom
			if m.logsOffset >= maxOffset {
				m.logsFollow = true
			}
		case "g":
			m.logsFollow = false
			m.logsOffset = 0
		case "G":
			m.logsFollow = true
			m.logsOffset = maxOffset
		case "f":
			m.logsFollow = !m.logsFollow
			if m.logsFollow {
				m.logsOffset = maxOffset
			}
		}
		return m, nil

	case LogsMsg:
		if msg.Err != nil {
			return m, m.scheduleLogsRefresh()
		}
		m.logs = msg.Logs
		if m.logsFollow {
			m.logsOffset = m.logsMaxOffset()
		}
		return m, m.scheduleLogsRefresh()

	case logsTickMsg:
		if m.mode == viewLogs {
			return m, m.fetchLogsCmd()
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.mode == viewLogs {
		return m.viewLogs()
	}
	return m.viewList()
}

func (m Model) viewLogs() string {
	var s strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(common.ColorPrimary)
	followIndicator := ""
	if m.logsFollow {
		followIndicator = common.StyleSuccess.Render("  ● " + i18n.T("follow_on"))
	} else {
		followIndicator = common.StyleSubtle.Render("  ○ " + i18n.T("follow_off"))
	}
	s.WriteString("\n" + titleStyle.Render(fmt.Sprintf("  Logs — %s", m.logsName)) + followIndicator + "\n")
	s.WriteString(common.StyleSubtle.Render("  "+strings.Repeat("─", 40)) + "\n")

	lines := strings.Split(m.logs, "\n")
	visible := m.logsVisible()
	start := m.logsOffset
	if start > len(lines) {
		start = len(lines)
	}
	end := start + visible
	if end > len(lines) {
		end = len(lines)
	}

	lineNoStyle := lipgloss.NewStyle().Foreground(common.ColorMuted)

	for i, line := range lines[start:end] {
		lineNo := lineNoStyle.Render(fmt.Sprintf(" %4d ", start+i+1))
		s.WriteString(lineNo + line + "\n")
	}

	total := len(lines)
	pos := ""
	if total > 0 {
		pct := 0
		if maxOff := total - visible; maxOff > 0 {
			pct = m.logsOffset * 100 / maxOff
		}
		pos = fmt.Sprintf("  %s %d-%d / %d (%d%%)", i18n.T("line"), start+1, end, total, pct)
	}

	s.WriteString("\n" + common.StyleSubtle.Render(pos+"  │  ↑/↓: "+i18n.T("scroll")+"  │  f: "+i18n.T("follow")+"  │  g/G: "+i18n.T("start_end")+"  │  esc: "+i18n.T("back")) + "\n")
	return s.String()
}

func (m Model) viewList() string {
	if m.loading {
		return "\n  " + m.spinner.View() + " " + i18n.T("loading_docker")
	}

	var s strings.Builder

	if m.err != nil {
		if strings.Contains(m.err.Error(), "not available") {
			s.WriteString("\n  " + common.StyleSubtle.Render(i18n.T("docker_disabled")) + "\n")
			s.WriteString("  " + common.StyleSubtle.Render(i18n.T("docker_enable")) + "\n")
			return s.String()
		}
		s.WriteString("\n  " + common.StyleError.Render("⚠ "+m.err.Error()) + "\n")
	}

	if m.statusMsg != "" {
		s.WriteString("\n  " + common.StyleSubtle.Render("→ "+m.statusMsg) + "\n")
	}

	running := 0
	updates := 0
	for _, c := range m.containers {
		if c.State == "running" {
			running++
		}
		if c.UpdateAvailable {
			updates++
		}
	}
	title := common.StyleTitle.Render(fmt.Sprintf("  %s (%d)", i18n.T("containers"), len(m.containers)))
	status := common.StyleSubtle.Render(fmt.Sprintf("  %d %s", running, i18n.T("running")))
	if updates > 0 {
		status += "  " + lipgloss.NewStyle().Foreground(common.ColorWarning).Render(fmt.Sprintf("⬆ %d %s", updates, i18n.T("update")))
	}
	s.WriteString("\n" + title + status + "\n\n")

	colName, colImage, colState, colStatus, colAuto, colPorts := m.colWidths()

	colDefs := []struct {
		col   sortColumn
		width int
	}{
		{sortName, colName},
		{sortImage, colImage},
		{sortState, colState},
		{sortStatus, colStatus},
		{sortAuto, colAuto},
		{sortPorts, colPorts},
	}

	lineCount := strings.Count(s.String(), "\n") + 1
	header := "  "
	columnZones = nil
	xCursor := 2
	for _, c := range colDefs {
		label := c.col.label()
		if c.col == m.sortCol {
			arrow := "▲"
			if !m.sortAsc {
				arrow = "▼"
			}
			label = label + " " + arrow
		}
		colStr := fmt.Sprintf("%-*s ", c.width, label)
		columnZones = append(columnZones, colZone{
			col:   c.col,
			start: xCursor,
			end:   xCursor + c.width,
		})
		header += colStr
		xCursor += c.width + 1
	}
	columnHeaderY = lineCount

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(common.ColorPrimary)
	s.WriteString(headerStyle.Render(header) + "\n")

	totalWidth := colName + colImage + colState + colStatus + colAuto + colPorts + 5
	sep := "  " + strings.Repeat("─", totalWidth)
	s.WriteString(common.StyleSubtle.Render(sep) + "\n")

	visible := m.visibleRows()
	end := m.offset + visible
	if end > len(m.sorted) {
		end = len(m.sorted)
	}

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#CC6A1E"))
	normalStyle := lipgloss.NewStyle()

	for idx := m.offset; idx < end; idx++ {
		c := m.sorted[idx]
		autoIcon := "○"
		if c.AutoStart {
			autoIcon = "●"
		}
		row := fmt.Sprintf("  %-*s %-*s %-*s %-*s %-*s %-*s",
			colName, truncate(c.Name, colName),
			colImage, truncate(c.Image, colImage),
			colState, stateIcon(c.State, c.UpdateAvailable),
			colStatus, truncate(c.Status, colStatus),
			colAuto, autoIcon,
			colPorts, truncate(formatPorts(c.Ports), colPorts),
		)
		if idx == m.cursor {
			s.WriteString(selectedStyle.Render(row) + "\n")
		} else {
			s.WriteString(normalStyle.Render(row) + "\n")
		}
	}

	// Actions for selected container
	if m.cursor < len(m.sorted) {
		c := m.sorted[m.cursor]
		var actions []string
		if c.State == "running" {
			actions = append(actions, "S: "+i18n.T("stop"), "P: "+i18n.T("pause"), "l: "+i18n.T("logs"), "c: "+i18n.T("console"))
		} else if c.State == "paused" {
			actions = append(actions, "P: "+i18n.T("unpause"))
		} else {
			actions = append(actions, "S: "+i18n.T("start"))
		}
		if c.UpdateAvailable {
			actions = append(actions, "u: "+i18n.T("update")+" ⬆")
		}
		if c.WebUI != "" || hasHTTPPort(c.Ports) {
			actions = append(actions, "w: "+i18n.T("webui"))
		}
		if c.AutoStart {
			actions = append(actions, "a: "+i18n.T("autostart")+" ●")
		} else {
			actions = append(actions, "a: "+i18n.T("autostart")+" ○")
		}
		if len(actions) > 0 {
			s.WriteString("\n  " + common.StyleSubtle.Render(strings.Join(actions, "  │  ")))
		}
	}

	s.WriteString("\n" + common.StyleSubtle.Render("  n/i/s/t/o/p: "+i18n.T("sort")+"  │  ↑/↓: "+i18n.T("navigate")+"  │  r: "+i18n.T("refresh")+"  │  U: "+i18n.T("update_all")) + "\n")

	return s.String()
}

// --- Actions ---

func (m Model) logsVisible() int {
	v := m.height - 6
	if v < 5 {
		v = 5
	}
	return v
}

func (m Model) logsMaxOffset() int {
	lines := strings.Split(m.logs, "\n")
	max := len(lines) - m.logsVisible()
	if max < 0 {
		max = 0
	}
	return max
}

func (m Model) scheduleLogsRefresh() tea.Cmd {
	return tea.Tick(3*time.Second, func(_ time.Time) tea.Msg {
		return logsTickMsg{}
	})
}

func (m Model) fetchLogsCmd() tea.Cmd {
	name := m.logsName
	host := extractHost(m.client.ServerURL())
	return m.doFetchLogs(name, host)
}

func (m Model) fetchLogs() tea.Cmd {
	if m.cursor >= len(m.sorted) {
		return nil
	}
	c := m.sorted[m.cursor]
	if c.State != "running" {
		return nil
	}
	name := c.Name
	host := extractHost(m.client.ServerURL())
	return m.doFetchLogs(name, host)
}

func (m Model) doFetchLogs(name, host string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("ssh",
			"-o", "StrictHostKeyChecking=no",
			"-o", "ConnectTimeout=5",
			"root@"+host,
			"docker", "logs", "--tail", "500", name,
		)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			return LogsMsg{Name: name, Err: fmt.Errorf("ssh: %s", firstNonEmpty(stderr.String(), err.Error()))}
		}
		output := stdout.String()
		if output == "" {
			output = stderr.String()
		}
		return LogsMsg{Name: name, Logs: output}
	}
}

func (m *Model) toggleStartStop() tea.Cmd {
	if m.cursor >= len(m.sorted) {
		return nil
	}
	c := m.sorted[m.cursor]
	id := c.ID
	name := c.Name
	client := m.client

	if c.State == "running" {
		m.statusMsg = "⏳ " + i18n.T("stop") + " " + name + "..."
		return func() tea.Msg {
			err := client.StopContainer(context.Background(), id)
			return ContainerActionMsg{Action: "Stop", Name: name, Err: err}
		}
	}
	m.statusMsg = "⏳ " + i18n.T("start") + " " + name + "..."
	return func() tea.Msg {
		err := client.StartContainer(context.Background(), id)
		return ContainerActionMsg{Action: "Start", Name: name, Err: err}
	}
}

func (m Model) togglePause() tea.Cmd {
	if m.cursor >= len(m.sorted) {
		return nil
	}
	c := m.sorted[m.cursor]
	id := c.ID
	name := c.Name
	client := m.client

	if c.State == "paused" {
		return func() tea.Msg {
			err := client.UnpauseContainer(context.Background(), id)
			return ContainerActionMsg{Action: "Unpause", Name: name, Err: err}
		}
	}
	if c.State == "running" {
		return func() tea.Msg {
			err := client.PauseContainer(context.Background(), id)
			return ContainerActionMsg{Action: "Pause", Name: name, Err: err}
		}
	}
	return nil
}

func (m *Model) updateContainer() tea.Cmd {
	if m.cursor >= len(m.sorted) {
		return nil
	}
	c := m.sorted[m.cursor]
	if !c.UpdateAvailable {
		m.statusMsg = c.Name + ": " + i18n.T("up_to_date")
		return nil
	}
	id := c.ID
	name := c.Name
	m.statusMsg = "⏳ " + i18n.T("updating") + " " + name + "..."
	client := m.client
	return func() tea.Msg {
		err := client.UpdateContainer(context.Background(), id)
		return ContainerActionMsg{Action: "Update", Name: name, Err: err}
	}
}

func (m *Model) updateAllContainers() tea.Cmd {
	m.statusMsg = "⏳ " + i18n.T("updating_all") + "..."
	client := m.client
	return func() tea.Msg {
		err := client.UpdateAllContainers(context.Background())
		return ContainerActionMsg{Action: "Update all", Name: "", Err: err}
	}
}

func (m *Model) toggleAutostart() tea.Cmd {
	if m.cursor >= len(m.sorted) {
		return nil
	}
	c := m.sorted[m.cursor]
	id := c.ID
	name := c.Name
	newState := !c.AutoStart
	client := m.client
	allContainers := m.containers

	label := i18n.T("autostart_on")
	if !newState {
		label = i18n.T("autostart_off")
	}
	m.statusMsg = "⏳ " + label + " " + name + "..."
	return func() tea.Msg {
		err := client.SetAutostart(context.Background(), allContainers, id, newState)
		return ContainerActionMsg{Action: "Autostart", Name: name, Err: err}
	}
}

func (m Model) execConsole() tea.Cmd {
	if m.cursor >= len(m.sorted) {
		return nil
	}
	c := m.sorted[m.cursor]
	if c.State != "running" {
		return nil
	}
	host := extractHost(m.client.ServerURL())

	cmd := exec.Command("ssh",
		"-o", "StrictHostKeyChecking=no",
		"-t",
		"root@"+host,
		"docker", "exec", "-it", c.Name, "/bin/sh",
	)

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return ConsoleOutputMsg{Err: err}
	})
}

func (m *Model) openWebUI() {
	if m.cursor >= len(m.sorted) {
		return
	}
	c := m.sorted[m.cursor]
	webURL := c.WebUI
	if webURL == "" {
		webURL = guessWebUI(m.client.ServerURL(), c.Ports)
	}
	if webURL == "" {
		m.statusMsg = fmt.Sprintf(i18n.T("no_webui"), c.Name)
		return
	}
	openBrowser(webURL)
	m.statusMsg = fmt.Sprintf(i18n.T("webui_opened"), c.Name)
}

// --- Sort ---

func (m *Model) toggleSort(col sortColumn) {
	if m.sortCol == col {
		m.sortAsc = !m.sortAsc
	} else {
		m.sortCol = col
		m.sortAsc = true
	}
	m.applySort()
	m.cursor = 0
	m.offset = 0
}

func (m *Model) applySort() {
	m.sorted = make([]model.Container, len(m.containers))
	copy(m.sorted, m.containers)

	sort.SliceStable(m.sorted, func(i, j int) bool {
		var less bool
		a, b := m.sorted[i], m.sorted[j]
		switch m.sortCol {
		case sortName:
			less = strings.ToLower(a.Name) < strings.ToLower(b.Name)
		case sortImage:
			less = strings.ToLower(a.Image) < strings.ToLower(b.Image)
		case sortState:
			less = a.State < b.State
		case sortStatus:
			less = strings.ToLower(a.Status) < strings.ToLower(b.Status)
		case sortAuto:
			less = !a.AutoStart && b.AutoStart
		case sortPorts:
			less = formatPorts(a.Ports) < formatPorts(b.Ports)
		}
		if !m.sortAsc {
			less = !less
		}
		return less
	})
}

// --- Helpers ---

func (m Model) visibleRows() int {
	v := m.height - 10
	if v < 5 {
		v = 5
	}
	return v
}

func (m Model) colWidths() (int, int, int, int, int, int) {
	available := m.width - 8
	if available < 60 {
		available = 60
	}
	return available * 18 / 100,
		available * 22 / 100,
		available * 12 / 100,
		available * 22 / 100,
		6,
		available * 15 / 100
}

func stateIcon(state string, updateAvailable bool) string {
	icon := ""
	switch state {
	case "running":
		icon = "● " + i18n.T("running")
	case "exited":
		icon = "○ " + i18n.T("exited")
	case "paused":
		icon = "◑ " + i18n.T("paused")
	default:
		icon = state
	}
	if updateAvailable {
		icon += " ⬆"
	}
	return icon
}

func formatPorts(ports []model.Port) string {
	if len(ports) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(ports))
	for _, p := range ports {
		if p.HostPort > 0 {
			parts = append(parts, fmt.Sprintf("%d:%d", p.HostPort, p.ContainerPort))
		} else {
			parts = append(parts, fmt.Sprintf("%d", p.ContainerPort))
		}
	}
	return strings.Join(parts, ", ")
}

func hasHTTPPort(ports []model.Port) bool {
	for _, p := range ports {
		if p.HostPort > 0 {
			return true
		}
	}
	return false
}

func guessWebUI(serverURL string, ports []model.Port) string {
	parsed, err := url.Parse(serverURL)
	if err != nil {
		return ""
	}
	host := parsed.Hostname()
	for _, p := range ports {
		if p.HostPort > 0 {
			scheme := "http"
			if p.ContainerPort == 443 || p.HostPort == 443 {
				scheme = "https"
			}
			return fmt.Sprintf("%s://%s:%d", scheme, host, p.HostPort)
		}
	}
	return ""
}

func extractHost(serverURL string) string {
	parsed, err := url.Parse(serverURL)
	if err != nil {
		return serverURL
	}
	return parsed.Hostname()
}

func firstNonEmpty(strs ...string) string {
	for _, s := range strs {
		s = strings.TrimSpace(s)
		if s != "" {
			return s
		}
	}
	return "unknown error"
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

func openBrowser(rawURL string) {
	cmd := exec.Command("open", rawURL)
	cmd.Start()
}

func (m Model) fetchContainers() tea.Msg {
	containers, err := m.client.GetContainers(context.Background())
	return common.ContainersMsg{Containers: containers, Err: err}
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// InSubView returns true when Docker is in logs or another sub-view where 'q' should not quit the app.
func (m Model) InSubView() bool {
	return m.mode != viewList
}

// UpdateCount returns the number of containers with updates available.
func (m Model) UpdateCount() int {
	count := 0
	for _, c := range m.containers {
		if c.UpdateAvailable {
			count++
		}
	}
	return count
}

func (m Model) Refresh() tea.Cmd {
	return m.fetchContainers
}
