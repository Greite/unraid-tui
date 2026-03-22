package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Greite/unraid-tui/internal/model"
)

var (
	ErrUnauthorized     = errors.New("unauthorized: invalid or missing API key")
	ErrConnectionFailed = errors.New("connection failed: unable to reach Unraid server")
)

// UnraidClient defines operations against the Unraid GraphQL API.
type UnraidClient interface {
	GetSystemInfo(ctx context.Context) (*model.SystemInfo, error)
	GetSystemMetrics(ctx context.Context) (*model.SystemMetrics, error)
	GetVMs(ctx context.Context) ([]model.VM, error)
	GetNotifications(ctx context.Context) ([]model.Notification, error)
	GetContainers(ctx context.Context) ([]model.Container, error)
	StartContainer(ctx context.Context, id string) error
	StopContainer(ctx context.Context, id string) error
	PauseContainer(ctx context.Context, id string) error
	UnpauseContainer(ctx context.Context, id string) error
	UpdateContainer(ctx context.Context, id string) error
	UpdateAllContainers(ctx context.Context) error
	GetDisks(ctx context.Context) ([]model.Disk, error)
	GetShares(ctx context.Context) ([]model.Share, error)
	GetArrayInfo(ctx context.Context) (*model.ArrayInfo, error)
	GetNotificationsOverview(ctx context.Context) (*model.NotificationOverview, error)
	GetNetwork(ctx context.Context) ([]model.NetworkAccess, error)
	ServerURL() string
}

type httpClient struct {
	endpoint string
	origin   string
	apiKey   string
	http     *http.Client
}

// NewClient creates a new Unraid API client.
func NewClient(serverURL, apiKey string) UnraidClient {
	return &httpClient{
		endpoint: serverURL + "/graphql",
		origin:   serverURL,
		apiKey:   apiKey,
		http:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *httpClient) GetSystemInfo(ctx context.Context) (*model.SystemInfo, error) {
	var result graphqlResponse[systemInfoData]
	if err := c.doQuery(ctx, querySystemInfo, &result); err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("graphql: %s", result.Errors[0].Message)
	}
	return result.Data.Info.toDomain(), nil
}

func (c *httpClient) GetSystemMetrics(ctx context.Context) (*model.SystemMetrics, error) {
	var result graphqlResponse[systemMetricsData]
	if err := c.doQuery(ctx, querySystemMetrics, &result); err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("graphql: %s", result.Errors[0].Message)
	}
	return result.Data.Metrics.toDomain(), nil
}

func (c *httpClient) GetVMs(ctx context.Context) ([]model.VM, error) {
	var result graphqlResponse[vmsData]
	if err := c.doQuery(ctx, queryVMs, &result); err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("graphql: %s", result.Errors[0].Message)
	}
	return vmsToDomain(result.Data.VMs.Domains), nil
}

func (c *httpClient) GetNotifications(ctx context.Context) ([]model.Notification, error) {
	var result graphqlResponse[notificationsListData]
	if err := c.doQuery(ctx, queryNotificationsList, &result); err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("graphql: %s", result.Errors[0].Message)
	}
	return notifsToDomain(result.Data.Notifications.List), nil
}

func (c *httpClient) GetContainers(ctx context.Context) ([]model.Container, error) {
	var result graphqlResponse[dockerData]
	if err := c.doQuery(ctx, queryContainers, &result); err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("graphql: %s", result.Errors[0].Message)
	}
	return containersToDomain(result.Data.Docker.Containers), nil
}

func (c *httpClient) StartContainer(ctx context.Context, id string) error {
	return c.doContainerAction(ctx, mutationStartContainer, id)
}

func (c *httpClient) StopContainer(ctx context.Context, id string) error {
	return c.doContainerAction(ctx, mutationStopContainer, id)
}

func (c *httpClient) PauseContainer(ctx context.Context, id string) error {
	return c.doContainerAction(ctx, mutationPauseContainer, id)
}

func (c *httpClient) UnpauseContainer(ctx context.Context, id string) error {
	return c.doContainerAction(ctx, mutationUnpauseContainer, id)
}

