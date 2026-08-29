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

func TestReturnsDiff(t *testing.T) {
	expectedDiff := "diff de teste"

	fake := &fakeGit{
		diff: expectedDiff,
	}

	service := NewService(fake)

	diff, err := service.Run()

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if diff != expectedDiff {
		t.Errorf("diff esperado %q, obtido: %q", expectedDiff, diff)
	}
}

func TestReturnsError(t *testing.T) {
	expectedError := errors.New("erro ao executar git")

	fake := &fakeGit{
		err: expectedError,
	}

	service := NewService(fake)

	_, err := service.Run()

	// se o erro for diferente do que qual espero o if é acionado
	if err != expectedError {
		t.Errorf("erro esperado %v, obtido: %v", expectedError, err)
	}
}

func TestReturnsEmptyDiff(t *testing.T) {

	fake := &fakeGit{
		diff: "",
	}

	service := NewService(fake)

	diff, err := service.Run()

	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}

	if diff != "" {
		t.Errorf("esperava diff vazio, obtido: %q", diff)
	}
}
