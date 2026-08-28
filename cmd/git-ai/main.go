package main

import (
	"fmt"
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
		fmt.Println("Comando commit executado!")
	default: 
		fmt.Printf("Comando desconhecido: %s\n", command)
		
	}
}