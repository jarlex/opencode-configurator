package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jarlex/opencode-configurator/internal/api"
	"github.com/jarlex/opencode-configurator/internal/merge"
	"github.com/jarlex/opencode-configurator/internal/model"
)

// DataLoadedMsg signals that AppState data is ready (after async enrichment).
type DataLoadedMsg struct {
	State *model.AppState
}

// EnrichErrorMsg signals that API enrichment failed.
type EnrichErrorMsg struct {
	Err error
}

// AppModel is the root Bubbletea model that owns all sub-components.
type AppModel struct {
	state     *model.AppState
	tabBar    TabBarModel
	listView  ListViewModel
	detail    DetailViewModel
	statusBar StatusBarModel
	help      HelpOverlayModel
	keys      KeyMap

	apiClient *api.Client
	showHelp  bool

	width  int
	height int
	ready  bool
}

// NewApp creates the root application model with initial state.
func NewApp(state *model.AppState, client *api.Client) AppModel {
	keys := DefaultKeyMap()
	app := AppModel{
		state:     state,
		tabBar:    NewTabBar(keys),
		listView:  NewListView(keys),
		detail:    NewDetailView(),
		statusBar: NewStatusBar(),
		help:      NewHelpOverlay(keys),
		keys:      keys,
		apiClient: client,
	}
	app.statusBar.SetOnline(state.Online)
	app.statusBar.SetAPIError(state.APIError)
	return app
}

// SetSkillsWarning passes a skills directory warning to the status bar.
func (a *AppModel) SetSkillsWarning(warning string) {
	a.statusBar.SetWarning(warning)
}

// Init satisfies tea.Model. Fires initial tab load + async API enrichment.
func (a AppModel) Init() tea.Cmd {
	return tea.Batch(
		// Populate the initial tab
		func() tea.Msg {
			return TabChangedMsg{ActiveTab: 0}
		},
		// Async API enrichment — does NOT block first render (NFR-1)
		a.enrichCmd(),
	)
}

// enrichCmd returns a tea.Cmd that performs API enrichment in the background.
func (a AppModel) enrichCmd() tea.Cmd {
	client := a.apiClient
	state := a.state
	if client == nil || state == nil {
		return nil
	}
	return func() tea.Msg {
		enriched := merge.Enrich(state, client)
		return DataLoadedMsg{State: enriched}
	}
}

// Update handles all messages for the root model.
func (a AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.ready = true
		a.updateLayout()
		return a, nil

	case tea.KeyMsg:
		// Global keys that are always active
		if key.Matches(msg, a.keys.ForceQuit) {
			return a, tea.Quit
		}

		// Help overlay intercepts everything when visible
		if a.showHelp {
			if key.Matches(msg, a.keys.Help) || key.Matches(msg, a.keys.Escape) {
				a.showHelp = false
			}
			return a, nil
		}

		// Don't intercept keys while filtering
		if a.listView.Filtering() {
			var cmd tea.Cmd
			a.listView, cmd = a.listView.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			// Update filter indicator in status bar
			a.statusBar.SetFilter(a.listView.FilterValue())
			return a, tea.Batch(cmds...)
		}

		if key.Matches(msg, a.keys.Quit) {
			return a, tea.Quit
		}

		// Help toggle
		if key.Matches(msg, a.keys.Help) {
			a.showHelp = !a.showHelp
			return a, nil
		}

		// Refresh — async API re-enrichment (SC-6, FR-14)
		if key.Matches(msg, a.keys.Refresh) {
			a.statusBar.SetRefreshing(true)
			return a, a.enrichCmd()
		}

		// Tab switching
		if key.Matches(msg, a.keys.Tab) || key.Matches(msg, a.keys.ShiftTab) {
			var cmd tea.Cmd
			a.tabBar, cmd = a.tabBar.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			return a, tea.Batch(cmds...)
		}

		if key.Matches(msg, a.keys.Tab1) {
			a.tabBar.ActiveTab = 0
			return a, func() tea.Msg { return TabChangedMsg{ActiveTab: 0} }
		}
		if key.Matches(msg, a.keys.Tab2) {
			a.tabBar.ActiveTab = 1
			return a, func() tea.Msg { return TabChangedMsg{ActiveTab: 1} }
		}
		if key.Matches(msg, a.keys.Tab3) {
			a.tabBar.ActiveTab = 2
			return a, func() tea.Msg { return TabChangedMsg{ActiveTab: 2} }
		}
		if key.Matches(msg, a.keys.Tab4) {
			a.tabBar.ActiveTab = 3
			return a, func() tea.Msg { return TabChangedMsg{ActiveTab: 3} }
		}

		// Filter key — forward to list
		if key.Matches(msg, a.keys.Filter) {
			var cmd tea.Cmd
			a.listView, cmd = a.listView.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			return a, tea.Batch(cmds...)
		}

		// Navigation — forward to list
		if key.Matches(msg, a.keys.Up) || key.Matches(msg, a.keys.Down) {
			var cmd tea.Cmd
			a.listView, cmd = a.listView.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			return a, tea.Batch(cmds...)
		}

		// Escape: clear filter if active
		if key.Matches(msg, a.keys.Escape) {
			if a.listView.FilterValue() != "" {
				var cmd tea.Cmd
				a.listView, cmd = a.listView.Update(msg)
				a.statusBar.SetFilter("")
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
				return a, tea.Batch(cmds...)
			}
		}

		// Detail viewport scroll — explicit keys for scrolling detail content.
		// PgUp/PgDn, Ctrl+U/Ctrl+D, and f/b all scroll the detail panel.
		if key.Matches(msg, a.keys.PageDown) || key.Matches(msg, a.keys.PageUp) ||
			key.Matches(msg, a.keys.HalfPageUp) || key.Matches(msg, a.keys.HalfPageDn) {
			var cmd tea.Cmd
			a.detail, cmd = a.detail.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			return a, tea.Batch(cmds...)
		}

		// Remaining keys — forward to detail viewport as catch-all
		var cmd tea.Cmd
		a.detail, cmd = a.detail.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return a, tea.Batch(cmds...)

	case tea.MouseMsg:
		// Forward mouse events to the detail viewport for scroll support.
		// The detail viewport handles mouse wheel internally.
		var cmd tea.Cmd
		a.detail, cmd = a.detail.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return a, tea.Batch(cmds...)

	case TabChangedMsg:
		a.populateList(msg.ActiveTab)
		return a, nil

	case SelectionChangedMsg:
		a.updateDetail()
		return a, nil

	case DataLoadedMsg:
		a.state = msg.State
		a.statusBar.SetOnline(msg.State.Online)
		a.statusBar.SetAPIError(msg.State.APIError)
		a.statusBar.SetRefreshing(false)

		counts := []int{
			len(a.state.Agents),
			len(a.state.Skills),
			len(a.state.MCPs),
			len(a.state.Providers),
		}
		a.tabBar.SetCounts(counts)

		a.populateList(a.tabBar.ActiveTab)
		return a, nil

	case EnrichErrorMsg:
		a.statusBar.SetRefreshing(false)
		a.statusBar.SetAPIError(msg.Err.Error())
		return a, nil
	}

	return a, tea.Batch(cmds...)
}

