package local_cache

import (
	"github.com/WiggidyW/weve-esi/client/crude_client/response"

	"context"
	"sync"
)

type Cache struct {
	reps            map[string]*response.EsiResponse
	rep_lock        sync.RWMutex
	head_reps       map[string]*response.EsiHeadResponse
	head_rep_lock   sync.RWMutex
	locks           map[string]*sync.Mutex
	locks_lock      sync.Mutex
	head_locks      map[string]*sync.Mutex
	head_locks_lock sync.Mutex
}

func NewCache() *Cache {
	return &Cache{
		reps:            make(map[string]*response.EsiResponse),
		rep_lock:        sync.RWMutex{},
		head_reps:       make(map[string]*response.EsiHeadResponse),
		head_rep_lock:   sync.RWMutex{},
		locks:           make(map[string]*sync.Mutex),
		locks_lock:      sync.Mutex{},
		head_locks:      make(map[string]*sync.Mutex),
		head_locks_lock: sync.Mutex{},
	}
}

func lock(l *sync.Mutex, m map[string]*sync.Mutex, key string) {
	l.Lock()
	defer l.Unlock()
	if _, ok := m[key]; !ok {
		m[key] = &sync.Mutex{}
	}
	m[key].Lock()
}

func unlock(m map[string]*sync.Mutex, key string) {
	m[key].Unlock()
}

func (c *Cache) Lock(key string) {
	lock(&c.locks_lock, c.locks, key)
}

func (c *Cache) Unlock(key string) {
	unlock(c.locks, key)
}

func (c *Cache) LockHead(key string) {
	lock(&c.head_locks_lock, c.head_locks, key)
}

func (c *Cache) UnlockHead(key string) {
	unlock(c.head_locks, key)
}

func (c *Cache) Get(
	ctx context.Context,
	key string,
) (*response.EsiResponse, error) {
	c.rep_lock.RLock()
	defer c.rep_lock.RUnlock()
	return c.reps[key], nil
}

func (c *Cache) GetHead(
	ctx context.Context,
	key string,
) (*response.EsiHeadResponse, error) {
	c.head_rep_lock.RLock()
	defer c.head_rep_lock.RUnlock()
	return c.head_reps[key], nil
}

func (c *Cache) Set(
	ctx context.Context,
	key string,
	val *response.EsiResponse,
) error {
	c.rep_lock.Lock()
	defer c.rep_lock.Unlock()
	c.reps[key] = val
	return nil
}

func (c *Cache) SetHead(
	ctx context.Context,
	key string,
	val *response.EsiHeadResponse,
) error {
	c.head_rep_lock.Lock()
	defer c.head_rep_lock.Unlock()
	c.head_reps[key] = val
	return nil
}
