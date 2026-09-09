package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

const apiKeyName = "GEMINI_API_KEY"

var ErrAPIKeyNotConfigured = errors.New("GEMINI_API_KEY não configurada; execute git-ai config")

type Store struct {
	path string
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

func DefaultStore() (*Store, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("não foi possível descobrir o diretório do usuário: %w", err)
	}

	return NewStore(filepath.Join(homeDir, ".config", "git-ai", "config")), nil
}

func (s *Store) Path() string {
	return s.path
}

func (s *Store) SaveAPIKey(apiKey string) error {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return errors.New("a API Key não pode estar vazia")
	}

	configDir := filepath.Dir(s.path)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("não foi possível criar o diretório de configuração: %w", err)
	}
	if err := os.Chmod(configDir, 0700); err != nil {
		return fmt.Errorf("não foi possível restringir o diretório de configuração: %w", err)
	}

	file, err := os.OpenFile(s.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("não foi possível abrir a configuração: %w", err)
	}
	defer file.Close()

	if err := file.Chmod(0600); err != nil {
		return fmt.Errorf("não foi possível restringir a configuração: %w", err)
	}

	if _, err := fmt.Fprintf(file, "%s=%s\n", apiKeyName, apiKey); err != nil {
		return fmt.Errorf("não foi possível salvar a API Key: %w", err)
	}

	return nil
}

func (s *Store) LoadAPIKey() (string, error) {
	values, err := godotenv.Read(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", os.ErrNotExist
		}
		return "", fmt.Errorf("não foi possível ler a configuração: %w", err)
	}

	apiKey := strings.TrimSpace(values[apiKeyName])
	if apiKey == "" {
		return "", ErrAPIKeyNotConfigured
	}
	return apiKey, nil
}

func (s *Store) ResolveAPIKey(envFile string) (string, error) {
	if apiKey := strings.TrimSpace(os.Getenv(apiKeyName)); apiKey != "" {
		return apiKey, nil
	}

	if envFile == "" {
		envFile = ".env"
	}
	if values, err := godotenv.Read(envFile); err == nil {
		if apiKey := strings.TrimSpace(values[apiKeyName]); apiKey != "" {
			return apiKey, nil
		}
	}

	if apiKey, err := s.LoadAPIKey(); err == nil {
		return apiKey, nil
	}

	return "", ErrAPIKeyNotConfigured
}

func ResolveAPIKey(envFile string) (string, error) {
	store, err := DefaultStore()
	if err != nil {
		return "", err
	}
	return store.ResolveAPIKey(envFile)
}
