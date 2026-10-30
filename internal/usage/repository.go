package usage

import (
	"context"
	"sync"
)

type Repository interface {
	Insert(ctx context.Context, r Record) error
}

type MemoryRepository struct {
	mu      sync.Mutex
	records []Record
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{}
}

func (m *MemoryRepository) Insert(ctx context.Context, r Record) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.records = append(m.records, r)
	return nil
}

type Recorder struct {
	repo Repository
}

func NewRecorder(repo Repository) *Recorder {
	return &Recorder{repo: repo}
}

func (r *Recorder) Record(ctx context.Context, rec Record) error {
	return r.repo.Insert(ctx, rec)
}
