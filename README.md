# opencode-configurator

A terminal dashboard for inspecting and exploring your [OpenCode](https://opencode.ai) configuration — agents, skills, MCP servers, and providers — all from a single TUI.

```
┌─────────────────────────────────────────────────────────────────────┐
│  Agents │ Skills │ MCP │ Providers                                  │
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

- **4 navigable tabs** — Agents, Skills, MCP Servers, and Providers
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

## License

[MIT](LICENSE)
