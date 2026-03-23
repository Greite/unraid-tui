package api

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Greite/unraid-tui/internal/model"
)

// Generic GraphQL response wrapper.
type graphqlResponse[T any] struct {
	Data   T              `json:"data"`
	Errors []graphqlError `json:"errors,omitempty"`
}

type graphqlError struct {
	Message string `json:"message"`
}

type graphqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

// System info response types.
type systemInfoData struct {
	Info systemInfoPayload `json:"info"`
}

type systemInfoPayload struct {
	CPU cpuPayload `json:"cpu"`
	OS  osPayload  `json:"os"`
}

type systemInfoExtraData struct {
	Info systemInfoExtraPayload `json:"info"`
}

type systemInfoExtraPayload struct {
	Versions *versionsPayload `json:"versions"`
	Devices  *devicesPayload  `json:"devices"`
	Memory   *infoMemPayload  `json:"memory"`
}

type versionsPayload struct {
	Core coreVersionsPayload `json:"core"`
}

type coreVersionsPayload struct {
	Unraid string `json:"unraid"`
	API    string `json:"api"`
	Kernel string `json:"kernel"`
}

type devicesPayload struct {
	GPU []devicePayload `json:"gpu"`
	PCI []devicePayload `json:"pci"`
	USB []devicePayload `json:"usb"`
}

type devicePayload struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Vendor string `json:"vendor"`
	Model  string `json:"model"`
}

type infoMemPayload struct {
	Layout []ramLayoutPayload `json:"layout"`
}

type ramLayoutPayload struct {
	Size         uint64 `json:"size"`
	Type         string `json:"type"`
	ClockSpeed   int    `json:"clockSpeed"`
	Manufacturer string `json:"manufacturer"`
	Bank         string `json:"bank"`
}

type cpuPayload struct {
	Manufacturer string              `json:"manufacturer"`
	Brand        string              `json:"brand"`
	Cores        int                 `json:"cores"`
	Threads      int                 `json:"threads"`
	Packages     *cpuPackagesPayload `json:"packages"`
}

type cpuPackagesPayload struct {
	TotalPower float64   `json:"totalPower"`
	Temp       []float64 `json:"temp"`
}

type osPayload struct {
	Platform string `json:"platform"`
	Distro   string `json:"distro"`
	Release  string `json:"release"`
	Uptime   string `json:"uptime"`
	Hostname string `json:"hostname"`
	Kernel   string `json:"kernel"`
}

func (p *systemInfoPayload) toDomain() *model.SystemInfo {
	return &model.SystemInfo{
		CPU: model.CPUInfo{
			Manufacturer: p.CPU.Manufacturer,
			Brand:        p.CPU.Brand,
			Cores:        p.CPU.Cores,
			Threads:      p.CPU.Threads,
			Temp:         cpuTemp(p.CPU.Packages),
			Power:        cpuPower(p.CPU.Packages),
		},
		OS: model.OSInfo{
			Platform: p.OS.Platform,
			Distro:   p.OS.Distro,
			Release:  p.OS.Release,
			Uptime:   parseUptimeISO(p.OS.Uptime),
			Hostname: p.OS.Hostname,
			Kernel:   p.OS.Kernel,
		},
	}
}

func (p *systemInfoExtraPayload) applyTo(info *model.SystemInfo) {
	if p.Versions != nil {
		info.Versions = model.VersionInfo{
			Unraid: p.Versions.Core.Unraid,
			API:    p.Versions.Core.API,
			Kernel: p.Versions.Core.Kernel,
		}
	}
	if p.Devices != nil {
		info.Hardware = model.HardwareInfo{
			GPUs: devicesToDomain(p.Devices.GPU),
			PCIs: devicesToDomain(p.Devices.PCI),
			USBs: devicesToDomain(p.Devices.USB),
		}
	}
	if p.Memory != nil {
		for _, l := range p.Memory.Layout {
			info.Hardware.RAM = append(info.Hardware.RAM, model.RAMModule{
				Size: l.Size, Type: l.Type, ClockSpeed: l.ClockSpeed,
				Manufacturer: l.Manufacturer, Bank: l.Bank,
			})
		}
	}
}

