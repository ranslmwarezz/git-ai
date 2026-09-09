package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type APIKeySaver interface {
	SaveAPIKey(string) error
}

type configViewState int

const (
	configInputState configViewState = iota
	configSavedState
	configErrorState
)

type ConfigModel struct {
	input textinput.Model
	saver APIKeySaver
	state configViewState
	err   error
	width int
}

func NewConfigModel(saver APIKeySaver) ConfigModel {
	input := textinput.New()
	input.Prompt = "> "
	input.EchoMode = textinput.EchoPassword
	input.EchoCharacter = '*'
	input.CharLimit = 0
	input.Focus()

	return ConfigModel{input: input, saver: saver}
}

func (m ConfigModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ConfigModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.input.Width = maxConfigInputWidth(msg.Width)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.state == configInputState {
				return m, tea.Quit
			}
			if m.state == configErrorState {
				return m, tea.Quit
			}
		case "q":
			if m.state != configInputState {
				return m, tea.Quit
			}
		case "enter":
			if m.state == configInputState {
				return m.save()
			}
			if m.state == configErrorState {
				m.err = nil
				m.state = configInputState
				m.input.Focus()
				return m, textinput.Blink
			}
			return m, tea.Quit
		}

		if m.state == configInputState {
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m ConfigModel) save() (tea.Model, tea.Cmd) {
	apiKey := strings.TrimSpace(m.input.Value())
	if apiKey == "" {
		m.err = fmt.Errorf("a API Key não pode estar vazia")
		m.state = configErrorState
		return m, nil
	}
	if m.saver == nil {
		m.err = fmt.Errorf("armazenamento de configuração não disponível")
		m.state = configErrorState
		return m, nil
	}
	if err := m.saver.SaveAPIKey(apiKey); err != nil {
		m.err = err
		m.state = configErrorState
		return m, nil
	}

	m.input.Blur()
	m.state = configSavedState
	return m, nil
}

func (m ConfigModel) View() string {
	title := titleStyle.Render("Git-AI")
	subtitle := subtitleStyle.Render("Configuração")

	switch m.state {
	case configSavedState:
		return fmt.Sprintf("%s\n%s\n\n%s\n\n%s", title, subtitle,
			successStyle.Render("API Key configurada com sucesso."), footerStyle.Render("Pressione Enter para sair."))
	case configErrorState:
		return fmt.Sprintf("%s\n%s\n\n%s\n\n%s", title, subtitle,
			errorStyle.Render(m.err.Error()), footerStyle.Render("Enter tentar novamente   Esc sair"))
	default:
		return fmt.Sprintf("%s\n%s\n\n%s\n\n%s\n\n%s", title, subtitle,
			messageLabelStyle.Render("Gemini API Key:"), m.input.View(), footerStyle.Render("Enter confirmar   Esc sair"))
	}
}

func maxConfigInputWidth(windowWidth int) int {
	const horizontalPadding = 4
	const minimumWidth = 24

	if windowWidth-horizontalPadding < minimumWidth {
		return minimumWidth
	}
	return windowWidth - horizontalPadding
}
