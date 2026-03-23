package api

import (
	"context"

	"github.com/Greite/unraid-tui/internal/model"
)

// MockClient is a test double for UnraidClient.
type MockClient struct {
	GetSystemInfoFn            func(ctx context.Context) (*model.SystemInfo, error)
	GetSystemMetricsFn         func(ctx context.Context) (*model.SystemMetrics, error)
	GetVMsFn                   func(ctx context.Context) ([]model.VM, error)
	GetNotificationsFn         func(ctx context.Context) ([]model.Notification, error)
	GetContainersFn            func(ctx context.Context) ([]model.Container, error)
	GetContainerStatsFn        func(ctx context.Context) (map[string]model.Container, error)
	GetDisksFn                 func(ctx context.Context) ([]model.Disk, error)
	GetSharesFn                func(ctx context.Context) ([]model.Share, error)
	GetArrayInfoFn             func(ctx context.Context) (*model.ArrayInfo, error)
	GetParityHistoryFn         func(ctx context.Context) ([]model.ParityHistoryEntry, error)
	GetNotificationsOverviewFn func(ctx context.Context) (*model.NotificationOverview, error)
	GetNetworkFn               func(ctx context.Context) ([]model.NetworkAccess, error)
}

func (m *MockClient) GetSystemInfo(ctx context.Context) (*model.SystemInfo, error) {
	if m.GetSystemInfoFn != nil {
		return m.GetSystemInfoFn(ctx)
	}
	return nil, nil
}
func (m *MockClient) GetSystemInfoExtra(_ context.Context, _ *model.SystemInfo) error { return nil }
func (m *MockClient) GetSystemMetrics(ctx context.Context) (*model.SystemMetrics, error) {
	if m.GetSystemMetricsFn != nil {
		return m.GetSystemMetricsFn(ctx)
	}
	return nil, nil
}
func (m *MockClient) GetVMs(ctx context.Context) ([]model.VM, error) {
	if m.GetVMsFn != nil {
		return m.GetVMsFn(ctx)
	}
	return nil, nil
}
func (m *MockClient) GetNotifications(ctx context.Context) ([]model.Notification, error) {
	if m.GetNotificationsFn != nil {
		return m.GetNotificationsFn(ctx)
	}
	return nil, nil
}
func (m *MockClient) GetContainers(ctx context.Context) ([]model.Container, error) {
	if m.GetContainersFn != nil {
		return m.GetContainersFn(ctx)
	}
	return nil, nil
}
func (m *MockClient) GetContainerStats(ctx context.Context) (map[string]model.Container, error) {
	if m.GetContainerStatsFn != nil {
		return m.GetContainerStatsFn(ctx)
	}
	return nil, nil
}
func (m *MockClient) GetDisks(ctx context.Context) ([]model.Disk, error) {
	if m.GetDisksFn != nil {
		return m.GetDisksFn(ctx)
	}
	return nil, nil
}
func (m *MockClient) GetShares(ctx context.Context) ([]model.Share, error) {
	if m.GetSharesFn != nil {
		return m.GetSharesFn(ctx)
	}
	return nil, nil
}
func (m *MockClient) GetArrayInfo(ctx context.Context) (*model.ArrayInfo, error) {
	if m.GetArrayInfoFn != nil {
		return m.GetArrayInfoFn(ctx)
	}
	return nil, nil
}
func (m *MockClient) GetParityHistory(ctx context.Context) ([]model.ParityHistoryEntry, error) {
	if m.GetParityHistoryFn != nil {
		return m.GetParityHistoryFn(ctx)
	}
	return nil, nil
}
func (m *MockClient) GetNotificationsOverview(ctx context.Context) (*model.NotificationOverview, error) {
	if m.GetNotificationsOverviewFn != nil {
		return m.GetNotificationsOverviewFn(ctx)
	}
	return nil, nil
}
func (m *MockClient) GetNetwork(ctx context.Context) ([]model.NetworkAccess, error) {
	if m.GetNetworkFn != nil {
		return m.GetNetworkFn(ctx)
	}
	return nil, nil
}
func (m *MockClient) StartContainer(_ context.Context, _ string) error   { return nil }
func (m *MockClient) StopContainer(_ context.Context, _ string) error    { return nil }
func (m *MockClient) PauseContainer(_ context.Context, _ string) error   { return nil }
func (m *MockClient) UnpauseContainer(_ context.Context, _ string) error { return nil }
func (m *MockClient) UpdateContainer(_ context.Context, _ string) error  { return nil }
func (m *MockClient) UpdateAllContainers(_ context.Context) error        { return nil }
func (m *MockClient) SetAutostart(_ context.Context, _ []model.Container, _ string, _ bool) error {
	return nil
}
func (m *MockClient) StartVM(_ context.Context, _ string) error             { return nil }
func (m *MockClient) StopVM(_ context.Context, _ string) error              { return nil }
func (m *MockClient) PauseVM(_ context.Context, _ string) error             { return nil }
func (m *MockClient) ResumeVM(_ context.Context, _ string) error            { return nil }
func (m *MockClient) ForceStopVM(_ context.Context, _ string) error         { return nil }
func (m *MockClient) RebootVM(_ context.Context, _ string) error            { return nil }
func (m *MockClient) ArchiveNotification(_ context.Context, _ string) error { return nil }
func (m *MockClient) ArchiveAllNotifications(_ context.Context) error       { return nil }
func (m *MockClient) ServerURL() string                                     { return "http://localhost" }
