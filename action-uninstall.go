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

type uninstallModel struct {
	sl           selector.Model
	uninstalling bool
	spinner      spinner.Model
	result       uninstallationResult
}

type uninstallationResult struct {
	version DiscordVersion
	err     error
}

func initUninstallModel() uninstallModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return uninstallModel{
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
						Render("\nChoose a Discord version to uninstall:")
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
		spinner:      s,
		uninstalling: false,
	}
}

func (m uninstallModel) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m uninstallModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" && !m.uninstalling {
			return m, func() tea.Msg { return switchToMainMenu{} }
		}

	case string:
		if msg == common.DONE && !m.uninstalling {
			selected := m.sl.Selected().(DiscordVersion)
			m.uninstalling = true

			return m, tea.Batch(
				m.spinner.Tick,
				func() tea.Msg {
					return uninstallationResult{
						version: selected,
						err:     RunUninstall(selected),
					}
				},
			)
		}

	case uninstallationResult:
		m.result = msg
		m.uninstalling = false
		return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
			return tea.Quit()
		})

	case spinner.TickMsg:
		if m.uninstalling && m.result.version.Name == "" {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	if !m.uninstalling {
		_, cmd := m.sl.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m uninstallModel) View() string {
	// When uninstalling, DON'T show the selector at all
	if m.uninstalling && m.result.version.Name == "" {
		content := fmt.Sprintf("%s  Uninstalling BetterDiscord...", m.spinner.View())
		return "\n" + installingBox.Render(content) + "\n\n"
	}

	// Show result
	if m.result.version.Name != "" {
		if m.result.err != nil {
			title := errorText.Render("✗ Uninstallation Failed")
			errorMsg := subtleText.Render("\n\n" + m.result.err.Error())
			return "\n" + errorBox.Render(title+errorMsg) + "\n\n"
		}

		title := successText.Render("✓ Uninstallation Complete!")
		details := fmt.Sprintf("\n\nBetterDiscord removed from %s", m.result.version.Name)
		return "\n" + successBox.Render(title+details) + "\n\n"
	}

	// Only show selector when NOT uninstalling
	return m.sl.View()
}
