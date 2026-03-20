package tui

import (
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jarlex/opencode-configurator/internal/api"
	"github.com/jarlex/opencode-configurator/internal/merge"
	"github.com/jarlex/opencode-configurator/internal/model"
	"sort"
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

	width            int
	height           int
	ready            bool
	ShowHidden       bool
	FullScreenDetail bool
	SplitRatio       int
	ConfigPath       string
	SuccessMsg       string
}

// NewApp creates the root application model with initial state.
func NewApp(state *model.AppState, client *api.Client, splitRatio int, configPath string) AppModel {
	keys := DefaultKeyMap()
	app := AppModel{
		SplitRatio: splitRatio,
		ConfigPath: configPath,
		state:      state,
		tabBar:     NewTabBar(keys),
		listView:   NewListView(keys),
		detail:     NewDetailView(),
		statusBar:  NewStatusBar(),
		help:       NewHelpOverlay(keys),
		keys:       keys,
		apiClient:  client,
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
		return a.handleWindowSize(msg)
	case tea.KeyMsg:
		return a.handleKeyMsg(msg)
	case tea.MouseMsg:
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
		return a.handleDataLoaded(msg)
	case EnrichErrorMsg:
		a.statusBar.SetRefreshing(false)
		a.statusBar.SetAPIError(msg.Err.Error())
		return a, nil
	}

	return a, tea.Batch(cmds...)
}

func (a AppModel) handleWindowSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	a.width = msg.Width
	a.height = msg.Height
	a.ready = true
	a.updateLayout()
	return a, nil
}

func (a AppModel) handleDataLoaded(msg DataLoadedMsg) (tea.Model, tea.Cmd) {
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
}

func (a AppModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if !key.Matches(msg, a.keys.CopyDetail) {
		a.SuccessMsg = ""
	}

	if key.Matches(msg, a.keys.ForceQuit) {
		return a, tea.Quit
	}

	if a.showHelp {
		if key.Matches(msg, a.keys.Help) || key.Matches(msg, a.keys.Escape) {
			a.showHelp = false
		}
		return a, nil
	}

	if a.listView.Filtering() {
		var cmd tea.Cmd
		a.listView, cmd = a.listView.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		a.statusBar.SetFilter(a.listView.FilterValue())
		return a, tea.Batch(cmds...)
	}

	if key.Matches(msg, a.keys.Quit) {
		return a, tea.Quit
	}

	if key.Matches(msg, a.keys.Help) {
		a.showHelp = !a.showHelp
		return a, nil
	}

	if key.Matches(msg, a.keys.Refresh) {
		a.statusBar.SetRefreshing(true)
		return a, a.enrichCmd()
	}

	if handled, model, cmd := a.handleTabKeys(msg); handled {
		return model, cmd
	}

	if handled, model, cmd := a.handleActionKeys(msg); handled {
		return model, cmd
	}

	if handled, model, cmd := a.handleListKeys(msg); handled {
		return model, cmd
	}

	// Remaining keys forward to detail
	var cmd tea.Cmd
	a.detail, cmd = a.detail.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return a, tea.Batch(cmds...)
}

func (a AppModel) handleTabKeys(msg tea.KeyMsg) (bool, tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	if key.Matches(msg, a.keys.Tab) || key.Matches(msg, a.keys.ShiftTab) {
		var cmd tea.Cmd
		a.tabBar, cmd = a.tabBar.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return true, a, tea.Batch(cmds...)
	}

	return false, a, nil
}

func (a AppModel) handleActionKeys(msg tea.KeyMsg) (bool, tea.Model, tea.Cmd) {
	if key.Matches(msg, a.keys.ToggleHidden) {
		a.ShowHidden = !a.ShowHidden
		a.populateList(a.tabBar.ActiveTab)
		return true, a, nil
	}
	if key.Matches(msg, a.keys.FullScreen) {
		a.FullScreenDetail = !a.FullScreenDetail
		return true, a, nil
	}
	if key.Matches(msg, a.keys.CopyDetail) {
		err := clipboard.WriteAll(a.detail.Content())
		if err == nil {
			a.SuccessMsg = "Copied to clipboard"
		} else {
			a.SuccessMsg = "Copy failed"
		}
		return true, a, nil
	}
	return false, a, nil
}

