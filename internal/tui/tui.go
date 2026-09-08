package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	message string 
	quitting bool
	confirmed bool
}

func NewModel(message string) Model {
	return Model{
		message: message,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {


	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.confirmed = true
			m.quitting = true
			return m, tea.Quit

		case "esc", "q":
			m.quitting = true
			return m, tea.Quit
	
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	return fmt.Sprintf("Git-AI\n\n"+"Menssagem sugerida:\n\n"+" %s\n\n"+"Enter → confirmar\n"+"Esc → cancelar\n", m.message)
}

func (m Model) Confirmed() bool {
	return m.confirmed
}