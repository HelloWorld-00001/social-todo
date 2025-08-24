package cache

import (
	"context"
	"sync"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string, value interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type caching struct {
	store  map[string]interface{}
	locker *sync.RWMutex
}

func NewCaching() *caching {
	return &caching{
		store:  make(map[string]interface{}),
		locker: new(sync.RWMutex),
	}
}

func (c *caching) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.locker.Lock()
	defer c.locker.Unlock()

	c.store[key] = value
	return nil
}

func (c *caching) Get(ctx context.Context, key string, value interface{}) error {
	c.locker.RLock()
	defer c.locker.RUnlock()

	if v, ok := c.store[key]; ok {
		value = v
	}
	return nil
}

func (c *caching) Delete(ctx context.Context, key string) error {
	c.locker.Lock()
	defer c.locker.Unlock()
	delete(c.store, key)
	return nil
}
