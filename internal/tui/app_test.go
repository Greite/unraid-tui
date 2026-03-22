package tui

import (
	"context"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/Greite/unraid-tui/internal/api"
	"github.com/Greite/unraid-tui/internal/model"
	"github.com/Greite/unraid-tui/internal/tui/common"
)

func newTestModel() Model {
	mock := &api.MockClient{
		GetSystemInfoFn: func(ctx context.Context) (*model.SystemInfo, error) {
			return &model.SystemInfo{
				CPU: model.CPUInfo{Brand: "Test CPU", Cores: 8},
			}, nil
		},
		GetSystemMetricsFn: func(ctx context.Context) (*model.SystemMetrics, error) {
			return &model.SystemMetrics{CPUUsage: 50}, nil
		},
		GetContainersFn: func(ctx context.Context) ([]model.Container, error) {
			return []model.Container{{Name: "test"}}, nil
		},
	}
	return NewModel(mock)
}

func TestNewModel_DefaultsToeDashboard(t *testing.T) {
	m := newTestModel()
	if m.ActivePage() != common.PageDashboard {
		t.Errorf("expected PageDashboard, got %v", m.ActivePage())
	}
}

func TestUpdate_QuitOnQ(t *testing.T) {
	m := newTestModel()
	updated, cmd := m.Update(tea.KeyPressMsg(tea.Key{Code: 'q'}))
	_ = updated
	if cmd == nil {
		t.Fatal("expected quit command, got nil")
	}
}

func TestUpdate_QuitOnCtrlC(t *testing.T) {
	m := newTestModel()
	_, cmd := m.Update(tea.KeyPressMsg(tea.Key{Code: 'c', Mod: tea.ModCtrl}))
	if cmd == nil {
		t.Fatal("expected quit command, got nil")
	}
}

func TestUpdate_TabSwitchesToDocker(t *testing.T) {
	m := newTestModel()
	updated, _ := m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyTab}))
	root := updated.(Model)
	if root.ActivePage() != common.PageDocker {
		t.Errorf("expected PageDocker after tab, got %v", root.ActivePage())
	}
}

func TestUpdate_TabCyclesThroughPages(t *testing.T) {
	m := newTestModel()
	// Tab through all pages and back to Dashboard
	var updated tea.Model = m
	for range int(common.PageCount) {
		updated, _ = updated.(Model).Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyTab}))
	}
	root := updated.(Model)
	if root.ActivePage() != common.PageDashboard {
		t.Errorf("expected PageDashboard after full cycle, got %v", root.ActivePage())
	}
}

func TestUpdate_FKeySwitchesPage(t *testing.T) {
	m := newTestModel()
	updated, _ := m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyF3}))
	root := updated.(Model)
	if root.ActivePage() != common.PageVMs {
		t.Errorf("expected PageVMs on F3, got %v", root.ActivePage())
	}

	updated, _ = root.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyF4}))
	root = updated.(Model)
	if root.ActivePage() != common.PageNotifications {
		t.Errorf("expected PageNotifications on F4, got %v", root.ActivePage())
	}

	updated, _ = root.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyF5}))
	root = updated.(Model)
	if root.ActivePage() != common.PageShares {
		t.Errorf("expected PageShares on F5, got %v", root.ActivePage())
	}
}

func TestUpdate_F1SwitchesToDashboard(t *testing.T) {
	m := newTestModel()
	// Switch to Docker first
	m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyTab}))
	// F1 to Dashboard
	updated, _ := m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyF1}))
	root := updated.(Model)
	if root.ActivePage() != common.PageDashboard {
		t.Errorf("expected PageDashboard on F1, got %v", root.ActivePage())
	}
}

func TestUpdate_F2SwitchesToDocker(t *testing.T) {
	m := newTestModel()
	updated, _ := m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyF2}))
	root := updated.(Model)
	if root.ActivePage() != common.PageDocker {
		t.Errorf("expected PageDocker on F2, got %v", root.ActivePage())
	}
}

func TestUpdate_WindowSizeMsg(t *testing.T) {
	m := newTestModel()
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	root := updated.(Model)
	if root.width != 120 || root.height != 40 {
		t.Errorf("expected 120x40, got %dx%d", root.width, root.height)
	}
}

func TestView_ReturnsNonEmpty(t *testing.T) {
	m := newTestModel()
	m.width = 80
	m.height = 24
	v := m.View()
	if v.Content == "" {
		t.Error("expected non-empty view content")
	}
}