// appDimensions calculates the effective app width and height.
// Width uses the full terminal. Height is reduced to ~88% of the
// terminal height (capped at 50 rows) so the TUI feels compact.
func (a AppModel) appDimensions() (appWidth, appHeight int) {
	// Full terminal width — no margin.
	appWidth = a.width
	if appWidth < 20 {
		appWidth = 20
	}

	// Effective height = ~88% of terminal, capped at 50 rows.
	appHeight = a.height * 88 / 100
	const maxHeight = 50
	if appHeight > maxHeight {
		appHeight = maxHeight
	}
	if appHeight < 10 {
		appHeight = 10
	}

	return appWidth, appHeight
}

// View renders the complete UI layout.
func (a AppModel) View() string {
	if !a.ready {
		return "Loading..."
	}

	appWidth, appHeight := a.appDimensions()

	// The outer frame border consumes 2 cols and 2 rows.
	innerWidth := appWidth - 2
	innerHeight := appHeight - 2

	if innerWidth < 4 {
		innerWidth = 4
	}
	if innerHeight < 4 {
		innerHeight = 4
	}

	// Layout constants within the frame
	tabBarHeight := lipgloss.Height(a.tabBar.View())
	statusBarHeight := 1
	contentHeight := innerHeight - tabBarHeight - statusBarHeight

	if contentHeight < 1 {
		contentHeight = 1
	}

	// panelHeight is the OUTER rendered height for both panels (including border).
	// Both panels share this single value to guarantee identical visual height.
	panelHeight := contentHeight
	if panelHeight < 3 {
		panelHeight = 3 // minimum: 1 row content + 2 rows border
	}

	// panelInnerH is the content height INSIDE each panel's border (top+bottom = 2).
	panelInnerH := panelHeight - 2
	if panelInnerH < 1 {
		panelInnerH = 1
	}

	// Split: list (30%) | detail (70%)
	listWidth := innerWidth * 30 / 100
	detailWidth := innerWidth - listWidth

	if listWidth < 10 {
		listWidth = 10
	}
	if detailWidth < 10 {
		detailWidth = 10
	}

	// Panel inner dimensions (widths).
	// ListPanel has border only (no padding): 1 left + 1 right = 2 horizontal.
	// DetailPanel has border (2) + PaddingLeft(1) + PaddingRight(1).
	// lipgloss Width(n) is the content width INSIDE the style — padding is
	// subtracted automatically, so we only subtract the border here.
	listInnerW := listWidth - 2        // Width for panel content (no padding on ListPanel)
	detailPanelW := detailWidth - 2    // Width for DetailPanel style (border only; padding handled by lipgloss)
	detailContentW := detailPanelW - 2 // Content width for viewport (inside padding)

	if listInnerW < 1 {
		listInnerW = 1
	}
	if detailPanelW < 1 {
		detailPanelW = 1
	}
	if detailContentW < 1 {
		detailContentW = 1
	}

	// Update sub-component sizes (inner dimensions for the content)
	a.listView.SetSize(listInnerW, panelInnerH)
	a.detail.SetSize(detailContentW, panelInnerH)

	// Render each section.
	// CRITICAL: Both panels use Height(panelHeight) on their style to FORCE
	// identical outer rendered height, regardless of content length.
	tabBarView := a.tabBar.View()

	listRendered := ListPanel.
		Width(listInnerW).
		Height(panelHeight).
		Render(a.listView.View())

	detailRendered := DetailPanel.
		Width(detailPanelW).
		Height(panelHeight).
		Render(a.detail.View())

	contentView := lipgloss.JoinHorizontal(lipgloss.Top, listRendered, detailRendered)

	statusBarView := a.statusBar.View()

	innerUI := lipgloss.JoinVertical(lipgloss.Left,
		tabBarView,
		contentView,
		statusBarView,
	)

	// Wrap everything in the outer app frame
	ui := AppFrame.
		Width(innerWidth).
		Height(innerHeight).
		Render(innerUI)

	// Help overlay on top of everything (SC-10, FR-15)
	if a.showHelp {
		a.help.SetSize(appWidth, appHeight)
		return lipgloss.Place(a.width, a.height,
			lipgloss.Left, lipgloss.Center,
			a.help.View(ui),
			lipgloss.WithWhitespaceChars(" "),
		)
	}

	// Return the framed app directly — no left margin.
	return ui
}

