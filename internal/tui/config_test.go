package tui

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type fakeAPIKeySaver struct {
	apiKey string
	err    error
}

func (f *fakeAPIKeySaver) SaveAPIKey(apiKey string) error {
	f.apiKey = apiKey
	return f.err
}

func TestConfigModelSavesMaskedAPIKey(t *testing.T) {
	saver := &fakeAPIKeySaver{}
	model := NewConfigModel(saver)
	model.input.SetValue("secret-api-key")

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(ConfigModel)

	if model.state != configSavedState {
		t.Fatalf("esperava estado salvo, obtive %v", model.state)
	}
	if saver.apiKey != "secret-api-key" {
		t.Fatalf("API Key esperada %q, obtida %q", "secret-api-key", saver.apiKey)
	}
	if strings.Contains(model.View(), "secret-api-key") {
		t.Fatal("a tela de sucesso não deveria exibir a API Key")
	}
}

func TestConfigModelMasksAPIKeyInput(t *testing.T) {
	model := NewConfigModel(&fakeAPIKeySaver{})
	model.input.SetValue("secret-api-key")

	if strings.Contains(model.input.View(), "secret-api-key") {
		t.Fatal("a entrada não deveria exibir a API Key em texto claro")
	}
}

func TestConfigModelCancelDoesNotSave(t *testing.T) {
	saver := &fakeAPIKeySaver{}
	model := NewConfigModel(saver)
	model.input.SetValue("secret-api-key")

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEscape})
	if cmd == nil {
		t.Fatal("esperava o encerramento da TUI")
	}
	if saver.apiKey != "" {
		t.Fatal("cancelar não deveria salvar a API Key")
	}
}

func TestConfigModelShowsSaveError(t *testing.T) {
	expectedError := errors.New("erro de gravação")
	model := NewConfigModel(&fakeAPIKeySaver{err: expectedError})
	model.input.SetValue("secret-api-key")

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(ConfigModel)

	if model.state != configErrorState {
		t.Fatalf("esperava estado de erro, obtive %v", model.state)
	}
	if !strings.Contains(model.View(), expectedError.Error()) {
		t.Fatal("o erro de configuração não foi exibido")
	}
}
