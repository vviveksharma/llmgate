package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sync"
)

type SemanticCache struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewSemanticCache() *SemanticCache {
	return &SemanticCache{data: make(map[string][]byte)}
}

func (c *SemanticCache) Get(ctx context.Context, key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	return v, ok
}

func (c *SemanticCache) Set(ctx context.Context, key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func HashKey(parts ...string) string {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(p))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}
