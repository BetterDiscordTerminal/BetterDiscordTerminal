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

type repairModel struct {
	sl        selector.Model
	repairing bool
	spinner   spinner.Model
	result    repairResult
}

type repairResult struct {
	version DiscordVersion
	err     error
}

func initRepairModel() repairModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return repairModel{
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
						Render("\nChoose a Discord version to repair:")
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
				return common.FontColor("\nUse ↑/↓ to navigate, enter to select, ESC to return to menu, q/ctrl-c to quit.", selector.ColorFooter)
			},
			FinishedFunc: func(s interface{}) string {
				return ""
			},
		},
		spinner:   s,
		repairing: false,
	}
}

func (m repairModel) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m repairModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" && !m.repairing {
			return m, func() tea.Msg { return switchToMainMenu{} }
		}

	case string:
		if msg == common.DONE && !m.repairing {
			selected := m.sl.Selected().(DiscordVersion)
			m.repairing = true

			return m, tea.Batch(
				m.spinner.Tick,
				func() tea.Msg {
					return repairResult{
						version: selected,
						err:     RunRepair(selected),
					}
				},
			)
		}

	case repairResult:
		m.result = msg
		m.repairing = false
		return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
			return tea.Quit()
		})

	case spinner.TickMsg:
		if m.repairing && m.result.version.Name == "" {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	if !m.repairing {
		_, cmd := m.sl.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m repairModel) View() string {
	// Repairing view with spinner
	if m.repairing && m.result.version.Name == "" {
		content := fmt.Sprintf("%s  Reinstalling BetterDiscord", m.spinner.View())
		return "\n" + installingBox.Render(content) + "\n\n"
	}

	// Result views
	if m.result.version.Name != "" {
		if m.result.err != nil {
			// Error view
			title := errorText.Render("Repair Failed")
			errorMsg := subtleText.Render("\n\n" + m.result.err.Error())
			return "\n" + errorBox.Render(title+errorMsg) + "\n\n"
		}

		// Success view
		title := successText.Render("Repair Complete!")
		details := fmt.Sprintf("\n\nBetterDiscord repaired for %s", m.result.version.Name)
		launching := subtleText.Render("\nDiscord is launching...")
		return "\n" + successBox.Render(title+details+launching) + "\n\n"
	}

	return m.sl.View()
}
