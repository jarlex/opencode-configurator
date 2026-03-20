package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/jarlex/opencode-configurator/internal/model"
)

func TestAppModel_TabSwitching(t *testing.T) {
	app := NewApp(&model.AppState{}, nil, 30, "")

	// Simulate pressing "tab"
	msg := tea.KeyMsg{Type: tea.KeyTab}
	appModel, _ := app.Update(msg)
	newApp := appModel.(AppModel)

	if newApp.tabBar.ActiveTab != 1 {
		t.Errorf("expected tab 1 (Skills), got %d", newApp.tabBar.ActiveTab)
	}

	// Simulate pressing "shift+tab"
	msg = tea.KeyMsg{Type: tea.KeyShiftTab}
	appModel, _ = newApp.Update(msg)
	newApp = appModel.(AppModel)

	if newApp.tabBar.ActiveTab != 0 {
		t.Errorf("expected tab 0 (Agents), got %d", newApp.tabBar.ActiveTab)
	}

	// Simulate pressing "shift+tab" again to wrap around
	msg = tea.KeyMsg{Type: tea.KeyShiftTab}
	appModel, _ = newApp.Update(msg)
	newApp = appModel.(AppModel)

	if newApp.tabBar.ActiveTab != 3 {
		t.Errorf("expected tab 3 (Providers) due to wraparound, got %d", newApp.tabBar.ActiveTab)
	}
}

func TestAppModel_PopulateListSorting(t *testing.T) {
	state := &model.AppState{
		Agents: []model.Agent{
			{Name: "Zebra"},
			{Name: "Apple"},
			{Name: "Mango"},
		},
	}

	app := NewApp(state, nil, 30, "")
	app.populateList(0)

	items := app.listView.Items()
	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}

	if items[0].FilterValue() != "Apple" {
		t.Errorf("expected Apple, got %s", items[0].FilterValue())
	}
	if items[1].FilterValue() != "Mango" {
		t.Errorf("expected Mango, got %s", items[1].FilterValue())
	}
	if items[2].FilterValue() != "Zebra" {
		t.Errorf("expected Zebra, got %s", items[2].FilterValue())
	}
}
