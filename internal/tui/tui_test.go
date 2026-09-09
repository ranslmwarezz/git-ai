package tui

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type fakeGitClient struct {
	commitMessage string
	commitError   error
	commitCalls   int
	pushError     error
	pushCalls     int
}

func (f *fakeGitClient) DiffCached() (string, error) {
	return "", nil
}

func (f *fakeGitClient) Commit(message string) error {
	f.commitCalls++
	f.commitMessage = message
	return f.commitError
}

func (f *fakeGitClient) Push() error {
	f.pushCalls++
	return f.pushError
}

func TestSelectingCommitRunsGitCommit(t *testing.T) {
	client := &fakeGitClient{}
	model := NewModel("docs: atualiza README", client)

	updated, cmd := model.Update(keyMsg("enter"))
	model = updated.(Model)
	if cmd == nil {
		t.Fatal("esperava um comando de commit")
	}

	result := cmd().(commitResultMsg)
	updated, _ = model.Update(result)
	model = updated.(Model)

	if client.commitCalls != 1 {
		t.Fatalf("esperava uma chamada de commit, obtive %d", client.commitCalls)
	}
	if client.commitMessage != "docs: atualiza README" {
		t.Fatalf("mensagem esperada %q, obtida %q", "docs: atualiza README", client.commitMessage)
	}
	if model.state != successState {
		t.Fatalf("esperava estado de sucesso, obtive %v", model.state)
	}
	if !strings.Contains(model.View(), "Commit realizado com sucesso") {
		t.Fatal("a tela de sucesso não foi exibida")
	}
}

func TestSelectingCancelDoesNotRunGitCommit(t *testing.T) {
	client := &fakeGitClient{}
	model := NewModel("feat: nova funcionalidade", client)

	updated, _ := model.Update(keyMsg("down"))
	model = updated.(Model)
	updated, _ = model.Update(keyMsg("down"))
	model = updated.(Model)
	updated, cmd := model.Update(keyMsg("enter"))
	model = updated.(Model)

	if cmd == nil {
		t.Fatal("esperava o encerramento da TUI")
	}
	if client.commitCalls != 0 {
		t.Fatal("cancelar não deveria executar commit")
	}
	if model.cursor != 2 {
		t.Fatalf("esperava a opção cancelar selecionada, obtive %d", model.cursor)
	}
}

func TestNavigatingSelectsEditOption(t *testing.T) {
	model := NewModel("feat: mensagem original", &fakeGitClient{})

	updated, _ := model.Update(keyMsg("down"))
	model = updated.(Model)
	if model.cursor != 1 {
		t.Fatalf("esperava a opção editar selecionada, obtive %d", model.cursor)
	}
	if !strings.Contains(model.View(), "❯ Editar mensagem") {
		t.Fatal("a opção editar não foi destacada")
	}
}

func TestEditingMessageAndSavingIt(t *testing.T) {
	model := NewModel("feat: mensagem original", &fakeGitClient{})

	updated, _ := model.Update(keyMsg("down"))
	model = updated.(Model)
	updated, _ = model.Update(keyMsg("enter"))
	model = updated.(Model)
	if model.state != editingState {
		t.Fatalf("esperava estado de edição, obtive %v", model.state)
	}

	model.editor.SetValue("fix: mensagem editada")
	updated, _ = model.Update(keyMsg("enter"))
	model = updated.(Model)

	if model.state != suggestionState {
		t.Fatalf("esperava voltar à sugestão, obtive %v", model.state)
	}
	if model.message != "fix: mensagem editada" {
		t.Fatalf("mensagem editada não foi salva: %q", model.message)
	}
}

func TestCancelEditingPreservesMessage(t *testing.T) {
	model := NewModel("feat: mensagem original", &fakeGitClient{})

	updated, _ := model.Update(keyMsg("down"))
	model = updated.(Model)
	updated, _ = model.Update(keyMsg("enter"))
	model = updated.(Model)
	model.editor.SetValue("fix: alteração cancelada")

	updated, _ = model.Update(keyMsg("esc"))
	model = updated.(Model)

	if model.state != suggestionState {
		t.Fatalf("esperava voltar à sugestão, obtive %v", model.state)
	}
	if model.message != "feat: mensagem original" {
		t.Fatalf("cancelar edição alterou a mensagem para %q", model.message)
	}
}

