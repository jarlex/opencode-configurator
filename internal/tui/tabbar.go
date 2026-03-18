package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TabNames defines the label for each tab.
var TabNames = []string{"Agents", "Skills", "MCP", "Providers"}

// TabChangedMsg is emitted when the active tab changes.
type TabChangedMsg struct {
	ActiveTab int
}

// TabBarModel tracks tab selection and renders the tab bar.
type TabBarModel struct {
	ActiveTab int
	width     int
	keys      KeyMap
}

// NewTabBar creates a new TabBarModel.
func NewTabBar(keys KeyMap) TabBarModel {
	return TabBarModel{
		ActiveTab: 0,
		keys:      keys,
	}
}

// Init satisfies tea.Model.
func (t TabBarModel) Init() tea.Cmd {
	return nil
}

// Update handles tab switching via Tab/ShiftTab.
func (t TabBarModel) Update(msg tea.Msg) (TabBarModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, t.keys.Tab):
			t.ActiveTab = (t.ActiveTab + 1) % len(TabNames)
			return t, func() tea.Msg { return TabChangedMsg{ActiveTab: t.ActiveTab} }
		case key.Matches(msg, t.keys.ShiftTab):
			t.ActiveTab = (t.ActiveTab - 1 + len(TabNames)) % len(TabNames)
			return t, func() tea.Msg { return TabChangedMsg{ActiveTab: t.ActiveTab} }
		}
	}
	return t, nil
}

// SetWidth updates the tab bar's available width.
func (t *TabBarModel) SetWidth(w int) {
	t.width = w
}

// View renders the tab bar.
func (t TabBarModel) View() string {
	var tabs []string
	for i, name := range TabNames {
		if i == t.ActiveTab {
			tabs = append(tabs, TabActive.Render(name))
		} else {
			tabs = append(tabs, TabInactive.Render(name))
		}
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	return TabBar.Width(t.width).Render(row + strings.Repeat(" ", max(0, t.width-lipgloss.Width(row))))
}
