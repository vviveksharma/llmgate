package provider

import (
	"context"
	"fmt"
)

type Registry struct {
	providers map[string]Provider
	defaultName string
}

func NewRegistry(defaultName string) *Registry {
	return &Registry{
		providers:   make(map[string]Provider),
		defaultName: defaultName,
	}
}

func (r *Registry) Register(p Provider) {
	r.providers[p.Name()] = p
}

func (r *Registry) Get(name string) (Provider, error) {
	if name == "" {
		name = r.defaultName
	}
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
	return p, nil
}

func (r *Registry) Default(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	p, err := r.Get("")
	if err != nil {
		return nil, err
	}
	return p.Chat(ctx, req)
}
