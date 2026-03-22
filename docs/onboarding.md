# Onboarding

The onboarding wizard is an interactive configuration assistant that launches automatically on first startup when no configuration is detected.

## Trigger

The onboarding launches if:
- The file `~/.unraid-tui/config.yaml` does not exist
- **And** the environment variables `UNRAID_SERVER_URL` / `UNRAID_API_KEY` are not defined

If the configuration exists (file or env vars), the dashboard launches directly.

## Steps

The onboarding consists of 4 steps with a progress bar at the top of the screen.

### Step 1: Server Name

```
╭────────────────────────────────────────────────────────╮
│  Welcome!                                              │
│                                                        │
│  This wizard will help you configure the connection    │
│  to your Unraid server in a few steps.                 │
│                                                        │
│  Enter a friendly name for your server:                │
│                                                        │
│  > tower                                               │
│                                                        │
╰────────────────────────────────────────────────────────╯

  enter next  esc back
```

The user enters a friendly name to identify the server (used for multi-server switching with `Ctrl+S`).

Validation: the name cannot be empty.

### Step 2: Server URL

The user enters the URL of their Unraid server. The input accepts several formats:

| Input                         | Normalized to                 |
|-------------------------------|-------------------------------|
| `192.168.1.100:3001`          | `http://192.168.1.100:3001`   |
| `http://tower:3001`           | `http://tower:3001`           |
| `https://secure.local:3001/`  | `https://secure.local:3001`   |

Automatic normalization:
- Adds `http://` if no scheme is present
- Removes trailing `/`

Validation: the URL cannot be empty.

A connection test is performed by sending a minimal GraphQL query (`{ __typename }`) to the server. The test only verifies that the server is reachable (even a 401 response is considered a success at this stage -- it means the API is present).

- **Success**: proceed to the next step
- **Failure**: return to URL input with the error message
- **Timeout**: 5 seconds

### Step 3: API Key Instructions

Displays detailed instructions for creating an API key from the Unraid web interface:

1. Open the Unraid web interface
2. Settings > Management Access > Developer Options
3. Open Apollo GraphQL Studio
4. Execute the `apiKey.create` mutation
5. Copy the returned key

The exact GraphQL mutation is displayed in the interface.

### Step 4: API Key Entry

Masked input field (password mode -- characters are replaced with `*`).

Validation: the key cannot be empty.

An authenticated request (`info { os { hostname } }`) is sent to verify that the key works:

- **200 OK**: the key is valid, proceed to save
- **401/403**: invalid key or insufficient permissions, return to input
- **Other**: error displayed, return to input

## Save

On successful validation, the configuration is saved:

1. The server entry (name + URL) is written to `~/.unraid-tui/config.yaml` with permissions `0600` (owner read/write only)
2. The API key is stored in the **system keychain** (macOS Keychain, Linux secret-service, Windows Credential Manager), not in the config file

### Confirmation Screen

```
╭────────────────────────────────────────────────────────╮
│  Configuration complete!                               │
│                                                        │
│  Your configuration has been saved:                    │
│    Config: ~/.unraid-tui/config.yaml                   │
│    API key: stored in system keychain                  │
│                                                        │
│  Server: tower (http://192.168.1.100:3001)             │
│                                                        │
│  The dashboard will now launch.                        │
╰────────────────────────────────────────────────────────╯

  enter launch dashboard
```

After confirmation, the main dashboard launches automatically.

## Navigation

| Key      | Action                                  |
|----------|-----------------------------------------|
| `Enter`  | Validate / Proceed to next step         |
| `Esc`    | Return to previous step                 |
| `Ctrl+C` | Cancel and quit                        |

## Progress Bar

A progress bar at the top of the screen indicates the current step:

```
  ● Name  —  ◉ URL  —  ○ API Info  —  ○ API Key
```

- `●` completed step (green)
- `◉` current step (purple)
- `○` upcoming step (gray)

## Error Handling

Errors are displayed in red below the step content. The user remains on the current step and can correct their input.

Possible errors:
- Empty name
- Empty URL
- Server unreachable (timeout, DNS, connection refused)
- Invalid URL
- Empty API key
- Invalid API key (401/403)
- File save error
- Keychain access error

## Related Files

- `internal/tui/onboarding/onboarding.go` — Multi-step Bubbletea model
- `internal/tui/onboarding/onboarding_test.go` — Unit tests
- `internal/config/config.go` — `Exists()`, `Save()`, `FilePath()`
- `internal/config/keychain.go` — System keychain integration
- `cmd/root.go` — Onboarding detection and launch
