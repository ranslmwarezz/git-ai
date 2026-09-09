package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultStoreUsesUserConfigDirectory(t *testing.T) {
	store, err := DefaultStore()
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	expected := filepath.Join(homeDir, ".config", "git-ai", "config")
	if store.Path() != expected {
		t.Fatalf("caminho esperado %q, obtido %q", expected, store.Path())
	}
}

func TestSaveAndLoadAPIKeyCreatesRestrictedConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "git-ai", "config")
	store := NewStore(path)

	if err := store.SaveAPIKey("test-api-key"); err != nil {
		t.Fatalf("erro ao salvar API Key: %v", err)
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("erro ao ler arquivo: %v", err)
	}
	if string(contents) != "GEMINI_API_KEY=test-api-key\n" {
		t.Fatalf("conteúdo inesperado: %q", contents)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("erro ao consultar arquivo: %v", err)
	}
	if permissions := info.Mode().Perm(); permissions != 0600 {
		t.Fatalf("permissões esperadas 0600, obtidas %04o", permissions)
	}

	directoryInfo, err := os.Stat(filepath.Dir(path))
	if err != nil {
		t.Fatalf("erro ao consultar diretório: %v", err)
	}
	if permissions := directoryInfo.Mode().Perm(); permissions != 0700 {
		t.Fatalf("permissões do diretório esperadas 0700, obtidas %04o", permissions)
	}

	apiKey, err := store.LoadAPIKey()
	if err != nil {
		t.Fatalf("erro ao carregar API Key: %v", err)
	}
	if apiKey != "test-api-key" {
		t.Fatalf("API Key esperada %q, obtida %q", "test-api-key", apiKey)
	}
}

func TestLoadAPIKeyWhenFileDoesNotExist(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "config"))

	_, err := store.LoadAPIKey()
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("esperava os.ErrNotExist, obtive %v", err)
	}
}

func TestResolveAPIKeyPrioritizesEnvironment(t *testing.T) {
	t.Setenv(apiKeyName, "environment-key")
	store := NewStore(filepath.Join(t.TempDir(), "config"))

	apiKey, err := store.ResolveAPIKey(filepath.Join(t.TempDir(), ".env"))
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if apiKey != "environment-key" {
		t.Fatalf("API Key esperada do ambiente, obtida %q", apiKey)
	}
}

func TestResolveAPIKeyUsesDevelopmentEnvBeforeGlobalConfig(t *testing.T) {
	t.Setenv(apiKeyName, "")
	directory := t.TempDir()
	envFile := filepath.Join(directory, ".env")
	if err := os.WriteFile(envFile, []byte("GEMINI_API_KEY=development-key\n"), 0600); err != nil {
		t.Fatalf("erro ao criar .env: %v", err)
	}

	store := NewStore(filepath.Join(directory, "global-config"))
	if err := store.SaveAPIKey("global-key"); err != nil {
		t.Fatalf("erro ao salvar configuração global: %v", err)
	}

	apiKey, err := store.ResolveAPIKey(envFile)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if apiKey != "development-key" {
		t.Fatalf("API Key esperada do .env, obtida %q", apiKey)
	}
}

func TestResolveAPIKeyFallsBackToGlobalConfig(t *testing.T) {
	t.Setenv(apiKeyName, "")
	directory := t.TempDir()
	store := NewStore(filepath.Join(directory, "global-config"))
	if err := store.SaveAPIKey("global-key"); err != nil {
		t.Fatalf("erro ao salvar configuração global: %v", err)
	}

	apiKey, err := store.ResolveAPIKey(filepath.Join(directory, ".env-inexistente"))
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if apiKey != "global-key" {
		t.Fatalf("API Key global esperada, obtida %q", apiKey)
	}
}

func TestResolveAPIKeyReturnsHelpfulErrorWhenMissing(t *testing.T) {
	t.Setenv(apiKeyName, "")
	store := NewStore(filepath.Join(t.TempDir(), "config"))

	_, err := store.ResolveAPIKey(filepath.Join(t.TempDir(), ".env-inexistente"))
	if !errors.Is(err, ErrAPIKeyNotConfigured) {
		t.Fatalf("erro esperado ErrAPIKeyNotConfigured, obtido %v", err)
	}
	if !strings.Contains(err.Error(), "git-ai config") {
		t.Fatalf("erro não orienta a configuração: %v", err)
	}
	if strings.Contains(err.Error(), "key") {
		t.Fatalf("erro não deveria expor credencial: %v", err)
	}
}