// updateLayout recalculates all sub-component sizes.
func (a *AppModel) updateLayout() {
	appWidth, appHeight := a.appDimensions()

	// Outer frame border: 2 cols, 2 rows
	innerWidth := appWidth - 2
	innerHeight := appHeight - 2

	if innerWidth < 4 {
		innerWidth = 4
	}
	if innerHeight < 4 {
		innerHeight = 4
	}

	a.tabBar.SetWidth(innerWidth)
	a.statusBar.SetWidth(innerWidth)

	tabBarHeight := 2 // tab bar is typically 2 lines
	statusBarHeight := 1
	contentHeight := innerHeight - tabBarHeight - statusBarHeight

	if contentHeight < 1 {
		contentHeight = 1
	}

	// panelHeight is the OUTER rendered height for both panels (including border).
	panelHeight := contentHeight
	if panelHeight < 3 {
		panelHeight = 3
	}

	// panelInnerH is the content height INSIDE each panel's border (top+bottom = 2).
	panelInnerH := panelHeight - 2
	if panelInnerH < 1 {
		panelInnerH = 1
	}

	// Each panel has its own border (2 chars each side)
	listWidth := innerWidth * 30 / 100
	detailWidth := innerWidth - listWidth

	// ListPanel has border only (no padding): 1 left + 1 right = 2 horizontal.
	// DetailPanel has border (2) + PaddingLeft(1) + PaddingRight(1).
	// lipgloss Width(n) is the content width INSIDE the style — padding is
	// subtracted automatically, so we only subtract the border here.
	listInnerW := listWidth - 2         // Width for panel (no padding on ListPanel)
	detailInnerW := detailWidth - 2 - 2 // border(2) + padding(2)

	if listInnerW < 1 {
		listInnerW = 1
	}
	if detailInnerW < 1 {
		detailInnerW = 1
	}

	a.listView.SetSize(listInnerW, panelInnerH)
	a.detail.SetSize(detailInnerW, panelInnerH)
}

// populateList fills the list with items for the given tab.
func (a *AppModel) populateList(tab int) {
	if a.state == nil {
		return
	}

	var items []list.Item
	switch tab {
	case 0:
		items = AgentsToItems(a.state.Agents)
	case 1:
		items = SkillsToItems(a.state.Skills)
	case 2:
		items = MCPsToItems(a.state.MCPs)
	case 3:
		items = ProvidersToItems(a.state.Providers)
	}

	a.listView.SetItems(items)
	a.updateDetail()
}

// updateDetail renders detail content for the currently selected item.
func (a *AppModel) updateDetail() {
	item := a.listView.SelectedItem()
	if item == nil {
		a.detail.SetContent(DetailMuted.Render("No item selected"))
		return
	}

	var content string
	switch e := item.Entity().(type) {
	case *model.Agent:
		content = RenderAgentDetail(e)
	case *model.Skill:
		content = RenderSkillDetail(e)
	case *model.MCPServer:
		content = RenderMCPDetail(e)
	case *model.Provider:
		content = RenderProviderDetail(e)
	default:
		content = DetailMuted.Render("Unknown entity type")
	}

	a.detail.SetContent(content)
}
