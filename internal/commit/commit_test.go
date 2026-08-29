package commit

import (
	"errors"
	"testing"
)

type fakeGit struct {
	diff string
	err  error
}

func (f *fakeGit) DiffCached() (string, error) {
	return f.diff, f.err
}

type fakeAI struct {
	message string
	err     error
	called  bool
}

func (a *fakeAI) GenerateCommitMessage(diff string) (string, error) {
	a.called = true
	return a.message, a.err
}

func TestRunReturnsCommitMessage(t *testing.T) {
	fakeGit := &fakeGit{
		diff: "diff de teste",
	}

	expectedMessage := "feat: adiciona teste"

	fakeAi := &fakeAI{
		message: expectedMessage,
	}

	service := NewService(fakeGit, fakeAi)

	message, err := service.Run()

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if message != expectedMessage {
		t.Fatalf("messagem esperada %q, obitida: %q", expectedMessage, message)
	}
}

func TestReturnsError(t *testing.T) {
	expectedError := errors.New("erro ao executar git")

	fakeGit := &fakeGit{
		err: expectedError,
	}

	fakeAi := &fakeAI{}

	service := NewService(fakeGit, fakeAi)

	_, err := service.Run()

	// se o erro for diferente do que espero, o if é acionado
	if err != expectedError {
		t.Errorf("erro esperado %v, obtido: %v", expectedError, err)
	}
}

func TestReturnsEmptyDiff(t *testing.T) {

	fakeGit := &fakeGit{
		diff: "",
	}

	fakeAi := &fakeAI{}

	service := NewService(fakeGit, fakeAi)

	diff, err := service.Run()

	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}

	if diff != "" {
		t.Errorf("esperava diff vazio, obtido: %q", diff)
	}

	if fakeAi.called {
		t.Errorf("a IA não deveria ser chamada")
	}
}

func TestReturnsAIError(t *testing.T) {
	expectedError := errors.New("erro ao gerar mensagem")

	fakeGit := &fakeGit{
		diff: "diff teste",
	}

	fakeAI := &fakeAI{
		err: expectedError,
	}

	service := NewService(fakeGit, fakeAI)

	_, err := service.Run()

	if err != expectedError {
		t.Errorf("erro esperado %v, obtido: %v", expectedError, err)
	}

	if !fakeAI.called {
		t.Errorf("a IA deveria ter sido chamada")
	}
}
