# API Client

The API client encapsulates all communication with the Unraid GraphQL API.

## Architecture

```
UnraidClient (interface)
    │
    ├── httpClient (real implementation)
    │   └── sends HTTP POST requests to /graphql
    │
    └── MockClient (tests)
        └── configurable functions per test
```

### Interface

```go
type UnraidClient interface {
    GetSystemInfo(ctx context.Context) (*model.SystemInfo, error)
    GetSystemMetrics(ctx context.Context) (*model.SystemMetrics, error)
    GetContainers(ctx context.Context) ([]model.Container, error)
    GetContainerStats(ctx context.Context, id string) (*model.ContainerStats, error)
    GetVMs(ctx context.Context) ([]model.VM, error)
    GetShares(ctx context.Context) ([]model.Share, error)
    GetArrayInfo(ctx context.Context) (*model.ArrayInfo, error)
    GetParityHistory(ctx context.Context) ([]model.ParityEvent, error)

    // Mutations
    StartContainer(ctx context.Context, id string) error
    StopContainer(ctx context.Context, id string) error
    PauseContainer(ctx context.Context, id string) error
    UpdateContainer(ctx context.Context, id string) error
    SetAutostart(ctx context.Context, id string, enabled bool) error

    StartVM(ctx context.Context, id string) error
    StopVM(ctx context.Context, id string) error
    PauseVM(ctx context.Context, id string) error

    DismissNotification(ctx context.Context, id string) error
    DismissAllNotifications(ctx context.Context) error
}
```

All TUI code depends on this interface, never on the concrete implementation.

## HTTP Implementation

### Requests

Each method sends a POST request to `<server_url>/graphql` with:

- **Header** `Content-Type: application/json`
- **Header** `Authorization: Bearer <api_key>`
- **Body**: `{"query": "<graphql_query>"}`

The API key is retrieved from the system keychain at runtime.

### Responses

GraphQL responses are decoded into Go-specific structs (package `api/types.go`), then converted to domain types (`model/model.go`) via `toDomain()` methods.

### Error Handling

| Situation                | Error returned              |
|--------------------------|-----------------------------|
| Server unreachable       | `ErrConnectionFailed`       |
| 401 or 403 response      | `ErrUnauthorized`           |
| Unexpected HTTP code      | `unexpected status <code>`  |
| Error in GraphQL JSON    | `graphql: <message>`        |
| Invalid JSON             | `decoding response: <err>`  |

The sentinel errors `ErrConnectionFailed` and `ErrUnauthorized` allow specific handling on the TUI side (adapted display, no retry on auth errors).

## Mock for Tests

`api/mock.go` provides a `MockClient` with replaceable function fields:

```go
mock := &api.MockClient{
    GetContainersFn: func(ctx context.Context) ([]model.Container, error) {
        return []model.Container{{Name: "test"}}, nil
    },
}
```

If a function is not defined, the mock returns `(nil, nil)`.

## GraphQL Queries

Queries are string constants in `api/queries.go`. They correspond to the Unraid GraphQL API:

| Constant               | Data retrieved                                       |
|------------------------|------------------------------------------------------|
| `querySystemInfo`      | CPU, memory, OS, baseboard                           |
| `querySystemMetrics`   | CPU usage %, memory usage %, CPU temp/power/per-core |
| `queryContainers`      | List of containers with ports and networks           |
| `queryContainerStats`  | Per-container CPU and memory statistics              |
| `queryVMs`             | List of VMs with state and configuration             |
| `queryShares`          | List of shares with size and access settings         |
| `queryArrayInfo`       | Array disks, cache, parity status and capacity       |
| `queryParityHistory`   | Parity check history with dates and results          |

## Tests

Tests cover the HTTP client:

- `TestGetSystemInfo_Success` — Verifies headers, parsing, domain conversion
- `TestGetSystemMetrics_Success` — Metrics parsing
- `TestGetContainers_Success` — Multi-container parsing with ports
- `TestGetSystemInfo_Unauthorized` — 401 detection
- `TestGetSystemInfo_GraphQLError` — Error in GraphQL response
- `TestGetSystemInfo_ConnectionError` — Unreachable server

## Related Files

- `internal/api/client.go` — Interface + HTTP implementation
- `internal/api/client_test.go` — Tests with httptest
- `internal/api/queries.go` — GraphQL queries
- `internal/api/types.go` — JSON response structs + domain conversion
- `internal/api/mock.go` — Mock client for TUI tests
