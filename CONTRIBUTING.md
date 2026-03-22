# Contributing to unraid-tui

Thank you for your interest in this project! This guide explains how to contribute.

## Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [GNU Make](https://www.gnu.org/software/make/)
- (Optional) [GoReleaser](https://goreleaser.com/) to test releases
- (Optional) [GitHub CLI](https://cli.github.com/) for PRs

## Setup

```bash
git clone https://github.com/Greite/unraid-tui.git
cd unraid-tui
make build
make test
```

If the build and tests pass, you are ready to go.

## Development workflow

### 1. Create a branch

```bash
git checkout -b feat/my-feature
```

Branch naming conventions:

| Prefix     | Usage                |
|------------|----------------------|
| `feat/`    | New feature          |
| `fix/`     | Bug fix              |
| `refactor/`| Refactoring          |
| `docs/`    | Documentation only   |
| `test/`    | Add/modify tests     |

### 2. Develop

```bash
# Run tests continuously during development
make test

# Verify the code compiles
make build

# Run the linter
make lint
```

Configuration is stored at `~/.unraid-tui/config.yaml`.

### 3. Commit

Commit messages must be in English, using the present imperative tense:

```bash
# Good
git commit -m "add VM monitoring page"
git commit -m "fix container port display when host port is 0"

# Bad
git commit -m "Added VM page"
git commit -m "WIP"
git commit -m "fix stuff"
```

### 4. Open a PR

```bash
git push -u origin feat/my-feature
gh pr create --title "Add VM monitoring page" --body "Description..."
```

## Project structure

```
cmd/                        -> Cobra entry point
internal/api/               -> GraphQL client (interface + HTTP)
internal/config/            -> Viper configuration
internal/i18n/              -> Internationalization (i18n) system
internal/model/             -> Domain types
internal/tui/               -> Bubbletea app (router)
internal/tui/common/        -> Shared styles, messages, helpers
internal/tui/dashboard/     -> Dashboard page
internal/tui/docker/        -> Docker page
internal/tui/vms/           -> VMs page
internal/tui/notifications/ -> Notifications page
internal/tui/shares/        -> Shares page
internal/tui/onboarding/    -> Configuration wizard
```

## Adding a new TUI page

1. Create a package in `internal/tui/<name>/`
2. Implement a `Model` with `Init()`, `Update()`, `View()` (Bubbletea pattern)
3. Import styles and messages from `internal/tui/common/` (never from `internal/tui/` directly to avoid import cycles)
4. Add the page to the router in `internal/tui/app.go`
5. Add GraphQL queries in `internal/api/queries.go` if needed
6. Add all user-facing strings to the i18n system (see below)
7. Write tests

## Internationalization (i18n)

The project uses an i18n system located in `internal/i18n/`. Every user-facing string displayed in the TUI must go through this system.

When adding new strings:

1. Add the string key and its English translation in `internal/i18n/`
2. Provide translations for all supported languages
3. Tests use English strings (the default language), so always ensure the English translation is present

## Adding an API query

1. Add the GraphQL query in `internal/api/queries.go`
2. Add response types in `internal/api/types.go` with a `toDomain()` method
3. Add domain types in `internal/model/model.go`
4. Add the method to the `UnraidClient` interface (`internal/api/client.go`)
5. Implement the method in `httpClient`
6. Add the method to `MockClient` (`internal/api/mock.go`)
7. Test with `httptest` in `internal/api/client_test.go`

## Multi-server support

The project supports connecting to multiple Unraid servers. When working on features, keep in mind that all API calls and state management must handle the multi-server context. Server configuration is managed in `~/.unraid-tui/config.yaml`.

## Tests

Tests are mandatory for every contribution.

```bash
make test          # Run tests
make test-verbose  # With verbose output
make test-cover    # With coverage
```

### Rules

- Use the standard `testing` package (no external framework)
- Mock the API with `httptest.NewServer` for client tests
- Use `api.MockClient` for TUI tests
- Test the Bubbletea model (Update/View), not the terminal rendering
- Tests use English strings (the default language)
- Every new feature must include tests

### Run a specific test

```bash
go test -v -run TestMyTestName ./internal/tui/docker/
```

## Code conventions

- **Charmbracelet v2 import paths**: `charm.land/bubbletea/v2`, `charm.land/lipgloss/v2`, `charm.land/bubbles/v2`
- **Bubbletea v2**:
  - `View()` returns `tea.View`, not `string`
  - No `tea.WithAltScreen()` -- use `v.AltScreen = true` in the View
  - Spinner: use `m.spinner.Tick` as the initial Cmd (not via `Init()`)
- **No import cycles**: TUI sub-packages import `tui/common`, never `tui`
- **`UnraidClient` interface**: all API access goes through this interface
- Format code with `gofmt` (applied automatically by Go)

## Documentation

- Document new features in `docs/`
- Update `README.md` if the feature is user-facing
- Update `CLAUDE.md` if the contribution changes architecture or conventions

## Reporting a bug

Open an issue on [GitHub](https://github.com/Greite/unraid-tui/issues) with:

- The version (`unraid-tui version`)
- Your OS and architecture
- Steps to reproduce
- Expected vs observed behavior
