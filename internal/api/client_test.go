//nolint:testpackage
package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_Ping(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{"success 200", http.StatusOK, true},
		{"success 204", http.StatusNoContent, true},
		{"failure 404", http.StatusNotFound, false},
		{"failure 500", http.StatusInternalServerError, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = r
				if r.URL.Path != "/config" {
					t.Errorf("expected path /config, got %s", r.URL.Path)
				}
				if r.Method != http.MethodHead {
					t.Errorf("expected HEAD method, got %s", r.Method)
				}
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := New(server.URL, 1*time.Second)
			got := client.Ping()
			if got != tt.want {
				t.Errorf("Ping() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Ping_Unreachable(t *testing.T) {
	client := New("http://127.0.0.1:0", 1*time.Second) // Invalid port
	if got := client.Ping(); got != false {
		t.Errorf("Ping() = %v, want false", got)
	}
}

func TestClient_FetchMCPStatus(t *testing.T) {
	mockStatus := map[string]MCPStatus{
		"server1": {Name: "server1", Status: "connected", Error: ""},
		"server2": {Name: "server2", Status: "failed", Error: "timeout"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r
		if r.URL.Path != "/mcp" {
			t.Errorf("expected path /mcp, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockStatus)
	}))
	defer server.Close()

	client := New(server.URL, 1*time.Second)
	got, err := client.FetchMCPStatus()
	if err != nil {
		t.Fatalf("FetchMCPStatus() error = %v", err)
	}

	if len(got) != len(mockStatus) {
		t.Errorf("got %d statuses, want %d", len(got), len(mockStatus))
	}

	for k, v := range mockStatus {
		if got[k] != v {
			t.Errorf("status[%s] = %v, want %v", k, got[k], v)
		}
	}
}

func TestClient_FetchMCPStatus_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := New(server.URL, 1*time.Second)
	_, err := client.FetchMCPStatus()
	if err == nil {
		t.Error("FetchMCPStatus() expected error, got nil")
	}
}
