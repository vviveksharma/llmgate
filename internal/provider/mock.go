package provider

import (
	"context"
	"fmt"
)

type Mock struct{}

func NewMock() *Mock { return &Mock{} }

func (m *Mock) Name() string { return "mock" }

func (m *Mock) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	return &ChatResponse{
		ID:      "mock-1",
		Model:   req.Model,
		Content: "mock response",
	}, nil
}

func (m *Mock) Embeddings(ctx context.Context, req EmbeddingsRequest) (*EmbeddingsResponse, error) {
	out := make([][]float32, len(req.Input))
	for i := range req.Input {
		out[i] = []float32{0.1, 0.2, 0.3}
	}
	return &EmbeddingsResponse{Model: req.Model, Data: out}, nil
}

func (m *Mock) StreamChat(ctx context.Context, req ChatRequest, ch chan<- string) error {
	defer close(ch)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case ch <- fmt.Sprintf("mock stream for %s", req.Model):
	}
	return nil
}