func (c *httpClient) UpdateContainer(ctx context.Context, id string) error {
	return c.doContainerAction(ctx, mutationUpdateContainer, id)
}

func (c *httpClient) UpdateAllContainers(ctx context.Context) error {
	return c.doContainerAction(ctx, mutationUpdateAllContainers, "")
}

func (c *httpClient) doContainerAction(ctx context.Context, mutation, id string) error {
	query := mutation
	if id != "" {
		query = fmt.Sprintf(mutation, id)
	}
	body, err := json.Marshal(graphqlRequest{Query: query})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("Origin", c.origin)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrConnectionFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp graphqlResponse[json.RawMessage]
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && len(errResp.Errors) > 0 {
			return fmt.Errorf("graphql: %s", errResp.Errors[0].Message)
		}
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var result graphqlResponse[json.RawMessage]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}
	if len(result.Errors) > 0 {
		return fmt.Errorf("graphql: %s", result.Errors[0].Message)
	}
	return nil
}

func (c *httpClient) GetShares(ctx context.Context) ([]model.Share, error) {
	var result graphqlResponse[sharesData]
	if err := c.doQuery(ctx, queryShares, &result); err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("graphql: %s", result.Errors[0].Message)
	}
	return sharesToDomain(result.Data.Shares), nil
}

func (c *httpClient) GetArrayInfo(ctx context.Context) (*model.ArrayInfo, error) {
	var result graphqlResponse[arrayStateData]
	if err := c.doQuery(ctx, queryArrayState, &result); err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("graphql: %s", result.Errors[0].Message)
	}
	a := result.Data.Array
	return &model.ArrayInfo{
		State:          a.State,
		Free:           a.Capacity.Kilobytes.Free * 1024,
		Used:           a.Capacity.Kilobytes.Used * 1024,
		Total:          a.Capacity.Kilobytes.Total * 1024,
		ParityStatus:   a.ParityCheckStatus.Status,
		ParityProgress: a.ParityCheckStatus.Progress,
		ParityRunning:  a.ParityCheckStatus.Running,
	}, nil
}

func (c *httpClient) GetNotificationsOverview(ctx context.Context) (*model.NotificationOverview, error) {
	var result graphqlResponse[notificationsOverviewData]
	if err := c.doQuery(ctx, queryNotificationsOverview, &result); err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("graphql: %s", result.Errors[0].Message)
	}
	u := result.Data.Notifications.Overview.Unread
	return &model.NotificationOverview{
		Info: u.Info, Warning: u.Warning, Alert: u.Alert, Total: u.Total,
	}, nil
}

func (c *httpClient) GetNetwork(ctx context.Context) ([]model.NetworkAccess, error) {
	var result graphqlResponse[networkData]
	if err := c.doQuery(ctx, queryNetwork, &result); err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("graphql: %s", result.Errors[0].Message)
	}
	return networkToDomain(result.Data.Network.AccessUrls), nil
}

func (c *httpClient) GetDisks(ctx context.Context) ([]model.Disk, error) {
	var result graphqlResponse[disksData]
	if err := c.doQuery(ctx, queryDisks, &result); err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("graphql: %s", result.Errors[0].Message)
	}
	return allDisksToDomain(result.Data.Array), nil
}

func (c *httpClient) ServerURL() string {
	return c.origin
}

func (c *httpClient) doQuery(ctx context.Context, query string, dest any) error {
	body, err := json.Marshal(graphqlRequest{Query: query})
	if err != nil {
		return fmt.Errorf("marshaling query: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("x-api-key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrConnectionFailed, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return ErrUnauthorized
	case http.StatusOK:
		// ok
	case http.StatusBadRequest:
		// GraphQL servers return error details in the body on 400
		var errResp graphqlResponse[json.RawMessage]
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && len(errResp.Errors) > 0 {
			return fmt.Errorf("graphql: %s", errResp.Errors[0].Message)
		}
		return fmt.Errorf("bad request (400): query rejected by server")
	default:
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}
	return nil
}
