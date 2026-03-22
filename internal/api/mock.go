package api

import (
	"context"

	"github.com/Greite/unraid-tui/internal/model"
)

// MockClient is a test double for UnraidClient.
type MockClient struct {
	GetSystemInfoFn    func(ctx context.Context) (*model.SystemInfo, error)
	GetSystemMetricsFn func(ctx context.Context) (*model.SystemMetrics, error)
	GetVMsFn            func(ctx context.Context) ([]model.VM, error)
	GetNotificationsFn  func(ctx context.Context) ([]model.Notification, error)
	GetContainersFn    func(ctx context.Context) ([]model.Container, error)
	GetDisksFn         func(ctx context.Context) ([]model.Disk, error)
	StartContainerFn           func(ctx context.Context, id string) error
	StopContainerFn            func(ctx context.Context, id string) error
	PauseContainerFn           func(ctx context.Context, id string) error
	UnpauseContainerFn         func(ctx context.Context, id string) error
	UpdateContainerFn          func(ctx context.Context, id string) error
	UpdateAllContainersFn      func(ctx context.Context) error
	GetSharesFn                func(ctx context.Context) ([]model.Share, error)
	GetArrayInfoFn             func(ctx context.Context) (*model.ArrayInfo, error)
	GetNotificationsOverviewFn func(ctx context.Context) (*model.NotificationOverview, error)
	GetNetworkFn               func(ctx context.Context) ([]model.NetworkAccess, error)
}

func (m *MockClient) GetSystemInfo(ctx context.Context) (*model.SystemInfo, error) {
	if m.GetSystemInfoFn != nil {
		return m.GetSystemInfoFn(ctx)
	}
	return nil, nil
}

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

func (m *MockClient) GetDisks(ctx context.Context) ([]model.Disk, error) {
	if m.GetDisksFn != nil {
		return m.GetDisksFn(ctx)
	}
	return nil, nil
}

func (m *MockClient) StartContainer(ctx context.Context, id string) error {
	if m.StartContainerFn != nil {
		return m.StartContainerFn(ctx, id)
	}
	return nil
}
func (m *MockClient) StopContainer(ctx context.Context, id string) error {
	if m.StopContainerFn != nil {
		return m.StopContainerFn(ctx, id)
	}
	return nil
}
func (m *MockClient) PauseContainer(ctx context.Context, id string) error {
	if m.PauseContainerFn != nil {
		return m.PauseContainerFn(ctx, id)
	}
	return nil
}
func (m *MockClient) UnpauseContainer(ctx context.Context, id string) error {
	if m.UnpauseContainerFn != nil {
		return m.UnpauseContainerFn(ctx, id)
	}
	return nil
}

func (m *MockClient) UpdateContainer(ctx context.Context, id string) error {
	if m.UpdateContainerFn != nil {
		return m.UpdateContainerFn(ctx, id)
	}
	return nil
}
func (m *MockClient) UpdateAllContainers(ctx context.Context) error {
	if m.UpdateAllContainersFn != nil {
		return m.UpdateAllContainersFn(ctx)
	}
	return nil
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

func (m *MockClient) ServerURL() string {
	return "http://localhost"
}