func devicesToDomain(payloads []devicePayload) []model.DeviceInfo {
	devs := make([]model.DeviceInfo, len(payloads))
	for i, p := range payloads {
		devs[i] = model.DeviceInfo{ID: p.ID, Name: p.Name, Vendor: p.Vendor, Model: p.Model}
	}
	return devs
}

func cpuTemp(pkg *cpuPackagesPayload) float64 {
	if pkg != nil && len(pkg.Temp) > 0 {
		return pkg.Temp[0]
	}
	return 0
}

func cpuPower(pkg *cpuPackagesPayload) float64 {
	if pkg != nil {
		return pkg.TotalPower
	}
	return 0
}

func parseUptimeISO(s string) int64 {
	if s == "" {
		return 0
	}
	// Try ISO 8601 timestamp (boot time) → compute seconds since boot
	t, err := time.Parse(time.RFC3339Nano, s)
	if err == nil {
		uptime := int64(time.Since(t).Seconds())
		if uptime < 0 {
			uptime = 0
		}
		return uptime
	}
	// Fallback: try as numeric seconds
	if v, err := strconv.ParseInt(s, 10, 64); err == nil {
		return v
	}
	if v, err := strconv.ParseFloat(s, 64); err == nil {
		return int64(v)
	}
	return 0
}

// System metrics response types.
type systemMetricsData struct {
	Metrics metricsPayload `json:"metrics"`
}

type metricsPayload struct {
	CPU    metricsCPUPayload    `json:"cpu"`
	Memory metricsMemoryPayload `json:"memory"`
}

type metricsCPUPayload struct {
	PercentTotal float64          `json:"percentTotal"`
	Cpus         []cpuCorePayload `json:"cpus"`
}

type cpuCorePayload struct {
	PercentTotal float64 `json:"percentTotal"`
}

type metricsMemoryPayload struct {
	Used         uint64  `json:"used"`
	Available    uint64  `json:"available"`
	Total        uint64  `json:"total"`
	PercentTotal float64 `json:"percentTotal"`
}

func (p *metricsPayload) toDomain() *model.SystemMetrics {
	cores := make([]model.CoreUsage, len(p.CPU.Cpus))
	for i, c := range p.CPU.Cpus {
		cores[i] = model.CoreUsage{Percent: c.PercentTotal}
	}
	return &model.SystemMetrics{
		CPUUsage:    p.CPU.PercentTotal,
		CPUCores:    cores,
		MemoryUsed:  p.Memory.Used,
		MemoryTotal: p.Memory.Total,
		MemoryPct:   p.Memory.PercentTotal,
	}
}

// Parity history response types.
type parityHistoryData struct {
	Array parityHistoryArrayPayload `json:"array"`
}

type parityHistoryArrayPayload struct {
	ParityHistory []parityHistoryPayload `json:"parityHistory"`
}

type parityHistoryPayload struct {
	Date     string `json:"date"`
	Status   string `json:"status"`
	Duration string `json:"duration"`
	Speed    string `json:"speed"`
	Errors   int    `json:"errors"`
}

// Container stats response types.
type containerStatsData struct {
	Docker containerStatsDockerPayload `json:"docker"`
}

type containerStatsDockerPayload struct {
	Containers []containerStatsPayload `json:"containers"`
}

type containerStatsPayload struct {
	ID         string   `json:"id"`
	Names      []string `json:"names"`
	State      string   `json:"state"`
	CPUPercent float64  `json:"cpuPercent"`
	MemUsage   uint64   `json:"memUsage"`
	MemPercent float64  `json:"memPercent"`
}

// Share response types.
type sharesData struct {
	Shares []sharePayload `json:"shares"`
}

type sharePayload struct {
	Name    string `json:"name"`
	Free    uint64 `json:"free"`
	Used    uint64 `json:"used"`
	Size    uint64 `json:"size"`
	Cache   string `json:"cache"`
	Comment string `json:"comment"`
}

