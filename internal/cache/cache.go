package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	rdb *redis.Client
	ctx context.Context
}

func New(addr, pass string, db int) *Cache {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})
	return &Cache{rdb: rdb, ctx: ctx}
}

func (c *Cache) Get(key string) (string, error) {
	return c.rdb.Get(c.ctx, key).Result()
}

func (c *Cache) Set(key, val string, ttl time.Duration) error {
	return c.rdb.Set(c.ctx, key, val, ttl).Err()
}
