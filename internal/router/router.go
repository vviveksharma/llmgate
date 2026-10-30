package router

import (
	"context"

	"github.com/vviveksharma/llmgate/internal/provider"
)

type ModelRouter struct {
	registry *provider.Registry
}

func NewModelRouter(registry *provider.Registry) *ModelRouter {
	return &ModelRouter{registry: registry}
}

func (r *ModelRouter) Chat(ctx context.Context, req provider.ChatRequest) (*provider.ChatResponse, error) {
	name := req.Model
	if name == "" {
		name = ""
	}
	p, err := r.registry.Get(name)
	if err != nil {
		p, err = r.registry.Get("")
		if err != nil {
			return nil, err
		}
	}
	return p.Chat(ctx, req)
}

func (r *ModelRouter) Embeddings(ctx context.Context, req provider.EmbeddingsRequest) (*provider.EmbeddingsResponse, error) {
	p, err := r.registry.Get(req.Model)
	if err != nil {
		p, err = r.registry.Get("")
		if err != nil {
			return nil, err
		}
	}
	return p.Embeddings(ctx, req)
}
