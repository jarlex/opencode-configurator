package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// StatusBarModel renders the bottom status bar.
type StatusBarModel struct {
	Online     bool
	APIError   string
	Refreshing bool
	FilterText string
	Warning    string
	width      int
}

// NewStatusBar creates a new StatusBarModel.
func NewStatusBar() StatusBarModel {
	return StatusBarModel{}
}

// SetWidth updates the status bar's available width.
func (s *StatusBarModel) SetWidth(w int) {
	s.width = w
}

// SetOnline updates the online status.
func (s *StatusBarModel) SetOnline(online bool) {
	s.Online = online
}

// SetAPIError updates the API error message.
func (s *StatusBarModel) SetAPIError(err string) {
	s.APIError = err
}

// SetRefreshing updates the refreshing state.
func (s *StatusBarModel) SetRefreshing(r bool) {
	s.Refreshing = r
}

// SetFilter updates the active filter text.
func (s *StatusBarModel) SetFilter(text string) {
	s.FilterText = text
}

// SetWarning sets a non-critical warning message (e.g., skills dir missing).
func (s *StatusBarModel) SetWarning(text string) {
	s.Warning = text
}

// View renders the status bar.
func (s StatusBarModel) View() string {
	// Left side: connection status
	var status string
	if s.Refreshing {
		status = StatusRefreshing.Render("\u27f3 Refreshing...")
	} else if s.Online {
		status = StatusOnline.Render("\u25cf Online")
	} else if s.APIError != "" && containsTimeout(s.APIError) {
		status = StatusOffline.Render("\u25cb Offline (timeout)")
	} else {
		status = StatusOffline.Render("\u25cb Offline")
	}

	// Add API error if present
	if s.APIError != "" && !s.Refreshing {
		status += "  " + DetailMuted.Render(s.APIError)
	} else if s.Warning != "" && !s.Refreshing {
		status += "  " + StatusFilter.Render(s.Warning)
	}

	// Center: filter indicator
	var filter string
	if s.FilterText != "" {
		filter = StatusFilter.Render("Filter: " + s.FilterText)
	}

	// Right side: key hints
	hints := StatusHint.Render("tab:switch  /:filter  pgdn/pgup:scroll  r:refresh  ?:help  q:quit")

	// Layout: status (left) | filter (center) | hints (right)
	leftWidth := lipgloss.Width(status)
	rightWidth := lipgloss.Width(hints)
	filterWidth := lipgloss.Width(filter)

	gap := s.width - leftWidth - rightWidth - filterWidth
	if gap < 0 {
		gap = 0
	}

	var bar string
	if filter != "" {
		leftGap := gap / 2
		rightGap := gap - leftGap
		bar = status +
			lipgloss.NewStyle().Width(leftGap).Render("") +
			filter +
			lipgloss.NewStyle().Width(rightGap).Render("") +
			hints
	} else {
		bar = status +
			lipgloss.NewStyle().Width(gap).Render("") +
			hints
	}

	return StatusBar.Width(s.width).Render(bar)
}

// containsTimeout checks if an error string indicates a timeout.
func containsTimeout(err string) bool {
	return strings.Contains(strings.ToLower(err), "timeout") ||
		strings.Contains(strings.ToLower(err), "deadline exceeded") ||
		strings.Contains(strings.ToLower(err), "not reachable")
}
