package tui

import (
	"fmt"
	"git-ai/internal/git"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type viewState int

const (
	suggestionState viewState = iota
	editingState
	committingState
	postCommitState
	pushConfirmState
	pushingState
	completedState
	errorState
)

const successState viewState = postCommitState

type errorStage int

const (
	commitError errorStage = iota
	pushError
)

type commitResultMsg struct {
	err error
}

type pushResultMsg struct {
	err error
}

type Model struct {
	message   string
	gitClient git.GitClient
	editor    textarea.Model
	cursor    int
	state     viewState
	err       error
	errStage  errorStage
	width     int
}

var (
	titleStyle         = lipgloss.NewStyle().Bold(true)
	subtitleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	messageLabelStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	messageStyle       = lipgloss.NewStyle().Bold(true)
	selectedStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	footerStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	successStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	errorStyle         = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("203"))
	secondaryTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

func NewModel(message string, gitClients ...git.GitClient) Model {
	var gitClient git.GitClient
	if len(gitClients) > 0 {
		gitClient = gitClients[0]
	}

	return Model{
		message:   message,
		gitClient: gitClient,
		editor:    newEditor(message),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.editor.SetWidth(maxMessageWidth(msg.Width))

	case commitResultMsg:
		if msg.err != nil {
			m.err = msg.err
			m.errStage = commitError
			m.state = errorState
		} else {
			m.cursor = 0
			m.state = postCommitState
		}

	case pushResultMsg:
		if msg.err != nil {
			m.err = msg.err
			m.errStage = pushError
			m.state = errorState
		} else {
			m.state = completedState
		}

	case tea.KeyMsg:
		return m.updateKey(msg)
	}

	return m, nil
}

func (m Model) updateKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	if m.state == editingState {
		switch key {
		case "enter":
			m.message = m.editor.Value()
			m.state = suggestionState
			m.cursor = 0
			return m, nil
		case "esc":
			m.state = suggestionState
			m.cursor = 0
			return m, nil
		}

		var cmd tea.Cmd
		m.editor, cmd = m.editor.Update(msg)
		return m, cmd
	}

	if m.state == committingState || m.state == pushingState {
		return m, nil
	}

	switch m.state {
	case suggestionState:
		return m.updateSuggestion(key)
	case postCommitState:
		return m.updatePostCommit(key)
	case pushConfirmState:
		return m.updatePushConfirmation(key)
	case errorState:
		if key == "enter" {
			if m.errStage == pushError {
				m.state = postCommitState
				m.cursor = 0
			} else {
				m.state = suggestionState
				m.cursor = 0
			}
			m.err = nil
			return m, nil
		}
	}

	if key == "esc" || key == "q" || (m.state == completedState && key == "enter") {
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) updateSuggestion(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "up":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down":
		if m.cursor < 2 {
			m.cursor++
		}
	case "enter":
		switch m.cursor {
		case 0:
			m.state = committingState
			return m, m.commitCmd()
		case 1:
			m.editor.SetValue(m.message)
			m.editor.Focus()
			m.state = editingState
		case 2:
			return m, tea.Quit
		}
	case "esc", "q":
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) updatePostCommit(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "up":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down":
		if m.cursor < 1 {
			m.cursor++
		}
	case "enter":
		if m.cursor == 0 {
			m.cursor = 0
			m.state = pushConfirmState
			return m, nil
		}
		return m, tea.Quit
	case "esc", "q":
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) updatePushConfirmation(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "up", "down":
		m.cursor = 1 - m.cursor
	case "enter":
		if m.cursor == 0 {
			m.state = pushingState
			return m, m.pushCmd()
		}
		m.state = postCommitState
		m.cursor = 0
	case "esc", "q":
		m.state = postCommitState
		m.cursor = 0
	}
	return m, nil
}

func (m Model) View() string {
	title := titleStyle.Render("Git-AI")
	subtitle := subtitleStyle.Render("Gerador de mensagens de commit com IA")

	switch m.state {
	case editingState:
		return fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s", title,
			messageLabelStyle.Render("Editar mensagem"), m.editor.View(),
			footerStyle.Render("Enter salvar   Esc cancelar"))
	case committingState:
		return fmt.Sprintf("%s\n%s\n\n%s", title, subtitle, secondaryTextStyle.Render("Realizando commit..."))
	case pushingState:
		return fmt.Sprintf("%s\n\n%s", title, secondaryTextStyle.Render("Realizando push..."))
	case postCommitState:
		return fmt.Sprintf("%s\n\n%s\n\n%s\n%s\n\n%s", title,
			successStyle.Render("✓ Commit realizado com sucesso"), messageStyle.Render(m.message),
			option("Fazer push", m.cursor == 0), option("Sair", m.cursor == 1))
	case pushConfirmState:
		return fmt.Sprintf("%s\n\n%s\n\n%s\n%s\n\n%s", title,
			"Deseja fazer push para o repositório remoto?", option("Sim", m.cursor == 0), option("Não", m.cursor == 1),
			footerStyle.Render("Enter confirmar   Esc cancelar"))
	case completedState:
		return fmt.Sprintf("%s\n\n%s\n\n%s", title,
			successStyle.Render("✓ Push realizado com sucesso"), footerStyle.Render("Pressione Enter para sair."))
	case errorState:
		failure := "✗ Falha ao realizar commit"
		if m.errStage == pushError {
			failure = "✗ Falha ao realizar push"
		}
		return fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s", title, errorStyle.Render(failure), m.err,
			footerStyle.Render("Enter voltar   Esc sair"))
	}

	message := messageStyle.Render(m.message)
	if m.width > 0 {
		message = lipgloss.NewStyle().Width(maxMessageWidth(m.width)).Render(message)
	}

	return fmt.Sprintf("%s\n%s\n\n%s\n\n%s\n%s\n%s\n%s\n\n%s", title, subtitle,
		messageLabelStyle.Render("Mensagem sugerida"), message,
		option("Realizar commit", m.cursor == 0), option("Editar mensagem", m.cursor == 1), option("Cancelar", m.cursor == 2),
		footerStyle.Render("↑/↓ navegar   Enter selecionar   Esc/q sair"))
}

func option(label string, selected bool) string {
	if selected {
		return selectedStyle.Render("❯ " + label)
	}
	return "  " + label
}

func (m Model) commitCmd() tea.Cmd {
	return func() tea.Msg {
		if m.gitClient == nil {
			return commitResultMsg{err: fmt.Errorf("cliente Git não configurado")}
		}
		return commitResultMsg{err: m.gitClient.Commit(m.message)}
	}
}

func (m Model) pushCmd() tea.Cmd {
	return func() tea.Msg {
		if m.gitClient == nil {
			return pushResultMsg{err: fmt.Errorf("cliente Git não configurado")}
		}
		return pushResultMsg{err: m.gitClient.Push()}
	}
}

func newEditor(message string) textarea.Model {
	editor := textarea.New()
	editor.SetValue(message)
	editor.Prompt = ""
	editor.CharLimit = 0
	editor.ShowLineNumbers = false
	editor.Blur()
	return editor
}

func maxMessageWidth(windowWidth int) int {
	const horizontalPadding = 4
	const minimumWidth = 20

	if windowWidth-horizontalPadding < minimumWidth {
		return minimumWidth
	}
	return windowWidth - horizontalPadding
}
