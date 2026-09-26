package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisCache{client: client}
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, bool) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", false
	}
	return val, true
}

func (c *RedisCache) Set(ctx context.Context, key string, value string) {
	c.client.Set(ctx, key, value, 0)
}

func (c *RedisCache) Delete(ctx context.Context, key string) {
	c.client.Del(ctx, key)
}