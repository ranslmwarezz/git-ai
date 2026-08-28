package main

import (
	"fmt"
	"git-ai/internal/commit"
	"git-ai/internal/git"
	"os"
)

func main() {

	args := os.Args

	if len(args) < 2 {
		fmt.Println("Uso: git-ai <comando>")
		return
	}

	command := args[1]

	switch command {
	case "commit":
		gitClient := &git.Client{}
		commitService := commit.NewService(gitClient)

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
