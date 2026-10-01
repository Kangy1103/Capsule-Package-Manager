package main

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
)

type searchItem struct {
	packageInfo DetailedPackageInfo
}

type model struct {
	list         list.Model
	detail       detailModel
	showDetails  bool
	termWidth    int
	termHeight   int
	packagesList []DetailedPackageInfo
}

const (
	maxListWidth  = 80
	maxListHeight = 40
)

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

func SearchList(packages []DetailedPackageInfo) model {
	items := []list.Item{}
	log.Debug("Compiling search results")
	for _, packageInfo := range packages {
		items = append(items, searchItem{packageInfo: packageInfo})
	}
	l := list.New(items, list.NewDefaultDelegate(), maxListWidth, maxListHeight)
	l.Title = "Capsule Search"

	return model{list: l, packagesList: packages}
}

func (i searchItem) Title() string {
	return i.packageInfo.Name
}

func (i searchItem) Version() string {
	return i.packageInfo.Desc
}

func (i searchItem) Description() string {
	return i.packageInfo.Desc
}

func (i searchItem) Repo() string {
	return i.packageInfo.Repo
}

func (i searchItem) FilterValue() string {
	return i.packageInfo.Name
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		m.detail = detail.(detailModel)
		return m, cmd
	} else {
		m.list, cmd = m.list.Update(msg)
		listWidth, listHeight := searchListSize(m.termWidth, m.termHeight)
		m.list.SetWidth(listWidth)
		m.list.SetHeight(listHeight)
		return m, cmd
	}
}

func (m model) View() tea.View {
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
