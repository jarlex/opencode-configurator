package tui

import "github.com/charmbracelet/lipgloss"

// Color palette — dark terminal friendly.
var (
	colorPrimary   = lipgloss.Color("#7C3AED") // violet
	colorSecondary = lipgloss.Color("#A78BFA") // light violet
	colorAccent    = lipgloss.Color("#10B981") // green
	colorMuted     = lipgloss.Color("#6B7280") // gray
	colorDanger    = lipgloss.Color("#EF4444") // red
	colorWarning   = lipgloss.Color("#F59E0B") // amber
	colorBg        = lipgloss.Color("#1F2937") // dark bg
	colorFg        = lipgloss.Color("#F9FAFB") // white-ish
	colorBorder    = lipgloss.Color("#374151") // subtle border
)

// Tab styles.
var (
	TabActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorFg).
			Background(colorPrimary).
			Padding(0, 2)

	TabInactive = lipgloss.NewStyle().
			Foreground(colorMuted).
			Padding(0, 2)

	TabBar = lipgloss.NewStyle().
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(colorBorder).
		MarginBottom(0)
)

// List styles.
var (
	ListItem = lipgloss.NewStyle().
			Foreground(colorFg).
			PaddingLeft(2)

	ListItemSelected = lipgloss.NewStyle().
				Foreground(colorFg).
				Background(colorPrimary).
				Bold(true).
				PaddingLeft(2)

	ListTitle = lipgloss.NewStyle().
			Foreground(colorSecondary).
			Bold(true).
			PaddingLeft(1)
)

// Detail pane styles.
var (
	DetailBorder = lipgloss.NewStyle().
			BorderLeft(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorBorder).
			PaddingLeft(2).
			PaddingRight(1)

	DetailTitle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true).
			MarginBottom(1)

	DetailLabel = lipgloss.NewStyle().
			Foreground(colorSecondary).
			Bold(true)

	DetailValue = lipgloss.NewStyle().
			Foreground(colorFg)

	DetailMuted = lipgloss.NewStyle().
			Foreground(colorMuted)
)

// Status bar styles.
var (
	StatusBar = lipgloss.NewStyle().
			Foreground(colorMuted).
			PaddingLeft(1).
			PaddingRight(1)

	StatusOnline = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	StatusOffline = lipgloss.NewStyle().
			Foreground(colorDanger).
			Bold(true)

	StatusRefreshing = lipgloss.NewStyle().
				Foreground(colorWarning).
				Bold(true)

	StatusHint = lipgloss.NewStyle().
			Foreground(colorMuted)

	StatusFilter = lipgloss.NewStyle().
			Foreground(colorWarning)
)

// Outer frame — wraps the entire app in a visible box.
var (
	AppFrame = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary)

	ListPanel = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(colorBorder)

	DetailPanel = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(colorBorder).
			PaddingLeft(1).
			PaddingRight(1)
)

// Help overlay styles.
var (
	HelpOverlay = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(1, 2).
			Background(colorBg)

	HelpKey = lipgloss.NewStyle().
		Foreground(colorSecondary).
		Bold(true)

	HelpDesc = lipgloss.NewStyle().
			Foreground(colorFg)
)
