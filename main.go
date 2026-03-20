package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/jarlex/opencode-configurator/internal/api"
	"github.com/jarlex/opencode-configurator/internal/config"
	"github.com/jarlex/opencode-configurator/internal/model"
	"github.com/jarlex/opencode-configurator/internal/tui"
)

func main() {
	splitRatio := flag.Int("split", 30, "split ratio for list view vs detail view")
	// Default paths
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot determine home directory: %v\n", err)
		os.Exit(1)
	}

	defaultConfig := filepath.Join(homeDir, ".config", "opencode", "opencode.json")
	defaultSkillsDir := filepath.Join(homeDir, ".config", "opencode", "skills")

	// Flag parsing (FR-11)
	urlFlag := flag.String("url", api.DefaultBaseURL, "OpenCode API base URL")
	configFlag := flag.String("config", defaultConfig, "Path to opencode.json config file")
	flag.Parse()

	// Parse config — exit 1 if missing (SC-3)
	state, err := config.Parse(*configFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Scan skills — warn but don't fail if directory missing (SC-9)
	var skillsWarning string
	skills, err := config.ScanSkills(defaultSkillsDir)
	if err != nil {
		// Non-fatal: show warning in status bar later
		state.Skills = []model.Skill{}
		skillsWarning = fmt.Sprintf("Skills scan: %v", err)
	} else {
		state.Skills = skills
	}

	// Create API client — will be used for async enrichment (NFR-1, NFR-2)
	client := api.New(*urlFlag, api.DefaultTimeout)

	// Create the TUI app with initial offline state
	// API enrichment happens ASYNC via tea.Cmd — not blocking first render (NFR-1)
	app := tui.NewApp(state, client, *splitRatio, *configFlag)

	// Pass skills warning through the proper Warning channel (not APIError)
	if skillsWarning != "" {
		app.SetSkillsWarning(skillsWarning)
	}

	// Launch Bubbletea program
	// WithMouseCellMotion captures mouse scroll events so the terminal
	// doesn't scroll the entire alt-screen buffer (fixes frame jitter).
	p := tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
