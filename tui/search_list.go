package tui

import (
	"capsule-package-manager/models"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
	"github.com/Jguer/go-alpm/v2"
)

const (
	ColorPrimary      = "#ff5555"
	HighlightColour   = "#ff5555"
	DescriptionColour = "#94a3b8"
	ColorText         = "#f8f9fa"
	ColorMuted        = "#6272a4"
	ColorBg           = "#151827"
)

type searchItem struct {
	packageInfo models.DetailedPackageInfo
	installed   bool
}

type Model struct {
	list         list.Model
	detail       DetailModel
	showDetails  bool
	termWidth    int
	termHeight   int
	packagesList []models.DetailedPackageInfo
}

const (
	maxListWidth  = 80
	maxListHeight = 40
)

// Calculates list dimensions from termWidth and termHeight
// Called by Model.Update() and Model.View() in search_list.go
// Returns listWidth and listHeight to the caller
func searchListSize(termWidth, termHeight int) (int, int) {
	listWidth := termWidth - 8
	listHeight := termHeight - 8

	if listWidth > maxListWidth {
		listWidth = maxListWidth
	}
	if listHeight > maxListHeight {
		listHeight = maxListHeight
	}
	if listWidth < 1 {
		listWidth = 1
	}
	if listHeight < 1 {
		listHeight = 1
	}

	return listWidth, listHeight
}

// Builds a Bubble Tea list from matching packages
// Gets packages from BasicPackageSearch() in search/search_packages.go
// Returns Model to BasicPackageSearch() in search/search_packages.go
func SearchList(packages []models.DetailedPackageInfo) Model {
	items := []list.Item{}
	installedPackages := installedPackageNames()

	log.Debug("Compiling search results")
	for _, packageInfo := range packages {
		items = append(items, searchItem{
			packageInfo: packageInfo,
			installed:   installedPackages[packageInfo.Name],
		})
	}

	// SearchList styling here
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(HighlightColour))

	delegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorMuted))

	delegate.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorText))

	delegate.Styles.NormalDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color(DescriptionColour))

	l := list.New(items, delegate, maxListWidth, maxListHeight)
	l.Title = "Capsule Search"

	styles := list.DefaultStyles(true)

	// Title bar styling
	styles.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#151827")).
		Background(lipgloss.Color(ColorPrimary)).
		Padding(0, 1)

	// Status bar / pagination counter styling
	styles.StatusBar = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorMuted)).
		PaddingLeft(2)

	// Help keys styling (bottom text)
	styles.HelpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorMuted)).
		Padding(1, 0, 0, 2)

	l.Styles = styles

	return Model{list: l, packagesList: packages}
}

// Returns search item title, marking installed packages
// Called by the Bubble Tea list delegate; uses searchItem.packageInfo
func (i searchItem) Title() string {
	if i.installed {
		return i.packageInfo.Name + " [Installed]"
	}
	return i.packageInfo.Name
}

// Returns package description from searchItem.packageInfo
// Called through searchItem's list item methods
func (i searchItem) Version() string {
	return i.packageInfo.Desc
}

// Returns package description from searchItem.packageInfo
// Called by the Bubble Tea list delegate
func (i searchItem) Description() string {
	return i.packageInfo.Desc
}

// Returns package repo from searchItem.packageInfo
// Available to callers of searchItem.Repo()
func (i searchItem) Repo() string {
	return i.packageInfo.Repo
}

// Returns package name for Bubble Tea list filtering
// Called by the Bubble Tea list
func (i searchItem) FilterValue() string {
	return i.packageInfo.Name
}

// Supplies the initial Bubble Tea command for the search Model
// Called by Bubble Tea; returns nil
func (m Model) Init() tea.Cmd {
	return nil
}

// Handles keyboard, resize and detail-panel events
// Gets msg from Bubble Tea; uses DetailView() in detailed_viewport.go
// Returns updated Model and command to Bubble Tea
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
	case tea.KeyPressMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			log.Debug("Quitting Capsule")
			return m, tea.Quit
		case "esc":
			if m.showDetails {
				log.Debug("Returning to search list")
				m.showDetails = false
			} else {
				log.Debug("Quitting Capsule")
				return m, tea.Quit
			}
		case "enter":
			_, ok := m.list.SelectedItem().(searchItem)
			if ok {
				log.Debug("Opening package details")
				m.detail, _ = DetailView(m.packagesList, m.list.GlobalIndex(), m.termWidth, m.termHeight)
				m.showDetails = true
			}
		}
	}
	var cmd tea.Cmd

	if m.showDetails {
		detail, cmd := m.detail.Update(msg)
		m.detail = detail.(DetailModel)
		return m, cmd
	} else {
		m.list, cmd = m.list.Update(msg)
		listWidth, listHeight := searchListSize(m.termWidth, m.termHeight)
		m.list.SetWidth(listWidth)
		m.list.SetHeight(listHeight)
		return m, cmd
	}
}

// Renders the search list and optional detail overlay
// Uses Model fields updated by Update() in search_list.go
// Returns tea.View to Bubble Tea
func (m Model) View() tea.View {
	listWidth, _ := searchListSize(m.termWidth, m.termHeight)
	listContent := lipgloss.NewStyle().
		Width(listWidth).
		Render(m.list.View())

	list := lipgloss.Place(
		m.termWidth,
		m.termHeight,
		lipgloss.Center,
		lipgloss.Center,
		listContent,
	)
	if m.showDetails {
		detailsView := m.detail.View()
		details := detailsView.Content
		posX := (m.termWidth - lipgloss.Width(details)) / 2
		posY := (m.termHeight - lipgloss.Height(details)) / 2

		compositor := lipgloss.NewCompositor(
			lipgloss.NewLayer(list).X(0).Y(0).Z(0),
			lipgloss.NewLayer(details).X(posX).Y(posY).Z(1),
		)
		composed := compositor.Render()
		finalView := tea.NewView(composed)
		finalView.AltScreen = true

		return finalView
	} else {
		listView := tea.NewView(list)
		listView.AltScreen = true
		return listView
	}
}

// Reads installed package names from the local ALPM database
// Called by SearchList() in search_list.go
// Returns installedPackages map to SearchList()
func installedPackageNames() map[string]bool {
	installedPackages := map[string]bool{}
	handle, err := alpm.Initialize("/", "/var/lib/pacman/")
	if err != nil {
		return installedPackages
	}
	defer handle.Release()

	localDB, err := handle.LocalDB()
	if err != nil {
		return installedPackages
	}

	for _, packageInfo := range localDB.PkgCache().Slice() {
		installedPackages[packageInfo.Name()] = true
	}
	return installedPackages
}
