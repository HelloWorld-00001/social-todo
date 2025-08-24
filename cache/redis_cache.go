package cache

import (
	"context"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/coderconquerer/social-todo/common"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type redisCache struct {
	store *cache.Cache
}

func NewRedisCache(sc goservice.ServiceContext) *redisCache {
	rdClient := sc.MustGet(common.PluginRedis).(*redis.Client)

	c := cache.New(&cache.Options{
		Redis:      rdClient,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return &redisCache{store: c}
}

func (rdc *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return rdc.store.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   ttl,
	})
}

func (rdc *redisCache) Get(ctx context.Context, key string, value interface{}) error {
	return rdc.store.Get(ctx, key, value)
}

func (rdc *redisCache) Delete(ctx context.Context, key string) error {
	return rdc.store.Delete(ctx, key)
}
