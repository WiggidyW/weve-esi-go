package local_cache

import (
	"github.com/WiggidyW/weve-esi/client/crude_client/response"

	"context"
	"sync"
)

type Cache struct {
	head_reps  map[string]*response.EsiHeadResponse
	reps       map[string]*response.EsiResponse
	locks      map[string]*sync.Mutex
	head_locks map[string]*sync.Mutex
}

func NewCache() *Cache {
	return &Cache{
		head_reps:  make(map[string]*response.EsiHeadResponse),
		reps:       make(map[string]*response.EsiResponse),
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
	return c.reps[key], nil
}

func (c *Cache) GetHead(
	ctx context.Context,
	key string,
) (*response.EsiHeadResponse, error) {
	return c.head_reps[key], nil
}

func (c *Cache) Set(
	ctx context.Context,
	key string,
	val *response.EsiResponse,
) error {
	c.reps[key] = val
	return nil
}

func (c *Cache) SetHead(
	ctx context.Context,
	key string,
	val *response.EsiHeadResponse,
) error {
	c.head_reps[key] = val
	return nil
}
