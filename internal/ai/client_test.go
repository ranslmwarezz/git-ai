package ai

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGeneratedCommitMessage(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			t.Errorf("método esperado POST, obtido: %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type esperado application/json, obtido: %s",
				r.Header.Get("Content-Type"))
		}

		if r.Header.Get("x-goog-api-key") != "fake-api-key" {
			t.Errorf("API key incorreta")
		}

		body, err := io.ReadAll(r.Body)

		if err != nil {
			t.Errorf("erro ao ler body: %v", err)
		}

		if !strings.Contains(string(body), "Retorne somente a mensagem de commit") {
			t.Fatal("prompt não contém a mensagem esperada")
		}

		var request requestBody

		err = json.Unmarshal(body, &request)
		if err != nil {
			t.Errorf("erro ao decodificar body: %v", err)
		}

		if request.Model != "gemini-3.6-flash" {
			t.Errorf("modelo esperado gemini-3.6-flash, obtido: %s", request.Model)
		}

		if !strings.Contains(request.Input, "diff de teste") {
			t.Fatal("input não contém o diff esperado")
		}

		w.WriteHeader(http.StatusOK)

		w.Write([]byte(`{
		"steps": [
		{
		"type": "model_output",
		"content": [
		{ "type": "text",
		 "text": "feat: adiciona autenticação JWT"
		}] 
		}
		]
		}`))

	}))

	defer server.Close()

	client := NewAPIClient(server.URL, "fake-api-key")

	result, err := client.GenerateCommitMessage("diff de teste")

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	expected := "feat: adiciona autenticação JWT"

	if result != expected {
		t.Fatalf("esperado %q, obtido: %q", expected, result)
	}
}

func TestGenerateCommitMessageAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		w.Write([]byte(`{
		"error": {
		"code": 500,
		"message": "erro interno de teste",
		"status": "INTERNAL"}}`))
	}))

	defer server.Close()

	client := NewAPIClient(server.URL, "fake-api-key")

	_, err := client.GenerateCommitMessage("diff de teste")

	if err == nil {
		t.Fatal("esperava um erro, mas não ocorreu")
	}

	expected := "API retornou status 500: erro interno de teste"

	if err.Error() != expected {
		t.Fatalf("erro esperado %q, obtido: %q", expected, err.Error())
	}

}

func TestGenerateCommitMessageInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(`isso não é um JSON válido`))
	}))

	defer server.Close()

	client := NewAPIClient(server.URL, "fake-api-key")

	_, err := client.GenerateCommitMessage("diff de teste")

	if err == nil {
		t.Fatal("esperava um erro ao receber o JSON inválido")
	}
}

func TestGenerateCommitMessageWithoutText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(`{
		"steps": []}`))
	}))

	defer server.Close()

	client := NewAPIClient(server.URL, "fake-api-key")

	_, err := client.GenerateCommitMessage("diff de teste")

	if err == nil {
		t.Fatal("esperava um erro quando a resposta não contém text")
	}
}