func TestEscapeAndQCancelWithoutCommit(t *testing.T) {
	for _, key := range []string{"esc", "q"} {
		t.Run(key, func(t *testing.T) {
			client := &fakeGitClient{}
			model := NewModel("fix: corrige erro", client)

			_, cmd := model.Update(keyMsg(key))
			if cmd == nil {
				t.Fatal("esperava o encerramento da TUI")
			}
			if client.commitCalls != 0 {
				t.Fatal("sair não deveria executar commit")
			}
		})
	}
}

func TestCommitErrorIsShown(t *testing.T) {
	expectedError := errors.New("git: pre-commit hook recusou o commit")
	client := &fakeGitClient{commitError: expectedError}
	model := NewModel("chore: ajusta configuração", client)

	updated, cmd := model.Update(keyMsg("enter"))
	model = updated.(Model)
	updated, _ = model.Update(cmd().(commitResultMsg))
	model = updated.(Model)

	if model.state != errorState {
		t.Fatalf("esperava estado de erro, obtive %v", model.state)
	}
	if !strings.Contains(model.View(), expectedError.Error()) {
		t.Fatal("a tela de erro não exibiu o erro original")
	}
}

func TestPushRequiresConfirmationAndCanSucceed(t *testing.T) {
	client := &fakeGitClient{}
	model := NewModel("docs: atualiza README", client)

	updated, cmd := model.Update(keyMsg("enter"))
	model = updated.(Model)
	updated, _ = model.Update(cmd().(commitResultMsg))
	model = updated.(Model)

	updated, _ = model.Update(keyMsg("enter"))
	model = updated.(Model)
	if model.state != pushConfirmState || client.pushCalls != 0 {
		t.Fatal("o push deveria aguardar confirmação explícita")
	}

	updated, cmd = model.Update(keyMsg("enter"))
	model = updated.(Model)
	if model.state != pushingState || cmd == nil {
		t.Fatal("esperava iniciar o push")
	}
	updated, _ = model.Update(cmd().(pushResultMsg))
	model = updated.(Model)

	if model.state != completedState || client.pushCalls != 1 {
		t.Fatal("esperava push concluído")
	}
	if !strings.Contains(model.View(), "Push realizado com sucesso") {
		t.Fatal("a tela de sucesso do push não foi exibida")
	}
}

func TestCancelPushDoesNotRunGitPush(t *testing.T) {
	client := &fakeGitClient{}
	model := NewModel("docs: atualiza README", client)

	updated, cmd := model.Update(keyMsg("enter"))
	model = updated.(Model)
	updated, _ = model.Update(cmd().(commitResultMsg))
	model = updated.(Model)
	updated, _ = model.Update(keyMsg("enter"))
	model = updated.(Model)
	updated, _ = model.Update(keyMsg("down"))
	model = updated.(Model)
	updated, _ = model.Update(keyMsg("enter"))
	model = updated.(Model)

	if model.state != postCommitState || client.pushCalls != 0 {
		t.Fatal("cancelar push não deveria executar git push")
	}
}

func TestPushErrorIsShown(t *testing.T) {
	expectedError := errors.New("fatal: não foi possível acessar o remoto")
	client := &fakeGitClient{pushError: expectedError}
	model := NewModel("docs: atualiza README", client)

	updated, cmd := model.Update(keyMsg("enter"))
	model = updated.(Model)
	updated, _ = model.Update(cmd().(commitResultMsg))
	model = updated.(Model)
	updated, _ = model.Update(keyMsg("enter"))
	model = updated.(Model)
	updated, cmd = model.Update(keyMsg("enter"))
	model = updated.(Model)
	updated, _ = model.Update(cmd().(pushResultMsg))
	model = updated.(Model)

	if model.state != errorState || !strings.Contains(model.View(), expectedError.Error()) {
		t.Fatal("a tela de erro do push não preservou o erro original")
	}
}

func keyMsg(key string) tea.KeyMsg {
	switch key {
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEscape}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
	}
}
