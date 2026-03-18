package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DefaultBaseURL is the default OpenCode API URL.
const DefaultBaseURL = "http://localhost:4096"

// DefaultTimeout is the default timeout for data requests.
const DefaultTimeout = 3 * time.Second

// PingTimeout is the timeout for the Ping health check.
const PingTimeout = 1 * time.Second

// Client communicates with the OpenCode HTTP API.
type Client struct {
	BaseURL string
	HTTP    *http.Client
}

// New creates a new API client with the given base URL and timeout.
func New(baseURL string, timeout time.Duration) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTP: &http.Client{
			Timeout: timeout,
		},
	}
}

// MCPStatus represents the live status of an MCP server from the API.
type MCPStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"` // connected | disabled | failed | needs_auth
	Error  string `json:"error"`
}

// Ping checks if the OpenCode API is reachable.
// Uses a shorter timeout (1s) to avoid blocking startup.
func (c *Client) Ping() bool {
	pingClient := &http.Client{Timeout: PingTimeout}
	req, err := http.NewRequest(http.MethodHead, c.BaseURL+"/config", nil)
	if err != nil {
		return false
	}
	resp, err := pingClient.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

// FetchMCPStatus fetches live MCP server statuses from GET /mcp.
// Returns a map of server name to MCPStatus.
func (c *Client) FetchMCPStatus() (map[string]MCPStatus, error) {
	body, err := c.get("/mcp")
	if err != nil {
		return nil, fmt.Errorf("fetching MCP status: %w", err)
	}

	var result map[string]MCPStatus
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing MCP status response: %w", err)
	}
	return result, nil
}

// get performs a GET request to the given path and returns the response body.
func (c *Client) get(path string) ([]byte, error) {
	resp, err := c.HTTP.Get(c.BaseURL + path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, path)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	return body, nil
}
