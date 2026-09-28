package main

import (
	"fmt"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type detailModel struct {
	viewport     viewport.Model
	lineCount    int
	panelWidth   int
	panelHeight  int
	packageName  string
	repo         string
	packagesList []DetailedPackageInfo
	selectedRepo int
	termWidth    int
	termHeight   int
	standalone   bool
}

func DetailView(packageInfo []DetailedPackageInfo, selectedRepo, width, height int) (detailModel, error) {
	viewPort := viewport.New()
	initialValues := detailModel{
		viewport:     viewPort,
		packagesList: packageInfo,
		selectedRepo: selectedRepo,
	}
	err := initialValues.DetailedViewUpdate(width, height)
	if err != nil {
		return detailModel{}, fmt.Errorf("unable to update viewport: %s", err)
	}
	return initialValues, nil
}

func (m *detailModel) DetailedViewUpdate(width, height int) error {
	m.termWidth, m.termHeight = width, height
	viewPortWidth := detailPanelWidth(width) - 6
	selectedPackage := m.packagesList[m.selectedRepo]
	formattedPackage, err := DetailedPackageOutput(selectedPackage.Name, []DetailedPackageInfo{m.packagesList[m.selectedRepo]}, viewPortWidth)
	if err != nil {
		return fmt.Errorf("unable to format package: %s | %s", selectedPackage.Name, err)
	}
	m.lineCount = lipgloss.Height(formattedPackage)
	m.panelWidth, m.panelHeight = detailPanelSize(width, height, m.lineCount)
	m.viewport.SetContent(formattedPackage)
	m.viewport.SetWidth(viewPortWidth)
	m.viewport.SetHeight(m.panelHeight - 9)
	m.packageName = selectedPackage.Name
	m.repo = selectedPackage.Repo
	return nil
}

func (m detailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth, m.termHeight = msg.Width, msg.Height
		m.panelWidth, m.panelHeight = detailPanelSize(msg.Width, msg.Height, m.lineCount)
		m.viewport.SetWidth(m.panelWidth - 6)
		m.viewport.SetHeight(m.panelHeight - 9)

	case tea.KeyPressMsg:
		switch msg.String() {
		case "left":
			if m.selectedRepo > 0 {
				m.selectedRepo -= 1
				m.DetailedViewUpdate(m.termWidth, m.termHeight)
			}
		case "right":
			if m.selectedRepo < len(m.packagesList)-1 {
				m.selectedRepo += 1
				m.DetailedViewUpdate(m.termWidth, m.termHeight)
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.standalone {
				return m, tea.Quit
			}
		}

	}
	return m, nil
}

func (m detailModel) Init() tea.Cmd {
	return nil
}

func (m detailModel) View() tea.View {
	const (
		panelBackground = "#151827"
		panelBorder     = "#a78bfa"
		textColour      = "#e2e8f0"
		titleColour     = "#f9a8d4"
		accentColour    = "#5eead4"
		mutedColour     = "#94a3b8"
	)

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(titleColour)).
		Render("  Capsule Package Details")

	meta := lipgloss.NewStyle().
		Foreground(lipgloss.Color(accentColour)).
		Render("  " + m.packageName + "  •  " + m.repo)

	body := lipgloss.NewStyle().
		Foreground(lipgloss.Color(textColour)).
		Render(m.viewport.View())

	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color(mutedColour)).
		Render("  <-/-> to switch repo  •  q/ctrl-c to quit  •  Esc to close panel")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		meta,
		"",
		body,
		"",
		footer,
	)

	modal := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(panelBorder)).
		Background(lipgloss.Color(panelBackground)).
		Foreground(lipgloss.Color(textColour)).
		Padding(1, 2).
		Width(m.panelWidth).
		Height(m.panelHeight).
		Render(content)

	if m.standalone {
		centeredModal := lipgloss.Place(m.termWidth, m.termHeight, lipgloss.Center, lipgloss.Center, modal)
		standaloneView := tea.NewView(centeredModal)
		standaloneView.AltScreen = true
		return standaloneView
	}

	return tea.NewView(modal)
}

func detailPanelWidth(width int) int {
	return width * 3 / 4
}

func detailPanelSize(width, height, contentLines int) (int, int) {
	panelWidth := detailPanelWidth(width)
	panelHeight := contentLines + 9
	maxPanelHeight := height - 4

	if panelHeight > maxPanelHeight {
		panelHeight = maxPanelHeight
	}
	if panelHeight < 12 {
		panelHeight = 12
	}

	return panelWidth, panelHeight
}
