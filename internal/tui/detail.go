package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// DetailViewModel wraps a Bubbles viewport for scrollable detail content.
type DetailViewModel struct {
	viewport viewport.Model
	content  string
	ready    bool
	width    int
	height   int
}

// NewDetailView creates a new DetailViewModel.
func NewDetailView() DetailViewModel {
	return DetailViewModel{}
}

// SetSize updates the detail viewport dimensions.
func (d *DetailViewModel) SetSize(w, h int) {
	d.width = w
	d.height = h

	// The content area is now the full inner space of the panel
	// (DetailPanel border and padding are handled externally in app.go)
	innerW := w
	if innerW < 1 {
		innerW = 1
	}

	if d.ready {
		d.viewport.Width = innerW
		d.viewport.Height = h
	} else {
		d.viewport = viewport.New(innerW, h)
		d.viewport.MouseWheelEnabled = true
		// Extend page down/up to also respond to f/b keys
		d.viewport.KeyMap.PageDown = key.NewBinding(key.WithKeys("pgdown", "f"))
		d.viewport.KeyMap.PageUp = key.NewBinding(key.WithKeys("pgup", "b"))
		d.viewport.SetContent(d.content)
		d.ready = true
	}
}

// SetContent updates the displayed content and resets scroll position.
func (d *DetailViewModel) SetContent(content string) {
	d.content = content
	if d.ready {
		d.viewport.SetContent(content)
		d.viewport.GotoTop()
	}
}

// Update handles viewport scroll events.
func (d DetailViewModel) Update(msg tea.Msg) (DetailViewModel, tea.Cmd) {
	if !d.ready {
		return d, nil
	}
	var cmd tea.Cmd
	d.viewport, cmd = d.viewport.Update(msg)
	return d, cmd
}

// View renders the scrollable detail content.
func (d DetailViewModel) View() string {
	if !d.ready {
		return ""
	}
	return d.viewport.View()
}

// Content returns the current detail text.
func (d DetailViewModel) Content() string {
	return d.content
}

// ScrollPercent returns the viewport scroll percentage.
func (d DetailViewModel) ScrollPercent() float64 {
	return d.viewport.ScrollPercent()
}
