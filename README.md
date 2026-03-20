# opencode-configurator

A terminal dashboard for inspecting and exploring your [OpenCode](https://opencode.ai) configuration — agents, skills, MCP servers, and providers — all from a single TUI.

```
┌─────────────────────────────────────────────────────────────────────┐
│  Agents (6) │ Skills (12) │ MCP (3) │ Providers (2)                 │
├───────────────┬─────────────────────────────────────────────────────┤
│ ▸ orchestrator│  Name:        orchestrator                         │
│   sdd-apply   │  Mode:        primary                              │
│   sdd-spec    │  Description: Coordinator agent that delegates...  │
│   sdd-design  │  Model:       claude-sonnet-4-20250514             │
│   sdd-verify  │  Tools:       engram_*, delegate, task, skill      │
│   sdd-tasks   │  Permissions: edit=allow bash=allow web=allow      │
│               │  Prompt:      ## Agent Teams Orchestrator [...]    │
│               │                                                     │
│               │                                                     │
├───────────────┴─────────────────────────────────────────────────────┤
│  Online │ 6 agents │ Tab: switch │ /: filter │ r: refresh │ ?: help│
└─────────────────────────────────────────────────────────────────────┘
```

## Features

- **4 navigable tabs** — Agents, Skills, MCP Servers, and Providers (switch quickly with `1`-`4` keys)
- **Dynamic tab counts** — instantly see how many items are in each section (e.g., "Agents (10)")
- **Alphabetical sorting** — list items are automatically sorted A-Z for easy scanning
- **Offline-first** — launches instantly from your `opencode.json` config file; no running server required
- **Live API enrichment** — connects to the OpenCode API for real-time MCP status and tool data
- **Filterable lists** — fuzzy search across all items with `/`
- **Scrollable detail panel** — full agent prompts, tool lists, and permissions at a glance
- **Keyboard-driven** — navigate everything without touching the mouse
- **Responsive layout** — adapts to terminal resize (minimum 80x24)
- **Single binary** — zero external runtime dependencies

## Requirements

- **Go 1.26+** (for building from source)
- An `opencode.json` configuration file (typically at `~/.config/opencode/opencode.json`)

## Installation

### From source

```bash
go install github.com/jarlex/opencode-configurator@latest
```

### Manual build

```bash
git clone https://github.com/jarlex/opencode-configurator.git
cd opencode-configurator
go build -o opencode-configurator .
```

## Usage

```bash
# Launch with default config (~/.config/opencode/opencode.json)
opencode-configurator

# Connect to a running OpenCode API server
opencode-configurator --url http://localhost:4096

# Use a specific config file
opencode-configurator --config /path/to/opencode.json

# Both options combined
opencode-configurator --url http://localhost:4096 --config ./opencode.json
```

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `1`-`4` | Quick switch to tabs 1 through 4 |
| `Tab` | Next tab |
| `Shift+Tab` | Previous tab |
| `j` / `↓` | Next item in list |
| `k` / `↑` | Previous item in list |
| `/` | Filter items (fuzzy search) |
| `Esc` | Clear filter / close overlay |
| `r` | Refresh data from API |
| `?` | Toggle help overlay |
| `PgDn` / `f` | Scroll detail page down |
| `PgUp` / `b` | Scroll detail page up |
| `Ctrl+d` | Scroll detail half-page down |
| `Ctrl+u` | Scroll detail half-page up |
| `q` | Quit |
| `Ctrl+c` | Force quit |

## Configuration Formats

OpenCode's configuration is not yet fully standardized. This section documents the expected formats that `opencode-configurator` reads and displays, based on real-world usage.

### Agent Configuration Format

Agents are defined in `opencode.json` under the top-level `"agent"` key. Each agent is a named entry with the following structure:

```json
{
  "agent": {
    "my-agent": {
      "mode": "subagent",
      "hidden": true,
      "description": "Short description of what this agent does",
      "prompt": "System prompt text that defines the agent's behavior and instructions.",
      "model": "claude-sonnet-4-20250514",
      "temperature": 0.7,
      "top_p": 0.9,
      "maxSteps": 50,
      "tools": {
        "read": true,
        "write": true,
        "edit": true,
        "bash": true,
        "delegate": false
      },
      "permission": {
        "edit": "allow",
        "bash": "allow",
        "webfetch": "deny",
        "doom_loop": "deny",
        "external_directory": "deny",
        "task": {
          "*": "deny",
          "sdd-*": "allow"
        }
      }
    }
  }
}
```

#### Agent Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `mode` | `string` | Yes | Agent mode: `"primary"` (main agent), `"subagent"` (delegated worker), or `"all"` (both). |
| `description` | `string` | Yes | Short description shown in agent listings. |
| `prompt` | `string` | Yes | System prompt defining the agent's behavior and instructions. |
| `tools` | `object` | No | Map of tool names to boolean enable/disable. Controls which tools the agent can use (e.g., `read`, `write`, `edit`, `bash`, `delegate`). |
| `permission` | `object` | No | Fine-grained permissions for actions like `edit`, `bash`, `webfetch`, `doom_loop`, `external_directory`, and `task` delegation patterns. Values are `"allow"` or `"deny"`. The `task` sub-field supports glob patterns (e.g., `"sdd-*": "allow"`). |
| `hidden` | `boolean` | No | If `true`, the agent is hidden from default listings. Default: `false`. |
| `model` | `string` | No | Override the default model for this agent (e.g., `"claude-sonnet-4-20250514"`). |
| `temperature` | `number` | No | Sampling temperature (0.0–2.0). |
| `top_p` | `number` | No | Nucleus sampling parameter (0.0–1.0). |
| `maxSteps` | `integer` | No | Maximum reasoning steps before the agent stops. |

### Skill Configuration Format

Skills are Markdown files named `SKILL.md` located in subdirectories of the skills directory (typically `~/.config/opencode/skills/`). Each skill uses YAML frontmatter between `---` delimiters, followed by the skill's prompt/instructions as Markdown content.

**Example** (`~/.config/opencode/skills/my-skill/SKILL.md`):

```markdown
---
name: my-skill
description: >
  A short description of what this skill does and when it should be triggered.
  Trigger: When the user asks to perform a specific task.
license: MIT
metadata:
  author: your-name
  version: "1.0"
---

## Purpose

You are a sub-agent responsible for [specific task]. This section contains
the full prompt and instructions that define the skill's behavior.

## What to Do

### Step 1: Analyze
- Read the relevant files
- Understand the context

### Step 2: Execute
- Perform the required actions
- Follow project conventions

## Rules
- Always follow the specs
- Never skip validation
```

#### Skill Frontmatter Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | `string` | Yes | Unique identifier for the skill. Should match the directory name (e.g., `my-skill/SKILL.md` → `name: my-skill`). |
| `description` | `string` | Yes | Description of what the skill does and when to trigger it. Multi-line supported via YAML `>` syntax. |
| `license` | `string` | No | License identifier (e.g., `MIT`, `Apache-2.0`). |
| `metadata.author` | `string` | No | Author name or handle. |
| `metadata.version` | `string` | No | Semantic version string (e.g., `"1.0"`, `"2.0"`). |

#### How Skills Are Parsed

The Markdown content **after** the closing `---` delimiter is the skill's body — this is the actual prompt/instructions that an agent receives when the skill is loaded. The frontmatter is metadata only; it is not passed to the agent.

Directory structure expected:

```
~/.config/opencode/skills/
├── my-skill/
│   └── SKILL.md            # Frontmatter + instructions
├── another-skill/
│   └── SKILL.md
├── _shared/                 # Shared resources (skipped during scan)
│   ├── common-rules.md
│   └── conventions.md
└── ...
```

> **Note:** Directories prefixed with `_` (e.g., `_shared/`) are skipped during skill scanning and treated as shared resource directories.

### MCP Server Configuration Format

MCP servers are defined in `opencode.json` under the `"mcp"` key:

```json
{
  "mcp": {
    "my-server": {
      "type": "local",
      "command": ["my-binary", "mcp", "--flag=value"],
      "environment": {
        "API_KEY": "value"
      },
      "enabled": true,
      "timeout": 30
    }
  }
}
```

Local servers use `command` (array of strings); remote servers use `url` instead.

### Provider Configuration Format

Providers are defined in `opencode.json` under the `"provider"` key:

```json
{
  "provider": {
    "ollama": {
      "name": "Ollama",
      "npm": "@ai-sdk/openai-compatible",
      "options": {
        "baseURL": "http://127.0.0.1:11434/v1"
      },
      "models": {
        "my-model": {
          "name": "my-model",
          "_launch": true
        }
      }
    }
  }
}
```

## Architecture

The project is organized into 5 internal packages following the Go `internal/` convention:

```
main.go                        Entry point: flag parsing, init, launch
internal/
  model/                       Domain types (Agent, Skill, MCPServer, Provider, AppState)
  config/                      Config parser (opencode.json) and skills scanner (SKILL.md)
  api/                         HTTP client for OpenCode API enrichment
  merge/                       Merges static config data with live API responses
  tui/                         Bubbletea TUI: app, tabs, list, detail, status bar, help
```

**Data flow:** `opencode.json` + `skills/` directory are parsed into an `AppState` (offline-first). If the OpenCode API is reachable, live data (MCP status, tool IDs) is fetched asynchronously and merged into the state without blocking the UI.

## Roadmap

See [ROADMAP.md](ROADMAP.md) for planned features and improvements across upcoming versions (v0.2.0 through v1.0.0).

## License

[MIT](LICENSE)
