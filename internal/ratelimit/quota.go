package ratelimit

import "context"

type Quota struct {
	MaxRequests int64
	Used        int64
}

type QuotaChecker interface {
	Check(ctx context.Context, tenantID string) (allowed bool, q Quota, err error)
	Record(ctx context.Context, tenantID string, n int64) error
}

type MemoryQuota struct {
	limits map[string]int64
	used   map[string]int64
}

func NewMemoryQuota(defaultLimit int64) *MemoryQuota {
	return &MemoryQuota{
		limits: make(map[string]int64),
		used:   make(map[string]int64),
	}
}

func (m *MemoryQuota) Check(ctx context.Context, tenantID string) (bool, Quota, error) {
	limit := m.limits[tenantID]
	if limit == 0 {
		limit = 10000
	}
	used := m.used[tenantID]
	return used < limit, Quota{MaxRequests: limit, Used: used}, nil
}

func (m *MemoryQuota) Record(ctx context.Context, tenantID string, n int64) error {
	m.used[tenantID] += n
	return nil
}