func (a AppModel) handleListKeys(msg tea.KeyMsg) (bool, tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	if key.Matches(msg, a.keys.Filter) || key.Matches(msg, a.keys.Up) || key.Matches(msg, a.keys.Down) {
		var cmd tea.Cmd
		a.listView, cmd = a.listView.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return true, a, tea.Batch(cmds...)
	}

	if key.Matches(msg, a.keys.Escape) {
		if a.listView.FilterValue() != "" {
			var cmd tea.Cmd
			a.listView, cmd = a.listView.Update(msg)
			a.statusBar.SetFilter("")
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			return true, a, tea.Batch(cmds...)
		}
	}

	if key.Matches(msg, a.keys.PageDown) || key.Matches(msg, a.keys.PageUp) ||
		key.Matches(msg, a.keys.HalfPageUp) || key.Matches(msg, a.keys.HalfPageDn) {
		var cmd tea.Cmd
		a.detail, cmd = a.detail.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return true, a, tea.Batch(cmds...)
	}

	return false, a, nil
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
	panelHeight := contentHeight
	if panelHeight < 3 {
		panelHeight = 3 // minimum: 1 row content + 2 rows border
	}

	// panelInnerH is the content height INSIDE each panel's border (top+bottom = 2).
	panelInnerH := panelHeight - 2
	if panelInnerH < 1 {
		panelInnerH = 1
	}

	listWidth, detailWidth := a.calculatePanelWidths(innerWidth)

	listInnerW := listWidth - 2        // Width for panel content
	detailPanelW := detailWidth - 2    // Width for DetailPanel style
	detailContentW := detailPanelW - 2 // Content width for viewport

	if listInnerW < 1 {
		listInnerW = 0
	}
	if detailPanelW < 1 {
		detailPanelW = 1
	}
	if detailContentW < 1 {
		detailContentW = 1
	}

	// Update sub-component sizes
	a.listView.SetSize(listInnerW, panelInnerH)
	a.detail.SetSize(detailContentW, panelInnerH)

	innerUI := a.renderInnerUI(listInnerW, detailPanelW, panelHeight)

	// Wrap everything in the outer app frame
	ui := AppFrame.
		Width(innerWidth).
		Height(innerHeight).
		Render(innerUI)

	// Help overlay on top of everything
	if a.showHelp {
		a.help.SetSize(appWidth, appHeight)
		return lipgloss.Place(a.width, a.height,
			lipgloss.Left, lipgloss.Center,
			a.help.View(ui),
			lipgloss.WithWhitespaceChars(" "),
		)
	}

	return ui
}

func (a AppModel) calculatePanelWidths(innerWidth int) (int, int) {
	ratio := a.SplitRatio
	if ratio <= 0 || ratio > 90 {
		ratio = 30
	}
	listWidth := innerWidth * ratio / 100
	if a.FullScreenDetail {
		listWidth = 0
	}
	detailWidth := innerWidth - listWidth

	if !a.FullScreenDetail && listWidth < 10 {
		listWidth = 10
	}
	if detailWidth < 10 {
		detailWidth = 10
	}
	return listWidth, detailWidth
}

func (a AppModel) renderInnerUI(listInnerW, detailPanelW, panelHeight int) string {
	tabBarView := a.tabBar.View()

	listRendered := ListPanel.
		Width(listInnerW).
		Height(panelHeight).
		Render(a.listView.View())

	detailRendered := DetailPanel.
		Width(detailPanelW).
		Height(panelHeight).
		Render(a.detail.View())

	contentView := detailRendered
	if !a.FullScreenDetail {
		contentView = lipgloss.JoinHorizontal(lipgloss.Top, listRendered, detailRendered)
	}

	a.statusBar.ConfigPath = a.ConfigPath
	a.statusBar.SuccessMsg = a.SuccessMsg
	a.statusBar.ScrollPercent = a.detail.ScrollPercent()

	statusBarView := a.statusBar.View()

	return lipgloss.JoinVertical(lipgloss.Left,
		tabBarView,
		contentView,
		statusBarView,
	)
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
	ratio := a.SplitRatio
	if ratio <= 0 || ratio > 90 {
		ratio = 30
	}
	listWidth := innerWidth * ratio / 100
	if a.FullScreenDetail {
		listWidth = 0
	}
	detailWidth := innerWidth - listWidth

	// ListPanel has border only (no padding): 1 left + 1 right = 2 horizontal.
	// DetailPanel has border (2) + PaddingLeft(1) + PaddingRight(1).
	// lipgloss Width(n) is the content width INSIDE the style — padding is
	// subtracted automatically, so we only subtract the border here.
	listInnerW := listWidth - 2         // Width for panel (no padding on ListPanel)
	detailInnerW := detailWidth - 2 - 2 // border(2) + padding(2)

	if listInnerW < 1 {
		listInnerW = 0
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
		var filtered []model.Agent
		for _, agent := range a.state.Agents {
			if !agent.Hidden || a.ShowHidden {
				filtered = append(filtered, agent)
			}
		}
		items = AgentsToItems(filtered)
	case 1:
		items = SkillsToItems(a.state.Skills)
	case 2:
		items = MCPsToItems(a.state.MCPs)
	case 3:
		items = ProvidersToItems(a.state.Providers)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].FilterValue() < items[j].FilterValue()
	})

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
