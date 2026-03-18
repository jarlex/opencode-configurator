package model

// Agent represents an OpenCode agent configuration.
type Agent struct {
	Name        string
	Mode        string // primary | subagent | all
	Description string
	Prompt      string
	Tools       map[string]bool
	Permission  AgentPermission
	Hidden      bool
	Model       string
	Temperature float64
	TopP        float64
	MaxSteps    int
}

// AgentPermission holds permission settings for an agent.
type AgentPermission struct {
	Edit              string // allow | deny
	Bash              string
	Webfetch          string
	DoomLoop          string
	ExternalDirectory string
	Task              map[string]string // pattern → allow|deny
}

// Skill represents a scanned skill from a SKILL.md file.
type Skill struct {
	Name        string
	Description string
	Author      string
	Version     string
	Path        string // filesystem path to the SKILL.md
	Content     string // markdown content after the YAML frontmatter
}

// MCPServer represents an MCP server configuration.
type MCPServer struct {
	Name        string
	Type        string // local | remote
	Command     []string
	URL         string
	Enabled     bool
	Environment map[string]string
	Timeout     int
	Status      string // connected | disabled | failed | needs_auth (live)
	Error       string // live error message
}

// Provider represents a model provider configuration.
type Provider struct {
	Name    string
	NPM     string
	BaseURL string
	Models  []ProviderModel
}

// ProviderModel represents a model within a provider.
type ProviderModel struct {
	Name   string
	Launch bool
}

// AppState holds the complete application state from all data sources.
type AppState struct {
	Agents    []Agent
	Skills    []Skill
	MCPs      []MCPServer
	Providers []Provider
	Online    bool
	APIError  string
}
