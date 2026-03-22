# Dashboard

The dashboard is the application's home page. It displays real-time system information from the Unraid server.

## Panels

### System

Displays general system information.

```
╭── System ───────────────────────────────────────────╮
│  Hostname: tower  │  OS: Unraid 6.12.6              │
│  Kernel: 6.1.64   │  Platform: linux                │
│  Board: ASRock X570 Taichi                          │
╰─────────────────────────────────────────────────────╯
```

- **Hostname**: Server name
- **OS**: Distribution and version
- **Kernel**: Linux kernel version
- **Platform**: Architecture
- **Board**: Motherboard (manufacturer + model)

### CPU

Displays processor information and real-time usage, including temperature, power draw, and per-core utilization.

```
╭── CPU ──────────────────────────────────────────────╮
│  AMD Ryzen 9 5900X                                  │
│  Cores: 16  @ 3.7 GHz                               │
│  ████████████░░░░░░░░ 62.5%   Temp: 58°C  85W       │
│                                                     │
│  Core  0 ████████████████░░░░ 78%                   │
│  Core  1 ██████████░░░░░░░░░░ 51%                   │
│  Core  2 ████████████████████ 98%                   │
│  ...                                                │
╰─────────────────────────────────────────────────────╯
```

- **Brand**: Processor model
- **Cores**: Number of physical cores
- **Speed**: Frequency in GHz
- **Progress bar**: Real-time CPU usage (%)
- **Temperature**: Current CPU temperature in degrees Celsius
- **Power**: Current power draw in watts
- **Per-core bars**: Individual utilization for each core

### Memory

Displays server memory usage.

```
╭── Memory ───────────────────╮
│  Used: 32.0 GB / 64.0 GB    │
│  ██████████░░░░░░░░░░ 50.0% │
╰─────────────────────────────╯
```

- **Used / Total**: Memory used out of total memory (human-readable format)
- **Progress bar**: Usage percentage

### Network

Displays network interface statistics and throughput.

```
╭── Network ──────────────────────────────────────────╮
│  eth0: 1 Gbps   ▲ 12.5 MB/s   ▼ 45.2 MB/s           │
│  br0:  1 Gbps   ▲  1.2 MB/s   ▼  3.4 MB/s           │
╰─────────────────────────────────────────────────────╯
```

### Disks

Displays information about array disks, cache disks, and parity disks with their status, temperature, and capacity.

```
╭── Disks ────────────────────────────────────────────╮
│  Array:                                             │
│    Disk 1  4 TB  ████████░░░░ 40%  38°C  Active     │
│    Disk 2  4 TB  ██████████░░ 52%  41°C  Active     │
│  Cache:                                             │
│    Cache  1 TB   ██░░░░░░░░░░ 12%  35°C  Active     │
│  Parity:                                            │
│    Parity 4 TB   Valid        42°C  Active          │
╰─────────────────────────────────────────────────────╯
```

### Hardware

Displays detected hardware devices including GPUs, PCI devices, USB devices, and RAM modules.

```
╭── Hardware ─────────────────────────────────────────╮
│  GPU:                                               │
│    NVIDIA GeForce RTX 3080  10 GB                   │
│  PCI:                                               │
│    Intel I211 Gigabit Network Connection            │
│  USB:                                               │
│    Logitech USB Receiver                            │
│  RAM:                                               │
│    4x 16 GB DDR4-3200 Corsair                       │
╰─────────────────────────────────────────────────────╯
```

### Array Capacity

Displays overall array capacity and usage with a summary bar.

```
╭── Array Capacity ───────────────────────────────────╮
│  Total: 16 TB   Used: 7.2 TB   Free: 8.8 TB         │
│  ████████████████████░░░░░░░░░░░░░░░░░░░░ 45%       │
╰─────────────────────────────────────────────────────╯
```

## Automatic Refresh

Metrics (CPU, memory, network, disks) are refreshed automatically every **3 seconds** via a polling mechanism. Static system information (CPU model, OS, baseboard, hardware) is loaded only once on first display.

## States

| State        | Display                                          |
|--------------|--------------------------------------------------|
| Loading      | Animated spinner + "Loading dashboard..."        |
| Data OK      | All panels: System, CPU, Memory, Network, Disks, Hardware, Array Capacity |
| Error        | Red banner with error message + cached data remains visible |

## GraphQL Queries Used

- `info` — System information (CPU, memory, OS, baseboard)
- `metrics` — Real-time metrics (CPU usage, memory usage, temperature, power, per-core)
- `network` — Network interface statistics
- `disks` — Array, cache, and parity disk information
- `hardware` — GPU, PCI, USB, and RAM devices
- `array` — Array capacity and status

## Related Files

- `internal/tui/dashboard/dashboard.go` — Bubbletea model, panel rendering
- `internal/tui/dashboard/dashboard_test.go` — Unit tests
- `internal/api/queries.go` — GraphQL queries (`querySystemInfo`, `querySystemMetrics`, `queryArrayInfo`)