func sharesToDomain(payloads []sharePayload) []model.Share {
	shares := make([]model.Share, len(payloads))
	for i, p := range payloads {
		free := p.Free * 1024 // API returns KiB
		used := p.Used * 1024
		size := p.Size * 1024
		if size == 0 && (free > 0 || used > 0) {
			size = free + used
		}
		shares[i] = model.Share{
			Name:    p.Name,
			Free:    free,
			Used:    used,
			Size:    size,
			Cache:   p.Cache,
			Comment: p.Comment,
		}
	}
	return shares
}

// Array state response types.
type arrayStateData struct {
	Array arrayStatePayload `json:"array"`
}

type arrayStatePayload struct {
	State             string               `json:"state"`
	Capacity          arrayCapacityPayload `json:"capacity"`
	ParityCheckStatus parityStatusPayload  `json:"parityCheckStatus"`
}

type arrayCapacityPayload struct {
	Kilobytes arrayKBPayload `json:"kilobytes"`
}

type arrayKBPayload struct {
	Free  json.Number `json:"free"`
	Used  json.Number `json:"used"`
	Total json.Number `json:"total"`
}

type parityStatusPayload struct {
	Status   string  `json:"status"`
	Progress float64 `json:"progress"`
	Running  bool    `json:"running"`
	Date     string  `json:"date"`
	Duration int     `json:"duration"`
	Speed    string  `json:"speed"`
	Errors   int     `json:"errors"`
}

// Notifications response types.
type notificationsOverviewData struct {
	Notifications notificationsPayload `json:"notifications"`
}

type notificationsPayload struct {
	Overview notifOverviewPayload `json:"overview"`
}

type notifOverviewPayload struct {
	Unread notifCountsPayload `json:"unread"`
}

type notifCountsPayload struct {
	Info    int `json:"info"`
	Warning int `json:"warning"`
	Alert   int `json:"alert"`
	Total   int `json:"total"`
}

// Notifications list response types.
type notificationsListData struct {
	Notifications notifListPayload `json:"notifications"`
}

type notifListPayload struct {
	List []notifItemPayload `json:"list"`
}

type notifItemPayload struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Importance  string `json:"importance"`
	Timestamp   string `json:"timestamp"`
}

func notifsToDomain(payloads []notifItemPayload) []model.Notification {
	notifs := make([]model.Notification, len(payloads))
	for i, p := range payloads {
		notifs[i] = model.Notification{
			ID: p.ID, Title: p.Title, Subject: p.Subject,
			Description: p.Description, Importance: p.Importance, Timestamp: p.Timestamp,
		}
	}
	return notifs
}

// Network response types.
type networkData struct {
	Network networkPayload `json:"network"`
}

type networkPayload struct {
	AccessUrls []accessUrlPayload `json:"accessUrls"`
}

type accessUrlPayload struct {
	Name string  `json:"name"`
	Type string  `json:"type"`
	IPv4 *string `json:"ipv4"`
	IPv6 *string `json:"ipv6"`
}

func networkToDomain(payloads []accessUrlPayload) []model.NetworkAccess {
	result := make([]model.NetworkAccess, len(payloads))
	for i, p := range payloads {
		ipv4 := ""
		if p.IPv4 != nil {
			ipv4 = *p.IPv4
		}
		ipv6 := ""
		if p.IPv6 != nil {
			ipv6 = *p.IPv6
		}
		result[i] = model.NetworkAccess{
			Name: p.Name,
			Type: p.Type,
			IPv4: ipv4,
			IPv6: ipv6,
		}
	}
	return result
}

// Disk response types (from array).
type disksData struct {
	Array arrayPayload `json:"array"`
}

type arrayPayload struct {
	Disks    []arrayDiskPayload `json:"disks"`
	Caches   []arrayDiskPayload `json:"caches"`
	Parities []arrayDiskPayload `json:"parities"`
}

