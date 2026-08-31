package ai

type AIClient interface {
	GenerateCommitMessage(diff string) (string, error)
}

