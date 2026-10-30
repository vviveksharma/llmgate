package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Ollama struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewOllama(baseURL string) *Ollama {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	return &Ollama{BaseURL: baseURL, HTTPClient: http.DefaultClient}
}

func (o *Ollama) Name() string { return "ollama" }

func (o *Ollama) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	body, _ := json.Marshal(map[string]any{
		"model":    req.Model,
		"messages": req.Messages,
		"stream":   false,
	})
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, o.BaseURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := o.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("ollama chat: status %d", resp.StatusCode)
	}
	var parsed struct {
		Model   string `json:"model"`
		Message Message `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}
	return &ChatResponse{ID: "ollama", Model: parsed.Model, Content: parsed.Message.Content}, nil
}

func (o *Ollama) Embeddings(ctx context.Context, req EmbeddingsRequest) (*EmbeddingsResponse, error) {
	return nil, fmt.Errorf("ollama embeddings: not implemented")
}

func (o *Ollama) StreamChat(ctx context.Context, req ChatRequest, ch chan<- string) error {
	defer close(ch)
	return fmt.Errorf("ollama stream: not implemented")
}
