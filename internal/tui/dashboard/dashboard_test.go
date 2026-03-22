package dashboard

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Greite/unraid-tui/internal/api"
	"github.com/Greite/unraid-tui/internal/model"
	"github.com/Greite/unraid-tui/internal/tui/common"
)

func newTestDashboard() Model {
	mock := &api.MockClient{
		GetSystemInfoFn: func(ctx context.Context) (*model.SystemInfo, error) {
			return &model.SystemInfo{
				CPU:    model.CPUInfo{Brand: "AMD Ryzen 9", Cores: 16, Threads: 32},
				Memory: model.MemoryInfo{Total: 68719476736, Used: 34359738368},
				OS:     model.OSInfo{Distro: "Unraid", Release: "6.12.6", Platform: "linux", Uptime: 1234567},
			}, nil
		},
		GetSystemMetricsFn: func(ctx context.Context) (*model.SystemMetrics, error) {
			return &model.SystemMetrics{
				CPUUsage:    42.5,
				MemoryUsed:  34359738368,
				MemoryTotal: 68719476736,
				MemoryPct:   50.0,
			}, nil
		},
	}
	return New(mock)
}

func TestDashboard_InitialState(t *testing.T) {
	m := newTestDashboard()
	if !m.loading {
		t.Error("expected loading to be true initially")
	}
	if m.systemInfo != nil {
		t.Error("expected systemInfo to be nil initially")
	}
}

func TestDashboard_SystemInfoMsg(t *testing.T) {
	m := newTestDashboard()
	m.width = 80
	m.height = 24

	info := &model.SystemInfo{
		CPU: model.CPUInfo{Brand: "AMD Ryzen 9", Cores: 16},
		OS:  model.OSInfo{Distro: "Unraid"},
	}
	updated, _ := m.Update(common.SystemInfoMsg{Info: info})

	if updated.loading {
		t.Error("expected loading to be false after SystemInfoMsg")
	}
	if updated.systemInfo == nil {
		t.Fatal("expected systemInfo to be set")
	}
	if updated.systemInfo.CPU.Brand != "AMD Ryzen 9" {
		t.Errorf("expected CPU brand AMD Ryzen 9, got %s", updated.systemInfo.CPU.Brand)
	}
}

func TestDashboard_SystemInfoMsg_Error(t *testing.T) {
	m := newTestDashboard()
	m.width = 80
	m.height = 24

	updated, _ := m.Update(common.SystemInfoMsg{Err: errors.New("connection failed")})

	if updated.loading {
		t.Error("expected loading to be false")
	}
	if updated.err == nil {
		t.Fatal("expected error to be set")
	}
	if updated.err.Error() != "connection failed" {
		t.Errorf("expected 'connection failed', got %s", updated.err.Error())
	}
}

func TestDashboard_MetricsMsg(t *testing.T) {
	m := newTestDashboard()
	m.width = 80
	m.height = 24
	m.loading = false

	metrics := &model.SystemMetrics{CPUUsage: 75.0, MemoryPct: 50.0, MemoryUsed: 32000000000, MemoryTotal: 64000000000}
	updated, cmd := m.Update(common.SystemMetricsMsg{Metrics: metrics})

	if updated.metrics == nil {
		t.Fatal("expected metrics to be set")
	}
	if updated.metrics.CPUUsage != 75.0 {
		t.Errorf("expected CPU usage 75.0, got %f", updated.metrics.CPUUsage)
	}
	if cmd == nil {
		t.Error("expected a refresh command after metrics update")
	}
}

func TestDashboard_MetricsMsg_ErrorStillSchedulesRefresh(t *testing.T) {
	m := newTestDashboard()
	m.width = 80
	m.height = 24
	m.loading = false

	updated, cmd := m.Update(common.SystemMetricsMsg{Err: errors.New("timeout")})

	if updated.err == nil {
		t.Error("expected error to be set")
	}
	if cmd == nil {
		t.Error("expected refresh command even on error")
	}
}

func TestDashboard_View_Loading(t *testing.T) {
	m := newTestDashboard()
	m.width = 80
	m.height = 24
	view := m.View()
	if !strings.Contains(view, "Chargement") {
		t.Error("expected loading message in view")
	}
}

func TestDashboard_View_WithData(t *testing.T) {
	m := newTestDashboard()
	m.width = 80
	m.height = 24
	m.loading = false
	m.systemInfo = &model.SystemInfo{
		CPU: model.CPUInfo{Brand: "AMD Ryzen 9", Cores: 16, Threads: 32},
		OS:  model.OSInfo{Distro: "Unraid", Release: "6.12.6", Platform: "linux", Uptime: 1234567},
	}
	m.metrics = &model.SystemMetrics{CPUUsage: 42.5, MemoryPct: 50.0, MemoryUsed: 34359738368, MemoryTotal: 68719476736}

	view := m.View()
	if !strings.Contains(view, "AMD Ryzen 9") {
		t.Error("expected CPU brand in view")
	}
	if !strings.Contains(view, "CPU") {
		t.Error("expected CPU panel title in view")
	}
	if !strings.Contains(view, "Memory") {
		t.Error("expected Memory panel title in view")
	}
	if !strings.Contains(view, "Unraid") {
		t.Error("expected distro in view")
	}
}

func TestDashboard_View_WithError(t *testing.T) {
	m := newTestDashboard()
	m.width = 80
	m.height = 24
	m.loading = false
	m.err = errors.New("server unreachable")
	m.systemInfo = &model.SystemInfo{CPU: model.CPUInfo{Brand: "CPU"}}

	view := m.View()
	if !strings.Contains(view, "server unreachable") {
		t.Error("expected error message in view")
	}
}
