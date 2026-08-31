package commit

import (
	"git-ai/internal/ai"
	"git-ai/internal/git"
)

type Service struct {
	git git.GitClient
	ai ai.AIClient
}

func NewService(gitClient git.GitClient, aiClient ai.AIClient) *Service {
	return &Service{
		git: gitClient,
		ai: aiClient,
	}
}

func (s *Service) Run() (string, error) {
	diff, err := s.git.DiffCached()
	if err != nil {
		return "", err
	}

	if diff == "" {
		return "", nil
	}

	message, err := s.ai.GenerateCommitMessage(diff)

	if err != nil {
		return "", err
	}

	return message, nil
}
