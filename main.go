package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type mainModel struct {
	activeModel tea.Model
}

type DiscordVersion struct {
	Name        string
	Description string
	Index       int
}

// Available Discord versions
var discordVersions = []interface{}{
	DiscordVersion{"Discord", "Standard Discord release", 0},
	DiscordVersion{"Discord PTB", "Public Test Build - Beta features", 1},
	DiscordVersion{"Discord Canary", "Experimental bleeding-edge build", 2},
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen, m.activeModel.Init())
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case switchToInstallModel:
		m.activeModel = initInstallModel()
		return m, tea.Batch(tea.ClearScreen, m.activeModel.Init())

	case switchToUninstallModel:
		m.activeModel = initUninstallModel()
		return m, tea.Batch(tea.ClearScreen, m.activeModel.Init())

	case switchToRepairModel:
		m.activeModel = initRepairModel()
		return m, tea.Batch(tea.ClearScreen, m.activeModel.Init())

	case switchToMainMenu:
		m.activeModel = initMenuModel()
		return m, tea.Batch(tea.ClearScreen, m.activeModel.Init())
	}

	var cmd tea.Cmd
	m.activeModel, cmd = m.activeModel.Update(msg)
	return m, cmd
}

func (m mainModel) View() string {
	return m.activeModel.View()
}

func main() {
	m := mainModel{
		activeModel: initMenuModel(),
	}

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
