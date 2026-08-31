package main

import (
	"fmt"
	"git-ai/internal/ai"
	"git-ai/internal/commit"
	"git-ai/internal/git"
	"os"

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
		gitClient := &git.Client{}
		aiClient := ai.NewAPIClient("https://generativelanguage.googleapis.com/v1beta/interactions",
			apiKey)
		commitService := commit.NewService(gitClient, aiClient)

		diff, err := commitService.Run()

		if err != nil {
			fmt.Println("Erro: ", err)
			return
		}

		if diff == "" {
			fmt.Println("Nenhuma alteração encontrada no staging.")
			fmt.Println("Execute: git add <arquivo>")
			return
		}

		fmt.Println(diff)

	default:
		fmt.Printf("Comando desconhecido: %s\n", command)

	}
}
