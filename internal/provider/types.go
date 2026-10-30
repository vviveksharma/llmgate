package provider

import "context"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream,omitempty"`
}

type ChatResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Content string `json:"content"`
}

type EmbeddingsRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type EmbeddingsResponse struct {
	Model string      `json:"model"`
	Data  [][]float32 `json:"data"`
}

type ChatStreamer interface {
	StreamChat(ctx context.Context, req ChatRequest, ch chan<- string) error
}

type Provider interface {
	Name() string
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
	Embeddings(ctx context.Context, req EmbeddingsRequest) (*EmbeddingsResponse, error)
	ChatStreamer
}
