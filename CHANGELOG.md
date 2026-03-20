# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [0.2.0] - 2026-03-20

### Added

#### TUI Features
- Number keybindings (`1`-`4`) for direct tab switching
- Item counts dynamically displayed in tab labels (e.g., "Agents (10)")
- Alphabetical list sorting for all items
- Hidden agent toggle via the `h` key
- Copy detail text to clipboard via the `y` key
- Configurable list/detail split ratio via the `--split` flag
- Full-screen detail view toggle via the `Enter` key
- Scroll position indicator added to the status bar
- Config file path added to the status bar

#### Engineering & CI
- Unit tests for `api` and `merge` packages
- GitHub Actions CI pipeline
- golangci-lint integration

## [0.1.0] - 2026-03-19

### Added

#### Documentation
- Configuration Formats reference in README documenting agent, skill, MCP, and provider JSON/YAML structures
- Field-level documentation with types, required/optional status, and descriptions
- Skill SKILL.md frontmatter format and directory layout guide

#### Configuration Parsing
- Parse `opencode.json` into typed domain models (agents, MCP servers, providers)
- Scan `~/.config/opencode/skills/*/SKILL.md` for skill metadata via YAML frontmatter
- Support for `--config` flag to specify a custom config file path
- Graceful handling of missing config (error + exit) and missing skills directory (warning)

#### API Integration
- HTTP client for OpenCode API live enrichment (`http://localhost:4096` by default)
- Fetch live MCP server status (connected, disabled, failed, needs_auth)
- Async enrichment via Bubbletea commands — never blocks UI rendering
- Configurable API URL via `--url` flag
- 3-second timeout with graceful fallback to offline mode

#### TUI Interface
- Four navigable tabs: Agents, Skills, MCP, Providers
- Split-pane layout: filterable list (left) + scrollable detail panel (right)
- Fuzzy filtering across all tabs with `/`
- Agent detail: name, mode, description, model, tools, permissions, prompt
- Skill detail: name, description, author, version, file path
- MCP detail: name, type, command/URL, enabled status, live connection status
- Provider detail: name, npm package, base URL, model list
- Tab navigation with `Tab` / `Shift+Tab`
- Item navigation with `j`/`k` or arrow keys
- Detail scrolling with `PgUp`/`PgDn`, `f`/`b`, `Ctrl+u`/`Ctrl+d`
- API data refresh with `r`
- Help overlay toggle with `?`
- Status bar showing online/offline state, item count, and key hints
- Responsive layout adapting to terminal resize (minimum 80x24)

#### CLI
- `--url` flag for OpenCode API base URL
- `--config` flag for custom config file path
- Single binary with zero external runtime dependencies
