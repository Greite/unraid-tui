# Unraid TUI

A terminal user interface (TUI) to monitor and manage [Unraid](https://unraid.net/) servers from the command line, via the GraphQL API.

## Disclaimer

This project is not affiliated with, endorsed by, or associated with [Lime Technology, Inc.](https://lime-technology.com/) or the [Unraid](https://unraid.net/) project. Unraid is a registered trademark of Lime Technology, Inc. This is an independent, open-source tool that interacts with the Unraid API.

```
╭─────────────────── UNRAID TUI ──────────────────────────╮
│  [Dashboard]  Docker  VMs  Shares  Notifications        │
├─────────────────────────────────────────────────────────┤
│  ╭── CPU ──────────────╮  ╭── Memory ────────────╮      │
│  │ AMD Ryzen 9 5900X   │  │ Used: 24.3 / 64.0 GB │      │
│  │ Cores: 12 / 24T     │  │ ████████░░░░ 38.0%   │      │
│  │ Usage: ███████░ 72% │  │                      │      │
│  │ Temp: 62°C  85W     │  │                      │      │
│  ╰─────────────────────╯  ╰──────────────────────╯      │
│  ╭── System ──────────────────────────────────────╮     │
│  │ Hostname: tower  │ OS: Unraid 7.2.0            │     │
│  │ Kernel: 6.12.8   │ Board: ASRock X570 Taichi   │     │
│  ╰────────────────────────────────────────────────╯     │
├─────────────────────────────────────────────────────────┤
│  Tab switch page │ Ctrl+S server │ Ctrl+L lang │ q quit │
╰─────────────────────────────────────────────────────────╯
```

## Features

- **Dashboard** -- System info, CPU (per-core usage, temperature, power draw), Memory, Network interfaces, Disks (array + cache + parity), Hardware (GPU/PCI/USB/RAM), Array capacity, Parity status, Unraid version. Auto-refreshes every 3 seconds.
- **Docker** -- Sortable table of all containers. Start/stop/pause/unpause, update single or update all, log viewer with follow mode, SSH console access, open WebUI. Update-available indicator per container.
- **VMs** -- List with current state. Start/stop/pause/resume/reboot/force-stop actions.
- **Shares** -- List of all shares with usage progress bars.
- **Notifications** -- List with importance levels. Archive a single notification or all at once. Unread badge displayed in the header.
- **Multi-server** -- Manage named servers. Switch with `Ctrl+S` picker. Add, delete, or set a default server. `--server` CLI flag for scripting.
- **Internationalization** -- English (default) and French. Switch with `Ctrl+L` picker. `--lang` CLI flag and automatic `LANG` environment variable detection.
- **Onboarding** -- Guided setup wizard on first launch: server name, URL connectivity test, API key entry.
- **Secure config** -- Configuration stored in `~/.unraid-tui/config.yaml`. API keys are saved in the system keychain (macOS Keychain, GNOME Keyring, Windows Credential Manager).

## Prerequisites

- **Unraid 7.2+** (built-in GraphQL API) or Unraid 6.x with the [Unraid Connect](https://docs.unraid.net/API/) plugin
- An **Unraid API key** (the onboarding wizard will guide you through creation)

## Installation

### Homebrew (macOS / Linux)

```bash
brew install Greite/tap/unraid-tui
```

Upgrade:

```bash
brew upgrade unraid-tui
```

### Pre-built binaries

Download the binary matching your OS and architecture from the [Releases](https://github.com/Greite/unraid-tui/releases) page:

| OS      | Architecture  | File                                   |
|---------|---------------|----------------------------------------|
| macOS   | Apple Silicon | `unraid-tui_x.x.x_darwin_arm64.tar.gz`|
| macOS   | Intel         | `unraid-tui_x.x.x_darwin_amd64.tar.gz`|
| Linux   | x86_64        | `unraid-tui_x.x.x_linux_amd64.tar.gz` |
| Linux   | ARM64         | `unraid-tui_x.x.x_linux_arm64.tar.gz` |
| Windows | x86_64        | `unraid-tui_x.x.x_windows_amd64.zip`  |

```bash
# Example: macOS Apple Silicon
tar xzf unraid-tui_*_darwin_arm64.tar.gz
sudo mv unraid-tui /usr/local/bin/
```

### Go install

Requires Go 1.26+:

```bash
go install github.com/Greite/unraid-tui@latest
```

### From source

```bash
git clone https://github.com/Greite/unraid-tui.git
cd unraid-tui
make install
```

This compiles the binary and copies it to `/usr/local/bin/`.

To uninstall:

```bash
make uninstall
```

## Quick start

```bash
unraid-tui
```

On first launch, an interactive wizard guides you through:

1. Choosing a name for your server
2. Entering the server URL and testing the connection
3. Creating and entering your API key
4. Saving the configuration

The configuration is stored in `~/.unraid-tui/config.yaml`, with API keys secured in the system keychain.

### CLI flags

```bash
unraid-tui --server myserver   # Connect to a specific named server
unraid-tui --lang fr           # Force a specific language
```

### Manual configuration

If you prefer to configure manually, create `~/.unraid-tui/config.yaml`:

```yaml
server_url: "http://192.168.1.100:3001"
api_key: "your-api-key"
```

Or via environment variables:

```bash
export UNRAID_SERVER_URL="http://192.168.1.100:3001"
export UNRAID_API_KEY="your-api-key"
```

### Obtaining an API key

1. Open the Unraid web interface
2. Go to **Settings > Management Access > Developer Options**
3. Open Apollo GraphQL Studio
4. Run:

```graphql
mutation {
  apiKey {
    create(input: {
      name: "unraid-tui"
      roles: [ADMIN]
    }) { key }
  }
}
```

5. Copy the returned key

## Usage

### Keyboard shortcuts

| Key          | Action                              |
|--------------|-------------------------------------|
| `Tab`        | Next page                           |
| `Shift+Tab`  | Previous page                       |
| `1`-`5`     | Jump to page by number              |
| `Ctrl+S`    | Open server picker (multi-server)   |
| `Ctrl+L`    | Open language picker                |
| `↑` / `↓`   | Navigate within tables/lists        |
| `Enter`     | Confirm / open action menu          |
| `r`         | Refresh current page                |
| `q`         | Quit                                |
| `Ctrl+C`    | Quit                                |

### Pages

#### Dashboard

Displays real-time system metrics:
- **CPU** -- Model, core count, frequency, per-core usage, temperature, and power draw
- **Memory** -- Used / total with progress bar
- **Network** -- Interface list with transfer rates
- **Disks** -- Array, cache, and parity drives with status
- **Hardware** -- GPU, PCI devices, USB devices, RAM modules
- **Array** -- Total capacity and usage
- **Parity** -- Status and last check result
- **System** -- Hostname, Unraid version, kernel, motherboard

Metrics auto-refresh every 3 seconds.

#### Docker

Interactive table of all containers:

| Column | Description                                         |
|--------|-----------------------------------------------------|
| NAME   | Container name                                      |
| IMAGE  | Docker image                                        |
| STATE  | Current state (running / exited / paused indicator) |
| STATUS | Detail ("Up 14 days", "Exited (0) 2 days ago")     |
| UPDATE | Update-available indicator                          |
| PORTS  | Host:container port mapping                         |

The table is sortable by clicking column headers. Actions available on a selected container: start, stop, pause, unpause, update, view logs (with follow mode), open SSH console, and open WebUI. An "Update All" action updates every container with a pending update.

#### VMs

Lists all virtual machines with their current state. Available actions: start, stop, pause, resume, reboot, and force-stop.

#### Shares

Displays all configured shares with usage progress bars showing consumed vs. available space.

#### Notifications

Lists system notifications with importance levels (normal, warning, alert). You can archive a single notification or archive all at once. An unread badge is shown in the header bar.

## Development

### Commands

```bash
make build         # Compile the binary
make test          # Run tests
make test-verbose  # Run tests with verbose output
make test-cover    # Run tests with HTML coverage report
make lint          # Run go vet
make run           # Build and run
make clean         # Clean build artifacts
```

### Architecture

```
cmd/                        -> Cobra CLI entry point
internal/api/               -> GraphQL client (interface + HTTP)
internal/config/            -> Viper-based configuration
internal/i18n/              -> Internationalization (en, fr)
internal/model/             -> Domain types
internal/tui/               -> Bubbletea app (router, header, footer)
internal/tui/common/        -> Shared styles, messages, helpers
internal/tui/dashboard/     -> Dashboard page
internal/tui/docker/        -> Docker page
internal/tui/vms/           -> VMs page
internal/tui/shares/        -> Shares page
internal/tui/notifications/ -> Notifications page
internal/tui/onboarding/    -> Guided setup wizard
```

### Tests

```bash
make test
```

### Release

Releases are managed by [GoReleaser](https://goreleaser.com/). On each Git tag, GoReleaser:

1. Compiles binaries for macOS, Linux, and Windows (amd64 + arm64)
2. Creates `.tar.gz` / `.zip` archives
3. Publishes the GitHub release
4. Updates the Homebrew formula in [Greite/homebrew-tap](https://github.com/Greite/homebrew-tap)

```bash
# Test the release locally (dry run, no publish)
make release-dry

# Publish a release
git tag v0.1.0
git push origin v0.1.0
goreleaser release --clean
```

### Tech stack

| Component  | Library                                                                   |
|------------|---------------------------------------------------------------------------|
| TUI        | [Bubbletea v2](https://github.com/charmbracelet/bubbletea)               |
| Styling    | [Lipgloss v2](https://github.com/charmbracelet/lipgloss)                 |
| Components | [Bubbles v2](https://github.com/charmbracelet/bubbles) (table, spinner, textinput) |
| CLI        | [Cobra](https://github.com/spf13/cobra)                                  |
| Config     | [Viper](https://github.com/spf13/viper)                                  |
| Keychain   | [go-keyring](https://github.com/zalando/go-keyring)                      |
| API        | [Unraid GraphQL API](https://docs.unraid.net/API/) via `net/http`        |

## Documentation

Detailed documentation for each feature is in the [`docs/`](docs/) directory:

- [Configuration](docs/configuration.md)
- [Onboarding](docs/onboarding.md)
- [Dashboard](docs/dashboard.md)
- [Docker](docs/docker.md)
- [Navigation](docs/navigation.md)
- [API Client](docs/api-client.md)

## License

MIT
