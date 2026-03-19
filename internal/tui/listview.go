package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jarlex/opencode-configurator/internal/model"
)

// SelectionChangedMsg is emitted when the selected list item changes.
type SelectionChangedMsg struct {
	Index int
}

// EntityItem wraps any domain entity to implement bubbles list.Item.
type EntityItem struct {
	title       string
	description string
	entity      interface{} // *model.Agent, *model.Skill, *model.MCPServer, *model.Provider
}

// Ensure EntityItem implements list.Item and list.DefaultItem.
var _ list.Item = EntityItem{}

func (e EntityItem) Title() string       { return e.title }
func (e EntityItem) Description() string { return e.description }
func (e EntityItem) FilterValue() string { return e.title }

// Entity returns the underlying domain entity.
func (e EntityItem) Entity() interface{} { return e.entity }

// entityItemDelegate renders items in the list.
type entityItemDelegate struct{}

func (d entityItemDelegate) Height() int                             { return 1 }
func (d entityItemDelegate) Spacing() int                            { return 0 }
func (d entityItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d entityItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(EntityItem)
	if !ok {
		return
	}

	var style lipgloss.Style
	if index == m.Index() {
		style = ListItemSelected
	} else {
		style = ListItem
	}

	line := i.title
	if i.description != "" {
		line += " " + DetailMuted.Render("("+i.description+")")
	}
	fmt.Fprint(w, style.Render(line))
}

// ListViewModel wraps a Bubbles list.Model for entity navigation.
type ListViewModel struct {
	list    list.Model
	keys    KeyMap
	width   int
	height  int
	lastIdx int
}

// NewListView creates a new ListViewModel.
func NewListView(keys KeyMap) ListViewModel {
	delegate := entityItemDelegate{}
	l := list.New(nil, delegate, 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.DisableQuitKeybindings()

	// Override list keybindings to use our KeyMap
	l.KeyMap.CursorUp = keys.Up
	l.KeyMap.CursorDown = keys.Down
	l.KeyMap.Filter = keys.Filter
	l.KeyMap.CancelWhileFiltering = keys.Escape

	return ListViewModel{
		list:    l,
		keys:    keys,
		lastIdx: -1,
	}
}

// SetSize updates list dimensions.
func (lv *ListViewModel) SetSize(w, h int) {
	lv.width = w
	lv.height = h
	lv.list.SetSize(w, h)
}

// SetItems replaces the list items with domain entities and selects the first.
func (lv *ListViewModel) SetItems(items []list.Item) {
	lv.list.SetItems(items)
	if len(items) > 0 {
		lv.list.Select(0)
	}
	lv.lastIdx = -1 // Force re-emit of selection
}

// SelectedItem returns the currently selected EntityItem, or nil.
func (lv *ListViewModel) SelectedItem() *EntityItem {
	item := lv.list.SelectedItem()
	if item == nil {
		return nil
	}
	ei, ok := item.(EntityItem)
	if !ok {
		return nil
	}
	return &ei
}

// Items returns the current list items.
func (lv *ListViewModel) Items() []list.Item {
	return lv.list.Items()
}

// Filtering returns true if the list is in filtering mode.
func (lv *ListViewModel) Filtering() bool {
	return lv.list.FilterState() == list.Filtering
}

// FilterValue returns the current filter input text.
func (lv *ListViewModel) FilterValue() string {
	return lv.list.FilterValue()
}

// Update handles list navigation and filter input.
func (lv ListViewModel) Update(msg tea.Msg) (ListViewModel, tea.Cmd) {
	var cmd tea.Cmd
	lv.list, cmd = lv.list.Update(msg)

	// Check if selection changed
	var cmds []tea.Cmd
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	idx := lv.list.Index()
	if idx != lv.lastIdx {
		lv.lastIdx = idx
		cmds = append(cmds, func() tea.Msg {
			return SelectionChangedMsg{Index: idx}
		})
	}

	return lv, tea.Batch(cmds...)
}

// View renders the list.
func (lv ListViewModel) View() string {
	return lv.list.View()
}

// --- Helper functions to convert domain types to EntityItems ---

// AgentsToItems converts a slice of Agents to list items.
func AgentsToItems(agents []model.Agent) []list.Item {
	items := make([]list.Item, len(agents))
	for i := range agents {
		items[i] = EntityItem{
			title:       agents[i].Name,
			description: agents[i].Mode,
			entity:      &agents[i],
		}
	}
	return items
}

// SkillsToItems converts a slice of Skills to list items.
func SkillsToItems(skills []model.Skill) []list.Item {
	items := make([]list.Item, len(skills))
	for i := range skills {
		items[i] = EntityItem{
			title:       skills[i].Name,
			description: "", // detail shown in right panel
			entity:      &skills[i],
		}
	}
	return items
}

// MCPsToItems converts a slice of MCPServers to list items.
func MCPsToItems(mcps []model.MCPServer) []list.Item {
	items := make([]list.Item, len(mcps))
	for i := range mcps {
		desc := mcps[i].Type
		if mcps[i].Status != "" {
			desc += " · " + mcps[i].Status
		}
		items[i] = EntityItem{
			title:       mcps[i].Name,
			description: desc,
			entity:      &mcps[i],
		}
	}
	return items
}

// ProvidersToItems converts a slice of Providers to list items.
func ProvidersToItems(providers []model.Provider) []list.Item {
	items := make([]list.Item, len(providers))
	for i := range providers {
		desc := fmt.Sprintf("%d models", len(providers[i].Models))
		if providers[i].NPM != "" {
			desc += " · " + providers[i].NPM
		}
		items[i] = EntityItem{
			title:       providers[i].Name,
			description: desc,
			entity:      &providers[i],
		}
	}
	return items
}
