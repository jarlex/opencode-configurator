# Roadmap

Overview of planned features and improvements for opencode-configurator.

> Last updated: 2026-03-19

## v0.2.0 — Quick Wins & Polish

| ID | Feature | Effort | Priority |
|----|---------|--------|----------|
| [x] A1 | Number keybindings (1-4) for direct tab switching | S | P1 |
| [x] A2 | Item counts in tab labels (e.g., "Agents (10)") | S | P1 |
| [x] A3 | Alphabetical list sorting | S | P1 |
| [x] A4 | Hidden agent toggle (`h` key) | S | P2 |
| [x] A5 | Copy detail to clipboard (`y` key) | S | P2 |
| [x] A6 | Configurable list/detail split ratio | S | P3 |
| [x] A7 | Full-screen detail view (`Enter` key) | M | P2 |
| [x] A8 | Scroll position indicator in status bar | S | P3 |
| [x] A9 | Config file path in status bar | S | P3 |
| [x] E1 | Unit tests for api and merge packages | M | P1 |
| [x] E2 | GitHub Actions CI pipeline | S | P1 |
| [x] E3 | golangci-lint integration | S | P2 |

## v0.3.0 — New Data Views

| ID | Feature | Effort | Priority |
|----|---------|--------|----------|
| B1 | Commands tab (parse `~/.config/opencode/commands/`) | M | P1 |
| B2 | Plugins tab (parse plugin configuration and .ts files) | M | P1 |
| B3 | Tool detail expansion with descriptions from API | M | P2 |
| B4 | Provider model details (deeper parsing) | S | P3 |
| B5 | Global search across all tabs | L | P2 |
| B6 | Agent relationship visualization (sub-agent dependencies) | M | P2 |
| A10 | Color themes (dark/light/custom, `--theme` flag) | M | P2 |
| A11 | Markdown syntax highlighting for prompts (glamour) | M | P3 |

## v0.4.0 — Configuration Management

| ID | Feature | Effort | Priority |
|----|---------|--------|----------|
| C1 | Toggle MCP server enabled/disabled (writes back to JSON) | M | P1 |
| C2 | Toggle agent hidden flag | M | P2 |
| C3 | Config validation and linting | L | P2 |
| C4 | Edit agent fields inline | XL | P3 |
| C5 | Create agent/skill from template | L | P3 |
| E4 | Config file watcher with auto-reload (fsnotify) | M | P2 |

## v0.5.0 — Monitoring & Live Features

| ID | Feature | Effort | Priority |
|----|---------|--------|----------|
| D1 | Auto-refresh MCP status polling | M | P2 |
| D2 | Session/conversation viewer tab | L | P2 |
| D3 | Token usage tracking | L | P3 |
| D4 | Agent execution history | XL | P3 |
| E5 | TUI visual tests (teatest) | L | P2 |
| E6 | Plugin system for custom tabs | XL | P3 |

## v1.0.0 — Full Release & Distribution

| ID | Feature | Effort | Priority |
|----|---------|--------|----------|
| F1 | goreleaser setup | M | P0 |
| F2 | Homebrew formula | M | P1 |
| F3 | AUR package | S | P2 |
| F4 | Docker image | S | P3 |
| F5 | Pre-built binaries on GitHub Releases | S | P0 |
| E7 | Comprehensive README rewrite | M | P1 |

## Legend

| Label | Meaning |
|-------|---------|
| **Effort**: S | Small — a few hours |
| **Effort**: M | Medium — 1-2 days |
| **Effort**: L | Large — 3-5 days |
| **Effort**: XL | Extra Large — 1+ week |
| **Priority**: P0 | Critical — blocks release |
| **Priority**: P1 | High — core feature |
| **Priority**: P2 | Medium — important improvement |
| **Priority**: P3 | Nice to have |

## Known Risks

- **v0.4.0**: Writing back to JSON must preserve existing formatting and comments
- **v0.5.0**: Session and token API endpoints may not be stable yet
- **Scope creep**: v0.4.0 should stay focused on simple toggles, not become a full config editor
- **Tab count growth**: More tabs changes navigation muscle memory; number keybindings (A1) mitigate this
