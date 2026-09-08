package main

import (
	"fmt"
	"git-ai/internal/ai"
	"git-ai/internal/commit"
	"git-ai/internal/git"
	"git-ai/internal/tui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		fmt.Println("Aviso: arquivo .env não encontrado")
	}

	args := os.Args

	if len(args) < 2 {
		fmt.Println("Uso: git-ai <comando>")
		return
	}

	command := args[1]

	apiKey := os.Getenv("GEMINI_API_KEY")

	if apiKey == "" {
		fmt.Println("Erro: GEMINI_API_KEY não configurada")
		fmt.Println("Defina a variável de ambiente ou crie um arquivo .env")
		return
	}

	switch command {
	case "commit":
		gitClient := git.NewClient()
		aiClient := ai.NewAPIClient("https://generativelanguage.googleapis.com/v1beta/interactions",
			apiKey)
		commitService := commit.NewService(gitClient, aiClient)

		message, err := commitService.Run()

		if err != nil {
			fmt.Println("Erro: ", err)
			return
		}

		if message == "" {
			fmt.Println("Nenhuma alteração encontrada no staging.")
			fmt.Println("Execute: git add <arquivo>")
			return
		}

		model := tui.NewModel(message)

		program := tea.NewProgram(model)

		finalModel, err := program.Run()
		if err != nil {
			fmt.Println("Erro ao inicar TUI:", err)
			return
		}

		final := finalModel.(tui.Model)

		if !final.Confirmed() {
			fmt.Println("Commit cancelado.")
			return
		}

		if err := gitClient.Commit(message); err != nil {
			fmt.Println("Erro ao executar commit:", err)
			return
		}
		fmt.Println("Commit realizado.")

	default:
		fmt.Printf("Comando desconhecido: %s\n", command)

	}
}
