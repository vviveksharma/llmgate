package tenant

import "context"

type Repository interface {
	GetByID(ctx context.Context, id string) (*Tenant, error)
	List(ctx context.Context) ([]Tenant, error)
}

type MemoryRepository struct {
	tenants map[string]Tenant
}

func NewMemoryRepository(seed []Tenant) *MemoryRepository {
	m := make(map[string]Tenant, len(seed))
	for _, t := range seed {
		m[t.ID] = t
	}
	return &MemoryRepository{tenants: m}
}

func (r *MemoryRepository) GetByID(ctx context.Context, id string) (*Tenant, error) {
	t, ok := r.tenants[id]
	if !ok {
		return nil, nil
	}
	return &t, nil
}

func (r *MemoryRepository) List(ctx context.Context) ([]Tenant, error) {
	out := make([]Tenant, 0, len(r.tenants))
	for _, t := range r.tenants {
		out = append(out, t)
	}
	return out, nil
}
