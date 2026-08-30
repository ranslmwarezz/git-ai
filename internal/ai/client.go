package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type errorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

type APIClient struct {
	baseURL string
	apiKey  string
}

type requestBody struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type responseBody struct {
	Steps []struct {
		Type    string `json:"type"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"steps"`
}

func NewAPIClient(baseURL, apiKey string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

// garante em tempo de compilação que APIClient implementa a interface AIClient
var _ AIClient = (*APIClient)(nil)

func (c *APIClient) GenerateCommitMessage(diff string) (string, error) {
	body := requestBody{
		Model: "gemini-3.6-flash",
		Input: diff,
	}

	jsonBody, err := json.Marshal(body)

	if err != nil {
		return "", err
	}

	// preparando a requisição
	req, err := http.NewRequest(
		http.MethodPost,
		c.baseURL,
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", c.apiKey)

	// enviando a requisição
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// o body é uma stream de dados, logo precisa ser fechado
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, err := io.ReadAll(resp.Body)

		if err != nil {
			return "", fmt.Errorf("API returnou status: %d", resp.StatusCode)
		}

		var apiError errorResponse

		if err := json.Unmarshal(respBody, &apiError); err != nil {
			return "", fmt.Errorf("API returnou status: %d", resp.StatusCode)
		}
		return "", fmt.Errorf(
			"API retornou status %d: %s",
			resp.StatusCode,
			apiError.Error.Message,
		)
	}

	// lendo a resposta de sucesso
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var response responseBody

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return "", err
	}

	for _, step := range response.Steps {
		if step.Type == "model_output" {
			for _, content := range step.Content {
				if content.Type == "text" {
					return content.Text, nil
				}
			}
		}
	}

	return "", fmt.Errorf("a resposta da IA não contém texto")
}
