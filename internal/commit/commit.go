package commit

import (
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

func (s *Service) Run() (string, error) {
	diff, err := s.git.DiffCached()
	if err != nil {
		return "", err
	}

	return diff, nil
}