type arrayDiskPayload struct {
	Name   string  `json:"name"`
	Device string  `json:"device"`
	Size   uint64  `json:"size"`
	FsSize *uint64 `json:"fsSize"`
	FsFree *uint64 `json:"fsFree"`
	FsUsed *uint64 `json:"fsUsed"`
	Status string  `json:"status"`
	Type   string  `json:"type"`
	Temp   int     `json:"temp"`
}

func allDisksToDomain(a arrayPayload) []model.Disk {
	var all []arrayDiskPayload
	all = append(all, a.Disks...)
	all = append(all, a.Caches...)
	all = append(all, a.Parities...)
	return disksToDomain(all)
}

func disksToDomain(payloads []arrayDiskPayload) []model.Disk {
	disks := make([]model.Disk, len(payloads))
	for i, p := range payloads {
		var fsSize, fsFree, fsUsed uint64
		if p.FsSize != nil {
			fsSize = *p.FsSize * 1024
		}
		if p.FsFree != nil {
			fsFree = *p.FsFree * 1024
		}
		if p.FsUsed != nil {
			fsUsed = *p.FsUsed * 1024
		}
		disks[i] = model.Disk{
			Name:   p.Name,
			Device: p.Device,
			Size:   p.Size * 1024,
			Type:   p.Type,
			FsSize: fsSize,
			FsFree: fsFree,
			FsUsed: fsUsed,
			Status: p.Status,
			Temp:   p.Temp,
		}
	}
	return disks
}

// VM response types.
type vmsData struct {
	VMs vmsPayload `json:"vms"`
}

type vmsPayload struct {
	Domains []vmDomainPayload `json:"domains"`
}

type vmDomainPayload struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}

func vmsToDomain(payloads []vmDomainPayload) []model.VM {
	vms := make([]model.VM, len(payloads))
	for i, p := range payloads {
		vms[i] = model.VM{ID: p.ID, Name: p.Name, State: p.State}
	}
	return vms
}

// Docker response types.
type dockerData struct {
	Docker dockerPayload `json:"docker"`
}

type dockerPayload struct {
	Containers              []containerPayload    `json:"containers"`
	ContainerUpdateStatuses []updateStatusPayload `json:"containerUpdateStatuses"`
}

type updateStatusPayload struct {
	Name         string `json:"name"`
	UpdateStatus string `json:"updateStatus"`
}

type containerPayload struct {
	ID        string        `json:"id"`
	Names     []string      `json:"names"`
	Image     string        `json:"image"`
	State     string        `json:"state"`
	Status    string        `json:"status"`
	AutoStart bool          `json:"autoStart"`
	Ports     []portPayload `json:"ports"`
	WebUI     string        `json:"webUiUrl"`
}

type portPayload struct {
	PrivatePort int    `json:"privatePort"`
	PublicPort  int    `json:"publicPort"`
	Type        string `json:"type"`
}

func containersToDomain(payloads []containerPayload, updateStatuses []updateStatusPayload) []model.Container {
	// Build update status map by name
	updateMap := make(map[string]bool)
	for _, u := range updateStatuses {
		updateMap[u.Name] = u.UpdateStatus == "UPDATE_AVAILABLE"
	}

	containers := make([]model.Container, len(payloads))
	for i, p := range payloads {
		ports := make([]model.Port, len(p.Ports))
		for j, port := range p.Ports {
			ports[j] = model.Port{
				HostPort:      port.PublicPort,
				ContainerPort: port.PrivatePort,
				Protocol:      port.Type,
			}
		}
		sort.Slice(ports, func(a, b int) bool {
			if ports[a].ContainerPort != ports[b].ContainerPort {
				return ports[a].ContainerPort < ports[b].ContainerPort
			}
			return ports[a].HostPort < ports[b].HostPort
		})
		name := ""
		if len(p.Names) > 0 {
			name = strings.TrimPrefix(p.Names[0], "/")
		}
		state := strings.ToLower(p.State)
		containers[i] = model.Container{
			ID:              p.ID,
			Name:            name,
			Image:           p.Image,
			State:           state,
			Status:          p.Status,
			Ports:           ports,
			WebUI:           p.WebUI,
			UpdateAvailable: updateMap[name],
		}
	}
	return containers
}
