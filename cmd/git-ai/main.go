package main

import (
	"fmt"
	"os"

	"github.com/ranslmwarezz/git-ai/internal/ai"
	"github.com/ranslmwarezz/git-ai/internal/commit"
	"github.com/ranslmwarezz/git-ai/internal/config"
	"github.com/ranslmwarezz/git-ai/internal/git"
	"github.com/ranslmwarezz/git-ai/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	command := ""
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	if command == "config" {
		runConfig()
		return
	}

	if command != "" && command != "commit" {
		fmt.Printf("Comando desconhecido: %s\n", command)
		return
	}

	apiKey, err := config.ResolveAPIKey(".env")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}
	runCommit(apiKey)
}

func runCommit(apiKey string) {
	gitClient := git.NewClient()
	aiClient := ai.NewAPIClient("https://generativelanguage.googleapis.com/v1beta/interactions", apiKey)
	commitService := commit.NewService(gitClient, aiClient)

	message, err := commitService.Run()
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}
	if message == "" {
		fmt.Println("Nenhuma alteração encontrada no staging.")
		fmt.Println("Execute: git add <arquivo>")
		return
	}

	if _, err := tea.NewProgram(tui.NewModel(message, gitClient)).Run(); err != nil {
		fmt.Println("Erro ao iniciar TUI:", err)
	}
}

func runConfig() {
	store, err := config.DefaultStore()
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	if _, err := tea.NewProgram(tui.NewConfigModel(store)).Run(); err != nil {
		fmt.Println("Erro ao iniciar TUI:", err)
	}
}
