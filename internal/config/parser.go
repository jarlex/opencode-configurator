package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jarlex/opencode-configurator/internal/model"
)

// configJSON mirrors the top-level opencode.json structure.
type configJSON struct {
	Agent    map[string]agentJSON    `json:"agent"`
	MCP      map[string]mcpJSON      `json:"mcp"`
	Provider map[string]providerJSON `json:"provider"`
}

type agentJSON struct {
	Mode        string          `json:"mode"`
	Description string          `json:"description"`
	Prompt      string          `json:"prompt"`
	Tools       map[string]bool `json:"tools"`
	Permission  permissionJSON  `json:"permission"`
	Hidden      bool            `json:"hidden"`
	Model       string          `json:"model"`
	Temperature float64         `json:"temperature"`
	TopP        float64         `json:"top_p"`
	MaxSteps    int             `json:"maxSteps"`
}

type permissionJSON struct {
	Edit              string            `json:"edit"`
	Bash              string            `json:"bash"`
	Webfetch          string            `json:"webfetch"`
	DoomLoop          string            `json:"doom_loop"`
	ExternalDirectory string            `json:"external_directory"`
	Task              map[string]string `json:"task"`
}

type mcpJSON struct {
	Type        string            `json:"type"`
	Command     []string          `json:"command"`
	URL         string            `json:"url"`
	Enabled     bool              `json:"enabled"`
	Environment map[string]string `json:"environment"`
	Timeout     int               `json:"timeout"`
}

type providerJSON struct {
	Name    string                       `json:"name"`
	NPM     string                       `json:"npm"`
	Models  map[string]providerModelJSON `json:"models"`
	Options providerOptionsJSON          `json:"options"`
}

type providerModelJSON struct {
	Name   string `json:"name"`
	Launch bool   `json:"_launch"`
}

type providerOptionsJSON struct {
	BaseURL string `json:"baseURL"`
}

// Parse reads an opencode.json file and returns the parsed AppState.
// Returns an error if the file does not exist or cannot be parsed.
func Parse(path string) (*model.AppState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found: %s", path)
		}
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg configJSON
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config JSON: %w", err)
	}

	state := &model.AppState{}

	// Convert agents
	for name, a := range cfg.Agent {
		agent := model.Agent{
			Name:        name,
			Mode:        a.Mode,
			Description: a.Description,
			Prompt:      a.Prompt,
			Tools:       a.Tools,
			Permission: model.AgentPermission{
				Edit:              a.Permission.Edit,
				Bash:              a.Permission.Bash,
				Webfetch:          a.Permission.Webfetch,
				DoomLoop:          a.Permission.DoomLoop,
				ExternalDirectory: a.Permission.ExternalDirectory,
				Task:              a.Permission.Task,
			},
			Hidden:      a.Hidden,
			Model:       a.Model,
			Temperature: a.Temperature,
			TopP:        a.TopP,
			MaxSteps:    a.MaxSteps,
		}
		state.Agents = append(state.Agents, agent)
	}

	// Convert MCP servers
	for name, m := range cfg.MCP {
		mcp := model.MCPServer{
			Name:        name,
			Type:        m.Type,
			Command:     m.Command,
			URL:         m.URL,
			Enabled:     m.Enabled,
			Environment: m.Environment,
			Timeout:     m.Timeout,
		}
		state.MCPs = append(state.MCPs, mcp)
	}

	// Convert providers
	for key, p := range cfg.Provider {
		provider := model.Provider{
			Name:    p.Name,
			NPM:     p.NPM,
			BaseURL: p.Options.BaseURL,
		}
		// If provider name is empty, use the map key
		if provider.Name == "" {
			provider.Name = key
		}
		for modelName, pm := range p.Models {
			name := pm.Name
			if name == "" {
				name = modelName
			}
			provider.Models = append(provider.Models, model.ProviderModel{
				Name:   name,
				Launch: pm.Launch,
			})
		}
		state.Providers = append(state.Providers, provider)
	}

	return state, nil
}
