package git

import (
	"os"
	"os/exec"
	"testing"
)

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	switch os.Getenv("HELPER_MODE") {
	case "diff":
		_, _ = os.Stdout.WriteString("diff --git a/test.txt b/test.txt\n+alteração\n")
		os.Exit(0)

	case "error":
		os.Exit(1)

	case "push":
		_, _ = os.Stdout.WriteString("push realizado\n")
		os.Exit(0)

	default:
		os.Exit(1)
	}
}

func helperCommand(mode string) *exec.Cmd {
	cmd := exec.Command(
		os.Args[0],
		"-test.run=TestHelperProcess",
	)

	cmd.Env = append(
		os.Environ(),
		"GO_WANT_HELPER_PROCESS=1",
		"HELPER_MODE="+mode,
	)

	return cmd
}

func TestClient_DiffCached(t *testing.T) {
	expectedDiff := "diff --git a/test.txt b/test.txt\n+alteração\n"

	client := &Client{
		command: func(name string, args ...string) *exec.Cmd {
			if name != "git" {
				t.Fatalf("esperava comando git, mas obteve %q", name)
			}

			if len(args) != 2 {
				t.Fatalf("esperava 2 argumentos, mas obteve %d", len(args))
			}

			if args[0] != "diff" || args[1] != "--cached" {
				t.Fatalf("argumentos inesperados: %v", args)
			}

			return helperCommand("diff")
		},
	}

	diff, err := client.DiffCached()

	if err != nil {
		t.Fatalf("esperava nenhum erro, mas obteve: %v", err)
	}

	if diff != expectedDiff {
		t.Fatalf("esperava %q, mas obteve %q", expectedDiff, diff)
	}
}

func TestClient_DiffCached_Error(t *testing.T) {
	client := &Client{
		command: func(name string, args ...string) *exec.Cmd {
			return helperCommand("error")
		},
	}

	_, err := client.DiffCached()

	if err == nil {
		t.Fatal("esperava um erro, mas nenhum erro ocorreu")
	}
}

func TestClient_Commit(t *testing.T) {
	message := "feat: adiciona nova funcionalidade"

	client := &Client{
		command: func(name string, args ...string) *exec.Cmd {
			if name != "git" {
				t.Fatalf("esperava comando git, mas obteve %q", name)
			}

			expectedArgs := []string{"commit", "-m", message}

			if len(args) != len(expectedArgs) {
				t.Fatalf(
					"esperava %d argumentos, mas obteve %d",
					len(expectedArgs),
					len(args),
				)
			}

			for i := range expectedArgs {
				if args[i] != expectedArgs[i] {
					t.Fatalf(
						"esperava argumento %d como %q, mas obteve %q",
						i,
						expectedArgs[i],
						args[i],
					)
				}
			}

			return helperCommand("diff")
		},
	}

	err := client.Commit(message)

	if err != nil {
		t.Fatalf("esperava nenhum erro, mas obteve: %v", err)
	}
}

func TestClient_Commit_Error(t *testing.T) {
	client := &Client{
		command: func(name string, args ...string) *exec.Cmd {
			return helperCommand("error")
		},
	}

	err := client.Commit("feat: teste")

	if err == nil {
		t.Fatal("esperava um erro, mas nenhum erro ocorreu")
	}
}

func TestClient_Push(t *testing.T) {
	client := &Client{
		command: func(name string, args ...string) *exec.Cmd {
			if name != "git" {
				t.Fatalf("esperava comando git, mas obteve %q", name)
			}
			if len(args) != 1 || args[0] != "push" {
				t.Fatalf("argumentos inesperados: %v", args)
			}
			return helperCommand("push")
		},
	}

	if err := client.Push(); err != nil {
		t.Fatalf("esperava nenhum erro, mas obteve: %v", err)
	}
}

func TestClient_Push_Error(t *testing.T) {
	client := &Client{
		command: func(name string, args ...string) *exec.Cmd {
			return helperCommand("error")
		},
	}

	if err := client.Push(); err == nil {
		t.Fatal("esperava um erro, mas nenhum erro ocorreu")
	}
}
