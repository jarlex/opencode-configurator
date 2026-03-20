package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all keybindings for the application.
type KeyMap struct {
	// Navigation
	Tab      key.Binding
	ShiftTab key.Binding
	Up       key.Binding
	Down     key.Binding

	// Detail scroll
	PageDown   key.Binding
	PageUp     key.Binding
	HalfPageUp key.Binding
	HalfPageDn key.Binding

	// Actions
	Filter       key.Binding
	Refresh      key.Binding
	Quit         key.Binding
	Help         key.Binding
	ToggleHidden key.Binding
	CopyDetail   key.Binding
	FullScreen   key.Binding

	ForceQuit key.Binding
	Escape    key.Binding
}

// DefaultKeyMap returns the default keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next tab"),
		),
		ShiftTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev tab"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up/down item"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "up/down item"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", "f"),
			key.WithHelp("pgdn/f", "detail page down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "b"),
			key.WithHelp("pgup/b", "detail page up"),
		),
		HalfPageUp: key.NewBinding(
			key.WithKeys("ctrl+u"),
			key.WithHelp("ctrl+u", "detail half-page up"),
		),
		HalfPageDn: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "detail half-page down"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		ToggleHidden: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "toggle hidden"),
		),
		CopyDetail: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "copy detail"),
		),
		FullScreen: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "fullscreen"),
		),

		ForceQuit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "force quit"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back/close"),
		),
	}
}
