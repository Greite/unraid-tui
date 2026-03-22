# Navigation

The application uses a multi-page TUI interface with keyboard navigation.

## Available Pages

| Page           | Shortcut | Description                            |
|----------------|----------|----------------------------------------|
| Dashboard      | `F1`     | CPU, memory, network, disks, hardware monitoring |
| Docker         | `F2`     | Docker container list and management   |
| VMs            | `F3`     | Virtual machine list and management    |
| Notifications  | `F4`     | System notifications                   |
| Shares         | `F5`     | Network shares                         |

## Global Shortcuts

These shortcuts work from any page:

| Key          | Action                                        |
|--------------|-----------------------------------------------|
| `Tab`        | Next page (cyclic)                            |
| `Shift+Tab`  | Previous page (cyclic)                        |
| `F1`         | Go to Dashboard                               |
| `F2`         | Go to Docker                                  |
| `F3`         | Go to VMs                                     |
| `F4`         | Go to Notifications                           |
| `F5`         | Go to Shares                                  |
| `Ctrl+S`     | Switch between configured servers             |
| `Ctrl+L`     | Switch language (EN/FR)                       |
| `q`          | Quit the application                          |
| `Ctrl+C`     | Quit the application                          |

## Page-Specific Shortcuts

### Docker

| Key     | Action                     |
|---------|----------------------------|
| `↑`/`↓` | Navigate the table         |
| `Enter` | Open action menu           |
| `s`     | Cycle sort mode            |
| `r`     | Refresh containers         |

### VMs

| Key     | Action                     |
|---------|----------------------------|
| `↑`/`↓` | Navigate the table         |
| `Enter` | Open action menu           |
| `r`     | Refresh VMs                |

### Notifications

| Key     | Action                     |
|---------|----------------------------|
| `↑`/`↓` | Navigate the list          |
| `d`     | Dismiss selected           |
| `D`     | Dismiss all                |

## Interface

```
┌──────────────────────────────────────────────────────────────────────────┐
│  UNRAID CLI    Dashboard   Docker   VMs   Notifications   Shares        │  ← Header with tabs
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  Active page content                                                     │  ← Content area
│                                                                          │
├──────────────────────────────────────────────────────────────────────────┤
│  F1-F5 pages  │  Tab switch  │  Ctrl+S server  │  Ctrl+L lang  │  q quit│  ← Footer with help
└──────────────────────────────────────────────────────────────────────────┘
```

### Header

- Title "UNRAID CLI" in purple
- Page tabs with the active tab highlighted
- Connected server name displayed

### Footer

- Reminder of main keyboard shortcuts

## Behavior on Page Change

- When navigating to a page, its `Init()` command is re-executed to refresh the data.
- Data from the previous page remains cached (no unnecessary reload when returning).

## Multi-Server Switching

Pressing `Ctrl+S` opens a server picker overlay listing all configured servers. Selecting a server reconnects the client and refreshes the current page.

## Language Switching

Pressing `Ctrl+L` toggles between English and French. All UI labels, messages, and keyboard hints update immediately without restarting the application.

## Related Files

- `internal/tui/app.go` — Page router, global key handling
- `internal/tui/app_test.go` — Navigation tests
- `internal/tui/header.go` — Header rendering with tabs
- `internal/tui/footer.go` — Footer rendering with shortcuts
- `internal/i18n/` — Internationalization strings (EN/FR)
