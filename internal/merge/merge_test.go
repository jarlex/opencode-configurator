//nolint:testpackage
package merge

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jarlex/opencode-configurator/internal/api"
	"github.com/jarlex/opencode-configurator/internal/model"
)

func TestEnrich_APIUnreachable(t *testing.T) {
	client := api.New("http://127.0.0.1:0", 1*time.Second)

	initialState := &model.AppState{
		MCPs: []model.MCPServer{{Name: "test-server"}},
	}

	enriched := Enrich(initialState, client)

	if enriched.Online {
		t.Error("Expected Online to be false")
	}
	if enriched.APIError == "" {
		t.Error("Expected APIError to be set")
	}
	if len(enriched.MCPs) != 1 || enriched.MCPs[0].Name != "test-server" {
		t.Error("Expected original state to remain intact")
	}

	// Ensure we did a shallow copy correctly
	if &initialState.MCPs[0] == &enriched.MCPs[0] {
		t.Error("Expected copies of slices, not references")
	}
}

func TestEnrich_OverlayMCPStatus(t *testing.T) {
	mockStatus := map[string]api.MCPStatus{
		"server1": {Name: "server1", Status: "connected", Error: ""},
		"server2": {Name: "server2", Status: "failed", Error: "timeout"},
		"server3": {Name: "server3", Status: "", Error: "ignore empty"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.URL.Path == "/mcp" {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(mockStatus)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := api.New(server.URL, 1*time.Second)

	initialState := &model.AppState{
		MCPs: []model.MCPServer{
			{Name: "server1", Status: "disabled", Error: "old"},
			{Name: "server2", Status: "disabled", Error: ""},
			{Name: "server3", Status: "disabled", Error: "keep me"},
			{Name: "server-unknown", Status: "disabled", Error: ""},
		},
	}

	enriched := Enrich(initialState, client)

	if !enriched.Online {
		t.Error("Expected Online to be true")
	}
	if enriched.APIError != "" {
		t.Errorf("Expected empty APIError, got %s", enriched.APIError)
	}

	verifyMCPStatus(t, enriched.MCPs)
}

func verifyMCPStatus(t *testing.T, mcps []model.MCPServer) {
	for _, mcp := range mcps {
		switch mcp.Name {
		case "server1":
			if mcp.Status != "connected" || mcp.Error != "old" { // Error not updated because mockStatus had empty error
				t.Errorf("server1: got status=%s error=%s", mcp.Status, mcp.Error)
			}
		case "server2":
			if mcp.Status != "failed" || mcp.Error != "timeout" {
				t.Errorf("server2: got status=%s error=%s", mcp.Status, mcp.Error)
			}
		case "server3":
			if mcp.Status != "disabled" || mcp.Error != "ignore empty" { // Status not updated because mockStatus had empty status
				t.Errorf("server3: got status=%s error=%s", mcp.Status, mcp.Error)
			}
		case "server-unknown":
			if mcp.Status != "disabled" || mcp.Error != "" {
				t.Errorf("server-unknown modified: status=%s error=%s", mcp.Status, mcp.Error)
			}
		}
	}
}
