# Contributing to opencode-configurator

Thank you for your interest in contributing! This guide will help you get started.

## Prerequisites

- **Go 1.26+** — [download](https://go.dev/dl/)
- **Git** — for version control
- A terminal that supports 256 colors (for TUI rendering)

## Getting Started

```bash
# Clone the repository
git clone https://github.com/jarlex/opencode-configurator.git
cd opencode-configurator

# Build
go build -o opencode-configurator .

# Run
./opencode-configurator

# Or run directly
go run main.go
```

## Project Structure

```
main.go                        Entry point: flag parsing, initialization, program launch
internal/
  model/types.go               Domain types: Agent, Skill, MCPServer, Provider, AppState
  config/parser.go             Parses opencode.json into domain types
  config/skills.go             Scans skills directory for SKILL.md frontmatter
  config/parser_test.go        Unit tests for config parsing
  api/client.go                HTTP client for OpenCode API enrichment
  merge/merge.go               Merges static config with live API data
  tui/
    app.go                     Root Bubbletea model: dispatches keys, manages layout
    tabbar.go                  Tab navigation component (4 tabs)
    listview.go                Filterable list panel (left side)
    detail.go                  Scrollable detail panel (right side)
    statusbar.go               Status bar: online/offline, counts, key hints
    help.go                    Help overlay with keybinding reference
    keys.go                    Key binding definitions
    styles.go                  Lipgloss style constants
    render.go                  Detail rendering functions per entity type
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./internal/config/...
```

## Code Style

This project follows standard Go conventions:

```bash
# Format code
go fmt ./...

# Run static analysis
go vet ./...
```

Please ensure your code passes both `go fmt` and `go vet` before submitting.

## How to Submit Changes

1. **Fork** the repository on GitHub
2. **Clone** your fork locally
3. **Create a branch** for your change:
   ```bash
   git checkout -b feat/my-feature
   ```
4. **Make your changes** — follow existing code patterns and conventions
5. **Test** your changes:
   ```bash
   go test ./...
   go vet ./...
   ```
6. **Commit** with a descriptive message following the project convention:
   ```bash
   git commit -m "feat: add support for custom themes"
   ```
   Commit prefixes: `feat:`, `fix:`, `refactor:`, `test:`, `docs:`, `chore:`
7. **Push** your branch and open a **Pull Request**

## Reporting Issues

If you find a bug or have a feature request, please open an issue on GitHub with:

- A clear description of the problem or suggestion
- Steps to reproduce (for bugs)
- Your Go version and operating system
