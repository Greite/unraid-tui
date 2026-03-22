# Docker

The Docker page displays all Docker containers on the Unraid server in an interactive table.

## Container Table

```
  Containers (5)  3 running                              ⟳ 1 update available

  ┌──────────────┬─────────────────────┬────────────┬──────────────────┬─────────────┐
  │ NAME         │ IMAGE               │ STATE      │ STATUS           │ PORTS       │
  ├──────────────┼─────────────────────┼────────────┼──────────────────┼─────────────┤
  │ plex         │ plexinc/pms:latest  │ ● running  │ Up 14 days       │ 32400:32400 │
  │ nextcloud  ⟳ │ nextcloud:28        │ ● running  │ Up 14 days       │ 443:443     │
  │ homeassistant│ ghcr.io/ha:latest   │ ● running  │ Up 14 days       │ 8123:8123   │
  │ pihole       │ pihole/pihole:lat.. │ ○ exited   │ Exited (0) 2d    │ -           │
  │ test-db      │ postgres:15         │ ○ exited   │ Exited (1) 5d    │ -           │
  └──────────────┴─────────────────────┴────────────┴──────────────────┴─────────────┘

  ↑/↓: navigate  │  s: sort  │  enter: actions  │  r: refresh
```

### Columns

| Column | Description                                           |
|--------|-------------------------------------------------------|
| NAME   | Container name (with ⟳ update indicator if available)  |
| IMAGE  | Docker image (truncated to 25 characters if needed)    |
| STATE  | State with visual indicator (● running, ○ exited, ◑ paused) |
| STATUS | Status detail (e.g., "Up 14 days", "Exited (0) 2 days ago") |
| PORTS  | Host:container port mapping (or `-` if none)           |

### Header

- **Total number** of containers in parentheses
- **Number of running containers** displayed alongside
- **Update indicator** showing how many container updates are available

## Sorting

Press `s` to cycle through sort modes:

- **Name** (alphabetical)
- **State** (running first, then paused, then exited)
- **Image** (alphabetical)

## Container Actions

Press `Enter` on a selected container to open the action menu:

| Action         | Description                                           |
|----------------|-------------------------------------------------------|
| Start          | Start a stopped container                             |
| Stop           | Stop a running container                              |
| Pause / Unpause| Pause or unpause a running container                  |
| Update         | Pull the latest image and recreate the container      |
| Logs           | View container logs with follow mode (live streaming)  |
| SSH Console    | Open an interactive SSH console into the container     |
| WebUI          | Open the container's web interface in the default browser |
| Set Autostart  | Toggle autostart on/off for the container              |

Available actions are context-sensitive based on the container's current state.

## Navigation

| Key     | Action                        |
|---------|-------------------------------|
| `↑`/`↓` | Navigate the table            |
| `Enter` | Open action menu              |
| `s`     | Cycle sort mode               |
| `r`     | Refresh the list              |
| `Esc`   | Close action menu / back      |

The table supports scrolling if the list exceeds the available height.

## Responsive Columns

Column widths adapt automatically to the terminal size. Each column uses a percentage of the available width:

- NAME: 20%
- IMAGE: 25%
- STATE: 10%
- STATUS: 25%
- PORTS: 15%

## States

| State        | Display                                          |
|--------------|--------------------------------------------------|
| Loading      | Animated spinner + "Loading containers..."       |
| Data OK      | Full table with counter                          |
| Error        | Red banner with error message                    |

## GraphQL Query Used

```graphql
query {
  docker {
    containers {
      id name image state status autostart
      ports { privatePort publicPort type }
      networks { networkId }
      updateAvailable
    }
  }
}
```

## Related Files

- `internal/tui/docker/docker.go` — Bubbletea model, Bubbles table, formatting
- `internal/tui/docker/docker_test.go` — Unit tests
- `internal/api/queries.go` — GraphQL query (`queryContainers`)
