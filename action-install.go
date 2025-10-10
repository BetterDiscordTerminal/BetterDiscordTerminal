package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mritd/bubbles/common"
	"github.com/mritd/bubbles/selector"
)

var (
	// Box styles
	installingBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("12")).
			Padding(1, 2).
			Width(40)

	successBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("10")).
			Padding(1, 2).
			Width(40)

	errorBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("9")).
			Padding(1, 2).
			Width(40)

	// Text styles
	successText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true)

	errorText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true)

	subtleText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))
)

type installModel struct {
	sl         selector.Model
	installing bool
	spinner    spinner.Model
	result     installationResult
}

type installationResult struct {
	version DiscordVersion
	err     error
}

func initInstallModel() installModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return installModel{
		sl: selector.Model{
			Data:    discordVersions,
			PerPage: 3,
			HeaderFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
				return lipgloss.NewStyle().
					Foreground(lipgloss.Color("10")).
					Bold(true).
					Render("BDTerm") +
					lipgloss.NewStyle().
						Foreground(lipgloss.Color("10")).
						Render("\nChoose a Discord version to install:")

			},
			SelectedFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
				v := obj.(DiscordVersion)
				return common.FontColor(fmt.Sprintf("[%d] %s (%s)", gdIndex+1, v.Name, v.Description), selector.ColorSelected)
			},
			UnSelectedFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
				v := obj.(DiscordVersion)
				return common.FontColor(fmt.Sprintf(" %d. %s (%s)", gdIndex+1, v.Name, v.Description), selector.ColorUnSelected)
			},
			FooterFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
				return common.FontColor("\nUse ↑/↓ to navigate, enter to select, ESC to return to menu, q/ctrl-c to quit.", selector.ColorFooter) + "\n" +
					common.FontColor("BetterDiscord Terminal by nmn (https://github.com/pandeynmn/bdterm)", "8")
			},
			FinishedFunc: func(s interface{}) string {
				return ""
			},
		},
		spinner:    s,
		installing: false,
	}
}

func (m installModel) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m installModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" && !m.installing {
			return m, func() tea.Msg { return switchToMainMenu{} }
		}

	case string:
		if msg == common.DONE && !m.installing {
			selected := m.sl.Selected().(DiscordVersion)
			m.installing = true

			return m, tea.Batch(
				m.spinner.Tick,
				func() tea.Msg {
					return installationResult{
						version: selected,
						err:     RunInstallation(selected),
					}
				},
			)
		}

	case installationResult:
		m.result = msg
		m.installing = false
		return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
			return tea.Quit()
		})

	case spinner.TickMsg:
		if m.installing && m.result.version.Name == "" {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	if !m.installing {
		_, cmd := m.sl.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m installModel) View() string {
	// Installing view with spinner
	if m.installing && m.result.version.Name == "" {
		content := fmt.Sprintf("%s  Installing BetterDiscord", m.spinner.View())
		return "\n" + installingBox.Render(content) + "\n\n"
	}

	// Result views
	if m.result.version.Name != "" {
		if m.result.err != nil {
			// Error view
			title := errorText.Render("Installation Failed")
			errorMsg := subtleText.Render("\n\n" + m.result.err.Error())
			return "\n" + errorBox.Render(title+errorMsg) + "\n\n"
		}

		// Success view
		title := successText.Render("Installation Complete!")
		details := fmt.Sprintf("\n\nBetterDiscord installed for %s", m.result.version.Name)
		launching := subtleText.Render("\nDiscord is launching...")
		return "\n" + successBox.Render(title+details+launching) + "\n\n"
	}

	return m.sl.View()
}
