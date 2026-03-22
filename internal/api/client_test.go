package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSystemInfo_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("expected Bearer test-key, got %s", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("expected application/json, got %s", got)
		}

		var req graphqlRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Query == "" {
			t.Error("expected non-empty query")
		}

		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"info": map[string]any{
					"cpu": map[string]any{"cores": 16, "threads": 32, "brand": "Ryzen 9 5900X", "manufacturer": "AMD"},
					"os":  map[string]any{"platform": "linux", "distro": "Unraid", "release": "6.12.6", "uptime": "2026-03-21T10:00:00.000Z", "hostname": "tower", "kernel": "6.1.64"},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	info, err := client.GetSystemInfo(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.CPU.Brand != "Ryzen 9 5900X" {
		t.Errorf("expected CPU brand Ryzen 9 5900X, got %s", info.CPU.Brand)
	}
	if info.CPU.Cores != 16 {
		t.Errorf("expected 16 cores, got %d", info.CPU.Cores)
	}
	if info.CPU.Manufacturer != "AMD" {
		t.Errorf("expected manufacturer AMD, got %s", info.CPU.Manufacturer)
	}
	if info.OS.Distro != "Unraid" {
		t.Errorf("expected distro Unraid, got %s", info.OS.Distro)
	}
}

func TestGetSystemMetrics_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"metrics": map[string]any{
					"cpu":    map[string]any{"percentTotal": 42.5},
					"memory": map[string]any{"used": 17179869184, "available": 55000000000, "total": 68719476736, "percentTotal": 25.0},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	metrics, err := client.GetSystemMetrics(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if metrics.CPUUsage != 42.5 {
		t.Errorf("expected CPU usage 42.5, got %f", metrics.CPUUsage)
	}
	if metrics.MemoryPct != 25.0 {
		t.Errorf("expected memory pct 25.0, got %f", metrics.MemoryPct)
	}
	if metrics.MemoryTotal != 68719476736 {
		t.Errorf("expected memory total 68719476736, got %d", metrics.MemoryTotal)
	}
}

func TestGetContainers_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"docker": map[string]any{
					"containers": []any{
						map[string]any{
							"id": "abc123", "names": []any{"plex"}, "image": "plexinc/pms:latest",
							"state": "running", "status": "Up 14 days", "autoStart": true,
							"ports": []any{
								map[string]any{"privatePort": 32400, "publicPort": 32400, "type": "tcp"},
							},
						},
						map[string]any{
							"id": "def456", "names": []any{"pihole"}, "image": "pihole/pihole:latest",
							"state": "exited", "status": "Exited (0) 2 days ago", "autoStart": false,
							"ports": []any{},
						},
					},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	containers, err := client.GetContainers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(containers) != 2 {
		t.Fatalf("expected 2 containers, got %d", len(containers))
	}

	c := containers[0]
	if c.Name != "plex" {
		t.Errorf("expected name plex, got %s", c.Name)
	}
	if c.State != "running" {
		t.Errorf("expected state running, got %s", c.State)
	}
	if len(c.Ports) != 1 || c.Ports[0].ContainerPort != 32400 {
		t.Errorf("expected port 32400, got %+v", c.Ports)
	}

	c2 := containers[1]
	if c2.Name != "pihole" {
		t.Errorf("expected name pihole, got %s", c2.Name)
	}
	if c2.State != "exited" {
		t.Errorf("expected state exited, got %s", c2.State)
	}
}

func TestGetSystemInfo_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient(server.URL, "bad-key")
	_, err := client.GetSystemInfo(context.Background())
	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestGetSystemInfo_GraphQLError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data":   nil,
			"errors": []any{map[string]any{"message": "field not found"}},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	_, err := client.GetSystemInfo(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "graphql: field not found" {
		t.Errorf("expected graphql error message, got %s", err.Error())
	}
}

func TestGetSystemInfo_ConnectionError(t *testing.T) {
	client := NewClient("http://localhost:1", "test-key")
	_, err := client.GetSystemInfo(context.Background())
	if !errors.Is(err, ErrConnectionFailed) {
		t.Errorf("expected ErrConnectionFailed, got %v", err)
	}
}
