package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mritd/bubbles/common"
	"github.com/mritd/bubbles/selector"
)

type menuModel struct {
	sl selector.Model
}

func initMenuModel() menuModel {
	return menuModel{
		sl: selector.Model{
			Data: []interface{}{
				TypeMessage{"Install", "Install BetterDiscord"},
				TypeMessage{"Uninstall", "Uninstall BetterDiscord"},
				TypeMessage{"Repair", "Repair BetterDiscord"},
			},
			PerPage: 5,
			HeaderFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
				return lipgloss.NewStyle().
					Foreground(lipgloss.Color("10")).
					Bold(true).
					Render("BDTerm") +
					lipgloss.NewStyle().
						Foreground(lipgloss.Color("10")).
						Render("\ngithub.com/pandeynmn/bdterm")

			},
			SelectedFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
				t := obj.(TypeMessage)
				return common.FontColor(fmt.Sprintf("[%d] %s (%s)", gdIndex+1, t.Type, t.Description), selector.ColorSelected)
			},
			UnSelectedFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
				t := obj.(TypeMessage)
				return common.FontColor(fmt.Sprintf(" %d. %s (%s)", gdIndex+1, t.Type, t.Description), selector.ColorUnSelected)
			},
			FooterFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
				footerTpl := `
Use ↑/↓ to navigate, enter to select, q/ctrl-c to quit.

By executing this application, you acknowledge and agree to be bound by the terms and conditions located at 
https://www.apache.org/licenses/LICENSE-2.0.txt
`
				return common.FontColor(fmt.Sprintf(footerTpl), selector.ColorFooter) +
					common.FontColor("BetterDiscord Terminal by nmn (https://github.com/pandeynmn/bdterm)", "8")
			},
			FinishedFunc: func(s interface{}) string {
				return common.FontColor("Current selected: ", selector.ColorFinished) + s.(TypeMessage).Type + "\n"
			},
		},
	}
}

func (m menuModel) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg {
	case common.DONE:
		selected := m.sl.Selected().(TypeMessage)
		switch selected.Type {
		case "Install":
			return m, func() tea.Msg { return switchToInstallModel{} }
		case "Uninstall":
			return m, func() tea.Msg { return switchToUninstallModel{} }
		case "Repair":
			return m, func() tea.Msg { return switchToRepairModel{} }
		}
	}

	_, cmd := m.sl.Update(msg)
	return m, cmd
}

func (m menuModel) View() string {
	return m.sl.View()
}

type TypeMessage struct {
	Type        string
	Description string
}
