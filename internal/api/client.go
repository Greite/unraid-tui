package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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
	GetSystemInfoExtra(ctx context.Context, info *model.SystemInfo) error
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
	GetContainerStats(ctx context.Context) (map[string]model.Container, error)
	GetParityHistory(ctx context.Context) ([]model.ParityHistoryEntry, error)
	SetAutostart(ctx context.Context, containers []model.Container, targetID string, autoStart bool) error
	// VM mutations
	StartVM(ctx context.Context, id string) error
	StopVM(ctx context.Context, id string) error
	PauseVM(ctx context.Context, id string) error
	ResumeVM(ctx context.Context, id string) error
	ForceStopVM(ctx context.Context, id string) error
	RebootVM(ctx context.Context, id string) error
	// Notification mutations
	ArchiveNotification(ctx context.Context, id string) error
	ArchiveAllNotifications(ctx context.Context) error
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

func (c *httpClient) GetSystemInfoExtra(ctx context.Context, info *model.SystemInfo) error {
	var result graphqlResponse[systemInfoExtraData]
	if err := c.doQuery(ctx, querySystemInfoExtra, &result); err != nil {
		return err
	}
	if len(result.Errors) > 0 {
		return nil // silently ignore — extra fields may not be available
	}
	result.Data.Info.applyTo(info)
	return nil
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
	return containersToDomain(result.Data.Docker.Containers, result.Data.Docker.ContainerUpdateStatuses), nil
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
	kb := a.Capacity.Kilobytes
	free, _ := kb.Free.Int64()
	used, _ := kb.Used.Int64()
	total, _ := kb.Total.Int64()
	return &model.ArrayInfo{
		State:          a.State,
		Free:           uint64(free) * 1024,
		Used:           uint64(used) * 1024,
		Total:          uint64(total) * 1024,
		ParityStatus:   a.ParityCheckStatus.Status,
		ParityProgress: a.ParityCheckStatus.Progress,
		ParityRunning:  a.ParityCheckStatus.Running,
		ParityDate:     a.ParityCheckStatus.Date,
		ParityDuration: a.ParityCheckStatus.Duration,
		ParitySpeed:    a.ParityCheckStatus.Speed,
		ParityErrors:   a.ParityCheckStatus.Errors,
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

func (c *httpClient) GetContainerStats(ctx context.Context) (map[string]model.Container, error) {
	var result graphqlResponse[containerStatsData]
	if err := c.doQuery(ctx, queryContainerStats, &result); err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("graphql: %s", result.Errors[0].Message)
	}
	stats := make(map[string]model.Container)
	for _, s := range result.Data.Docker.Containers {
		name := ""
		if len(s.Names) > 0 {
			name = strings.TrimPrefix(s.Names[0], "/")
		}
		stats[s.ID] = model.Container{
			ID: s.ID, Name: name,
			CPUPercent: s.CPUPercent, MemUsage: s.MemUsage, MemPercent: s.MemPercent,
		}
	}
	return stats, nil
}

func (c *httpClient) GetParityHistory(ctx context.Context) ([]model.ParityHistoryEntry, error) {
	var result graphqlResponse[parityHistoryData]
	if err := c.doQuery(ctx, queryParityHistory, &result); err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("graphql: %s", result.Errors[0].Message)
	}
	entries := make([]model.ParityHistoryEntry, len(result.Data.Array.ParityHistory))
	for i, p := range result.Data.Array.ParityHistory {
		entries[i] = model.ParityHistoryEntry{
			Date: p.Date, Status: p.Status, Duration: p.Duration,
			Speed: p.Speed, Errors: p.Errors,
		}
	}
	return entries, nil
}

func (c *httpClient) SetAutostart(ctx context.Context, containers []model.Container, targetID string, autoStart bool) error {
	var entries []string
	for _, ct := range containers {
		as := ct.AutoStart
		if ct.ID == targetID {
			as = autoStart
		}
		entries = append(entries, fmt.Sprintf(`{ id: "%s", autoStart: %t, wait: 0 }`, ct.ID, as))
	}
	query := mutationAutostartPrefix + strings.Join(entries, ", ") + mutationAutostartSuffix
	return c.doContainerAction(ctx, query, "")
}

func (c *httpClient) StartVM(ctx context.Context, id string) error {
	return c.doContainerAction(ctx, mutationVMStart, id)
}
func (c *httpClient) StopVM(ctx context.Context, id string) error {
	return c.doContainerAction(ctx, mutationVMStop, id)
}
func (c *httpClient) PauseVM(ctx context.Context, id string) error {
	return c.doContainerAction(ctx, mutationVMPause, id)
}
func (c *httpClient) ResumeVM(ctx context.Context, id string) error {
	return c.doContainerAction(ctx, mutationVMResume, id)
}
func (c *httpClient) ForceStopVM(ctx context.Context, id string) error {
	return c.doContainerAction(ctx, mutationVMForceStop, id)
}
func (c *httpClient) RebootVM(ctx context.Context, id string) error {
	return c.doContainerAction(ctx, mutationVMReboot, id)
}

func (c *httpClient) ArchiveNotification(ctx context.Context, id string) error {
	return c.doContainerAction(ctx, mutationArchiveNotification, id)
}
func (c *httpClient) ArchiveAllNotifications(ctx context.Context) error {
	return c.doContainerAction(ctx, mutationArchiveAllNotifications, "")
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
