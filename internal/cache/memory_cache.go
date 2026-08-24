package cache

import (
	"context"
	"sync"
)

// InMemoryCache is a temporary stand-in for Redis or another real cache.
type InMemoryCache struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]string),
	}
}

func (c *InMemoryCache) Get(ctx context.Context, key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.data[key]
	return value, ok
}

func (c *InMemoryCache) Set(ctx context.Context, key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}

func (c *InMemoryCache) Delete(ctx context.Context, key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}