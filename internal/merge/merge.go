package merge

import (
	"github.com/jarlex/opencode-configurator/internal/api"
	"github.com/jarlex/opencode-configurator/internal/model"
)

// Enrich overlays live API data onto a static AppState.
// If the API is unreachable or returns an error, the original state
// is returned unchanged with Online=false and the error captured in APIError.
// Non-destructive: only non-zero API values overwrite static fields.
func Enrich(state *model.AppState, client *api.Client) *model.AppState {
	if state == nil {
		return &model.AppState{}
	}

	// Make a shallow copy so we don't mutate the caller's state.
	enriched := *state
	enriched.Agents = copyAgents(state.Agents)
	enriched.MCPs = copyMCPs(state.MCPs)
	enriched.Providers = copyProviders(state.Providers)
	enriched.Skills = copySkills(state.Skills)

	// Check connectivity first.
	if !client.Ping() {
		enriched.Online = false
		enriched.APIError = "OpenCode API is not reachable"
		return &enriched
	}

	enriched.Online = true
	enriched.APIError = ""

	// Overlay MCP live status.
	mcpStatus, err := client.FetchMCPStatus()
	if err != nil {
		enriched.APIError = "MCP status: " + err.Error()
		return &enriched
	}

	overlayMCPStatus(enriched.MCPs, mcpStatus)

	return &enriched
}

// overlayMCPStatus merges live MCP status onto static MCPServer entries.
// Only non-empty fields from the API overwrite static values.
func overlayMCPStatus(mcps []model.MCPServer, status map[string]api.MCPStatus) {
	for i := range mcps {
		live, ok := status[mcps[i].Name]
		if !ok {
			continue
		}
		// Only overlay non-empty values — never replace static with empty.
		if live.Status != "" {
			mcps[i].Status = live.Status
		}
		if live.Error != "" {
			mcps[i].Error = live.Error
		}
	}
}

// Copy helpers to avoid mutating the original state slices.

func copyAgents(src []model.Agent) []model.Agent {
	if src == nil {
		return nil
	}
	dst := make([]model.Agent, len(src))
	copy(dst, src)
	// Deep copy the Tools map for each agent to avoid shared references.
	for i := range dst {
		if dst[i].Tools != nil {
			tools := make(map[string]bool, len(dst[i].Tools))
			for k, v := range dst[i].Tools {
				tools[k] = v
			}
			dst[i].Tools = tools
		}
	}
	return dst
}

func copyMCPs(src []model.MCPServer) []model.MCPServer {
	if src == nil {
		return nil
	}
	dst := make([]model.MCPServer, len(src))
	copy(dst, src)
	return dst
}

func copyProviders(src []model.Provider) []model.Provider {
	if src == nil {
		return nil
	}
	dst := make([]model.Provider, len(src))
	copy(dst, src)
	return dst
}

func copySkills(src []model.Skill) []model.Skill {
	if src == nil {
		return nil
	}
	dst := make([]model.Skill, len(src))
	copy(dst, src)
	return dst
}
