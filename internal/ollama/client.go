package ollama

import (
	"context"
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
)

// Client wraps the official Ollama API client for use by llmgate.
type Client struct {
	api *api.Client
}

// NewClient connects using OLLAMA_HOST (default http://127.0.0.1:11434).
func NewClient() (*Client, error) {
	c, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("ollama client: %w", err)
	}
	return &Client{api: c}, nil
}

// ModelNames returns models available on the local Ollama instance.
func (c *Client) ModelNames(ctx context.Context) ([]string, error) {
	list, err := c.api.List(ctx)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(list.Models))
	for _, m := range list.Models {
		names = append(names, m.Name)
	}
	return names, nil
}

// Chat sends messages to a model and returns the assistant reply (non-streaming).
func (c *Client) Chat(ctx context.Context, model string, messages []api.Message) (string, error) {
	model = strings.TrimSpace(model)
	if model == "" {
		return "", fmt.Errorf("model is required")
	}
	if len(messages) == 0 {
		return "", fmt.Errorf("at least one message is required")
	}

	stream := false
	req := &api.ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   &stream,
	}

	var reply strings.Builder
	err := c.api.Chat(ctx, req, func(resp api.ChatResponse) error {
		reply.WriteString(resp.Message.Content)
		return nil
	})
	if err != nil {
		return "", err
	}
	return reply.String(), nil
}
