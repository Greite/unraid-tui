package docker

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Greite/unraid-tui/internal/api"
	"github.com/Greite/unraid-tui/internal/model"
	"github.com/Greite/unraid-tui/internal/tui/common"
)

func newTestDocker() Model {
	mock := &api.MockClient{
		GetContainersFn: func(ctx context.Context) ([]model.Container, error) {
			return []model.Container{
				{Name: "plex", Image: "plexinc/pms:latest", State: "running", Status: "Up 14 days",
					Ports: []model.Port{{HostPort: 32400, ContainerPort: 32400, Protocol: "tcp"}}},
				{Name: "pihole", Image: "pihole/pihole:latest", State: "exited", Status: "Exited (0)"},
			}, nil
		},
	}
	return New(mock)
}

func TestDocker_InitialState(t *testing.T) {
	m := newTestDocker()
	if !m.loading {
		t.Error("expected loading to be true initially")
	}
	if len(m.containers) != 0 {
		t.Error("expected no containers initially")
	}
}

func TestDocker_ContainersMsg(t *testing.T) {
	m := newTestDocker()
	m.width = 100
	m.height = 30

	containers := []model.Container{
		{Name: "nginx", Image: "nginx:latest", State: "running", Status: "Up 2 hours",
			Ports: []model.Port{{HostPort: 80, ContainerPort: 80, Protocol: "tcp"}}},
		{Name: "redis", Image: "redis:alpine", State: "running", Status: "Up 1 hour"},
		{Name: "postgres", Image: "postgres:15", State: "exited", Status: "Exited (0)"},
	}

	updated, _ := m.Update(common.ContainersMsg{Containers: containers})

	if updated.loading {
		t.Error("expected loading to be false")
	}
	if len(updated.containers) != 3 {
		t.Errorf("expected 3 containers, got %d", len(updated.containers))
	}
	if updated.err != nil {
		t.Errorf("expected no error, got %v", updated.err)
	}
}

func TestDocker_ContainersMsg_Error(t *testing.T) {
	m := newTestDocker()
	m.width = 100
	m.height = 30

	updated, _ := m.Update(common.ContainersMsg{Err: errors.New("API error")})

	if updated.loading {
		t.Error("expected loading to be false")
	}
	if updated.err == nil {
		t.Fatal("expected error to be set")
	}
}

func TestDocker_View_Loading(t *testing.T) {
	m := newTestDocker()
	m.width = 100
	m.height = 30
	view := m.View()
	if !strings.Contains(view, "Loading") {
		t.Error("expected loading message")
	}
}

func TestDocker_View_WithContainers(t *testing.T) {
	m := newTestDocker()
	m.width = 100
	m.height = 30
	m.loading = false
	m.containers = []model.Container{
		{Name: "plex", Image: "plexinc/pms:latest", State: "running", Status: "Up 14 days",
			Ports: []model.Port{{HostPort: 32400, ContainerPort: 32400, Protocol: "tcp"}}},
		{Name: "pihole", Image: "pihole/pihole:latest", State: "exited", Status: "Exited (0)"},
	}
	m.applySort()
	view := m.View()
	if !strings.Contains(view, "Containers (2)") {
		t.Error("expected container count in view")
	}
	if !strings.Contains(view, "1 running") {
		t.Error("expected running count in view")
	}
}

func TestDocker_View_WithError(t *testing.T) {
	m := newTestDocker()
	m.width = 100
	m.height = 30
	m.loading = false
	m.err = errors.New("connection refused")

	view := m.View()
	if !strings.Contains(view, "connection refused") {
		t.Error("expected error message in view")
	}
}

func TestFormatPorts(t *testing.T) {
	tests := []struct {
		name  string
		ports []model.Port
		want  string
	}{
		{"empty", nil, "-"},
		{"single", []model.Port{{HostPort: 80, ContainerPort: 80}}, "80:80"},
		{"no host port", []model.Port{{ContainerPort: 3306}}, "3306"},
		{"multiple", []model.Port{
			{HostPort: 80, ContainerPort: 80},
			{HostPort: 443, ContainerPort: 443},
		}, "80:80, 443:443"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatPorts(tt.ports)
			if got != tt.want {
				t.Errorf("formatPorts() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input string
		max   int
		want  string
	}{
		{"short", 10, "short"},
		{"a-very-long-image-name", 10, "a-very-l.."},
		{"exact-len", 9, "exact-len"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := truncate(tt.input, tt.max)
			if got != tt.want {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.max, got, tt.want)
			}
		})
	}
}

func TestSort_ByName(t *testing.T) {
	m := newTestDocker()
	m.width = 100
	m.height = 30
	m.loading = false
	m.containers = []model.Container{
		{Name: "zebra", State: "running"},
		{Name: "alpha", State: "exited"},
		{Name: "middle", State: "running"},
	}
	m.sortCol = sortName
	m.sortAsc = true
	m.applySort()

	if m.sorted[0].Name != "alpha" {
		t.Errorf("expected first sorted container to be alpha, got %s", m.sorted[0].Name)
	}
	if m.sorted[2].Name != "zebra" {
		t.Errorf("expected last sorted container to be zebra, got %s", m.sorted[2].Name)
	}
}

func TestSort_ByNameDesc(t *testing.T) {
	m := newTestDocker()
	m.containers = []model.Container{
		{Name: "alpha"},
		{Name: "zebra"},
	}
	m.sortCol = sortName
	m.sortAsc = false
	m.applySort()

	if m.sorted[0].Name != "zebra" {
		t.Errorf("expected first sorted container to be zebra, got %s", m.sorted[0].Name)
	}
}

func TestSort_ByState(t *testing.T) {
	m := newTestDocker()
	m.containers = []model.Container{
		{Name: "a", State: "running"},
		{Name: "b", State: "exited"},
		{Name: "c", State: "running"},
	}
	m.sortCol = sortState
	m.sortAsc = true
	m.applySort()

	if m.sorted[0].State != "exited" {
		t.Errorf("expected exited first, got %s", m.sorted[0].State)
	}
}

func TestToggleSort_SameColumnReversesOrder(t *testing.T) {
	m := newTestDocker()
	m.containers = []model.Container{{Name: "a"}, {Name: "b"}}
	m.sortCol = sortName
	m.sortAsc = true

	m.toggleSort(sortName)
	if m.sortAsc {
		t.Error("expected sortAsc to be false after toggle")
	}

	m.toggleSort(sortName)
	if !m.sortAsc {
		t.Error("expected sortAsc to be true after second toggle")
	}
}

func TestToggleSort_DifferentColumnResetsToAsc(t *testing.T) {
	m := newTestDocker()
	m.containers = []model.Container{{Name: "a"}}
	m.sortCol = sortName
	m.sortAsc = false

	m.toggleSort(sortState)
	if m.sortCol != sortState {
		t.Errorf("expected sortCol to be sortState, got %d", m.sortCol)
	}
	if !m.sortAsc {
		t.Error("expected sortAsc to be true when switching columns")
	}
}

func TestStateStyle(t *testing.T) {
	if got := stateIcon("running", false); !strings.Contains(got, "running") {
		t.Errorf("expected 'running' in output, got %q", got)
	}
	if got := stateIcon("exited", false); !strings.Contains(got, "exited") {
		t.Errorf("expected 'exited' in output, got %q", got)
	}
	if got := stateIcon("running", true); !strings.Contains(got, "⬆") {
		t.Errorf("expected update indicator in output, got %q", got)
	}
	if got := stateIcon("unknown", false); got != "unknown" {
		t.Errorf("expected 'unknown', got %q", got)
	}
}
