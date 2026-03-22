package model

type SystemInfo struct {
	CPU    CPUInfo
	Memory MemoryInfo
	OS     OSInfo
}

type CPUInfo struct {
	Manufacturer string
	Brand        string
	Cores        int
	Threads      int
	Temp         float64 // °C (first package)
	Power        float64 // Watts (total)
}

type MemoryInfo struct {
	Total uint64
	Used  uint64
	Free  uint64
}

type OSInfo struct {
	Platform string
	Distro   string
	Release  string
	Uptime   int64 // seconds since boot
	Hostname string
	Kernel   string
}

type SystemMetrics struct {
	CPUUsage    float64 // 0-100
	CPUCores    []CoreUsage
	MemoryUsed  uint64
	MemoryTotal uint64
	MemoryPct   float64
}

type CoreUsage struct {
	Percent float64
}

type ArrayInfo struct {
	State          string
	Free           uint64 // bytes
	Used           uint64
	Total          uint64
	ParityStatus   string
	ParityProgress float64
	ParityRunning  bool
}

type NotificationOverview struct {
	Info    int
	Warning int
	Alert   int
	Total   int
}

type Share struct {
	Name    string
	Free    uint64
	Used    uint64
	Size    uint64
	Cache   string
	Comment string
}

type NetworkAccess struct {
	Name string
	Type string
	IPv4 string
	IPv6 string
}

type Disk struct {
	Name   string
	Device string
	Size   uint64
	Type   string // DATA, PARITY, CACHE, FLASH
	FsSize uint64
	FsFree uint64
	FsUsed uint64
	Status string
	Temp   int
}

type Notification struct {
	ID          string
	Title       string
	Subject     string
	Description string
	Importance  string // INFO, WARNING, ALERT
	Timestamp   string
}

type VM struct {
	ID    string
	Name  string
	State string
}

type Container struct {
	ID     string
	Name   string
	Image  string
	State  string // running, stopped, paused
	Status string // "Up 2 hours", "Exited (0) 3 days ago"
	Ports           []Port
	WebUI           string // URL if available
	UpdateAvailable bool
}

type Port struct {
	HostPort      int
	ContainerPort int
	Protocol      string // tcp, udp
}
