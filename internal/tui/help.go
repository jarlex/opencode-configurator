package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// HelpOverlayModel renders a centered help overlay showing all keybindings.
type HelpOverlayModel struct {
	keys   KeyMap
	width  int
	height int
}

// NewHelpOverlay creates a new HelpOverlayModel with the given keybindings.
func NewHelpOverlay(keys KeyMap) HelpOverlayModel {
	return HelpOverlayModel{keys: keys}
}

// SetSize updates the overlay's available dimensions.
func (h *HelpOverlayModel) SetSize(w, h2 int) {
	h.width = w
	h.height = h2
}

// View renders the help overlay on top of the given background content.
func (h HelpOverlayModel) View(background string) string {
	// Build help content
	bindings := []struct {
		key  string
		desc string
	}{
		{h.keys.Tab.Help().Key, h.keys.Tab.Help().Desc},
		{h.keys.ShiftTab.Help().Key, h.keys.ShiftTab.Help().Desc},
		{h.keys.Up.Help().Key, h.keys.Up.Help().Desc},
		{h.keys.Down.Help().Key, h.keys.Down.Help().Desc},
		{h.keys.PageDown.Help().Key, h.keys.PageDown.Help().Desc},
		{h.keys.PageUp.Help().Key, h.keys.PageUp.Help().Desc},
		{h.keys.HalfPageUp.Help().Key, h.keys.HalfPageUp.Help().Desc},
		{h.keys.HalfPageDn.Help().Key, h.keys.HalfPageDn.Help().Desc},
		{h.keys.Filter.Help().Key, h.keys.Filter.Help().Desc},
		{h.keys.Refresh.Help().Key, h.keys.Refresh.Help().Desc},
		{h.keys.Help.Help().Key, h.keys.Help.Help().Desc},
		{h.keys.Escape.Help().Key, h.keys.Escape.Help().Desc},
		{h.keys.Quit.Help().Key, h.keys.Quit.Help().Desc},
		{h.keys.ForceQuit.Help().Key, h.keys.ForceQuit.Help().Desc},
	}

	title := DetailTitle.Render("Keyboard Shortcuts") + "\n\n"

	var rows string
	for _, b := range bindings {
		k := HelpKey.Width(14).Render(b.key)
		d := HelpDesc.Render(b.desc)
		rows += k + "  " + d + "\n"
	}

	content := title + rows + "\n" + DetailMuted.Render("Press ? or Esc to close")

	// Calculate overlay dimensions
	overlayWidth := 50
	if overlayWidth > h.width-4 {
		overlayWidth = h.width - 4
	}

	overlay := HelpOverlay.
		Width(overlayWidth).
		Render(content)

	// Center the overlay on the screen
	overlayH := lipgloss.Height(overlay)
	overlayW := lipgloss.Width(overlay)

	topPad := (h.height - overlayH) / 2
	leftPad := (h.width - overlayW) / 2

	if topPad < 0 {
		topPad = 0
	}
	if leftPad < 0 {
		leftPad = 0
	}

	return lipgloss.Place(
		h.width, h.height,
		lipgloss.Center, lipgloss.Center,
		overlay,
		lipgloss.WithWhitespaceChars(" "),
	)
}
