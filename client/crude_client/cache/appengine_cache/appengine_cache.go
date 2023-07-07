package appengine_cache

import (
	"github.com/WiggidyW/weve-esi/client/crude_client/response"

	"context"
	"sync"

	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
)

const HEAD_PREFIX = "head/"

type Cache struct {
	locks      map[string]*sync.Mutex
	head_locks map[string]*sync.Mutex
}

func NewCache() *Cache {
	return &Cache{
		locks:      make(map[string]*sync.Mutex),
		head_locks: make(map[string]*sync.Mutex),
	}
}

func lock(m map[string]*sync.Mutex, key string) {
	if _, ok := m[key]; !ok {
		m[key] = &sync.Mutex{}
	}
	m[key].Lock()
}

func unlock(m map[string]*sync.Mutex, key string) {
	m[key].Unlock()
}

func get[T interface{}](
	ctx context.Context,
	key string,
) (*T, error) {
	v := new(T)
	appengineCtx := appengine.BackgroundContext()
	_, err := memcache.Gob.Get(appengineCtx, key, v)
	if err == nil {
		return v, nil
	} else if err == memcache.ErrCacheMiss {
		return nil, nil
	} else {
		return nil, err
	}
}

func set[T interface{}](
	ctx context.Context,
	key string,
	v *T,
) error {
	appengineCtx := appengine.BackgroundContext()
	return memcache.Gob.Set(appengineCtx, &memcache.Item{
		Key:    key,
		Object: v,
	})
}

func (c *Cache) Lock(key string) {
	lock(c.locks, key)
}

func (c *Cache) Unlock(key string) {
	unlock(c.locks, key)
}

func (c *Cache) LockHead(key string) {
	lock(c.head_locks, key)
}

func (c *Cache) UnlockHead(key string) {
	unlock(c.head_locks, key)
}

func (c *Cache) Get(
	ctx context.Context,
	key string,
) (*response.EsiResponse, error) {
	return get[response.EsiResponse](ctx, key)
}

func (c *Cache) GetHead(
	ctx context.Context,
	key string,
) (*response.EsiHeadResponse, error) {
	return get[response.EsiHeadResponse](ctx, HEAD_PREFIX+key)
}

func (c *Cache) Set(
	ctx context.Context,
	key string,
	val *response.EsiResponse,
) error {
	return set(ctx, key, val)
}

func (c *Cache) SetHead(
	ctx context.Context,
	key string,
	val *response.EsiHeadResponse,
) error {
	return set(ctx, HEAD_PREFIX+key, val)
}
