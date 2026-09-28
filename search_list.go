package main

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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

func SearchList(packages []DetailedPackageInfo) model {
	items := []list.Item{}

	for _, packageInfo := range packages {
		items = append(items, searchItem{packageInfo: packageInfo})
	}
	l := list.New(items, list.NewDefaultDelegate(), 80, 20)
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
			return m, tea.Quit
		case "esc":
			if m.showDetails {
				m.showDetails = false
			} else {
				return m, tea.Quit
			}
		case "enter":
			_, ok := m.list.SelectedItem().(searchItem)
			if ok {
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
		return m, cmd
	}
}

func (m model) View() tea.View {
	list := lipgloss.Place(
		m.termWidth,
		m.termHeight,
		lipgloss.Center,
		lipgloss.Center,
		m.list.View(),
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

/*
delegate := list.NewDefaultDelegate()
	delegate.SetHeight(3)
	delegate.SetSpacing(1)
	delegate.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#e2e8f0")).
		PaddingLeft(2)
	delegate.Styles.NormalDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#94a3b8")).
		PaddingLeft(2)
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#f9a8d4")).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color("#5eead4")).
		PaddingLeft(1)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#5eead4")).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color("#5eead4")).
		PaddingLeft(1)

	l := list.New(items, delegate, 80, 20)
	styles := list.DefaultStyles(true)
	styles.TitleBar = lipgloss.NewStyle().PaddingLeft(2)
	styles.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#151827")).
		Background(lipgloss.Color("#a78bfa")).
		Padding(0, 1)
	styles.StatusBar = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#94a3b8")).
		PaddingLeft(2)
	styles.HelpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#94a3b8")).
		Padding(1, 0, 0, 2)
	l.Styles = styles
	l.Title = "✦ Capsule Search" */
