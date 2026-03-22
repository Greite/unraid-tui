# Configuration

## Configuration File

The application looks for a configuration file at `~/.unraid-tui/config.yaml` in the user's home directory.

### Format

```yaml
servers:
  - name: "tower"
    url: "http://192.168.1.100"
  - name: "backup"
    url: "http://192.168.1.200"

language: "en"
```

API keys are **not** stored in the configuration file. They are stored securely in the system keychain (macOS Keychain, Linux secret-service, Windows Credential Manager), keyed by server name.

### Parameters

| Parameter  | Required | Description                                           |
|------------|----------|-------------------------------------------------------|
| `servers`  | Yes      | List of named Unraid servers                          |
| `name`     | Yes      | Friendly name for the server                          |
| `url`      | Yes      | URL of the Unraid server (including port)             |
| `language` | No       | UI language (`en` or `fr`). Defaults to `en`.         |

### Multi-Server Support

The configuration supports multiple named servers. Use `Ctrl+S` in the TUI to switch between configured servers. Each server has its own API key stored independently in the system keychain.

## Environment Variables

Each parameter can be overridden by an environment variable prefixed with `UNRAID_`:

| Variable              | Overrides            |
|-----------------------|----------------------|
| `UNRAID_SERVER_URL`   | Primary server URL   |
| `UNRAID_API_KEY`      | Primary server key   |

Environment variables take priority over the configuration file.

### Example

```bash
export UNRAID_SERVER_URL="http://10.0.0.5"
export UNRAID_API_KEY="my-secret-key"
unraid-tui
```

## API Key Storage

API keys are stored in the system keychain for security. During onboarding, the key is saved to the keychain automatically. You can also manage keys manually:

```bash
# The application handles keychain storage transparently
# Keys are stored under the service name "unraid-tui" with the server name as account
```

## Obtaining an API Key

1. Open the Unraid web interface
2. Go to **Settings > Management Access > Developer Options**
3. Open Apollo GraphQL Studio
4. Create an API key via the GraphQL mutation:

```graphql
mutation {
  apiKey {
    create(input: {
      name: "unraid-tui"
      description: "CLI monitoring tool"
      roles: [VIEWER]
      permissions: [
        { resource: INFO, actions: [READ_ANY] }
        { resource: DOCKER, actions: [READ_ANY, WRITE_ANY] }
        { resource: VMS, actions: [READ_ANY, WRITE_ANY] }
        { resource: SHARES, actions: [READ_ANY] }
        { resource: NOTIFICATIONS, actions: [READ_ANY, WRITE_ANY] }
      ]
    }) {
      key
    }
  }
}
```

5. The key will be stored in your system keychain during onboarding

## Internationalization (i18n)

The application supports multiple languages:

- **English** (`en`) — Default
- **French** (`fr`)

Switch language at any time using `Ctrl+L` in the TUI, or set it in the configuration file.

## Validation

At startup, the application verifies that at least one server is configured with a URL and that the corresponding API key exists in the keychain. If anything is missing, the onboarding wizard is launched.

## Related Files

- `internal/config/config.go` — Loading and validation
- `internal/config/config_test.go` — Tests (file, env vars, override)
