package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/samber/do"
)

var _ CacheStore = (*ImMemoryCacheStore)(nil)

type CacheStore interface {
	Set(key string, value any)
	Get(key string) (any, bool)
}

type ImMemoryCacheStore struct {
	store           *cache.Cache
	cacheExpiration time.Duration
}

func NewCacheStore(expiration, clean time.Duration) do.Provider[CacheStore] {
	return func(i *do.Injector) (CacheStore, error) {
		return NewImMemoryCacheStore(expiration, clean), nil
	}
}

func NewImMemoryCacheStore(expiration, clean time.Duration) *ImMemoryCacheStore {
	store := cache.New(expiration, clean)
	return &ImMemoryCacheStore{
		store:           store,
		cacheExpiration: 5 * time.Minute,
	}
}

func (c *ImMemoryCacheStore) Set(key string, value any) {
	c.store.Set(key, value, c.cacheExpiration)
}

func (c *ImMemoryCacheStore) Get(key string) (any, bool) {
	return c.store.Get(key)
}
