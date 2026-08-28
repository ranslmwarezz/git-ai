package commit

import (
	"fmt"
	"git-ai/internal/git"
)

type Service struct {
	git git.GitClient
}

func NewService(gitClient git.GitClient) *Service {
	return &Service{
		git: gitClient,
	}
}

func (s *Service) Run() error {
	diff, err := s.git.DiffCached()
	if err != nil {
		return err
	}

	if diff == "" {
		fmt.Println("Nenhuma alteração encontrada no staging.")
		fmt.Println("Execute: git add <arquivo>")
		return nil
	}

	fmt.Println(diff)
	return nil
}
