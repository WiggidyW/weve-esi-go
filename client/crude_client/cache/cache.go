package cache

import (
	"context"

	"github.com/WiggidyW/weve-esi/client/crude_client/cache/appengine_cache"
	"github.com/WiggidyW/weve-esi/client/crude_client/cache/local_cache"
	"github.com/WiggidyW/weve-esi/client/crude_client/response"
)

type Cache interface {
	Lock(key string)
	Unlock(key string)
	LockHead(key string)
	UnlockHead(key string)
	Get(
		ctx context.Context,
		key string,
	) (*response.EsiResponse, error)
	GetHead(
		ctx context.Context,
		key string,
	) (*response.EsiHeadResponse, error)
	Set(
		ctx context.Context,
		key string,
		val *response.EsiResponse,
	) error
	SetHead(
		ctx context.Context,
		key string,
		val *response.EsiHeadResponse,
	) error
}

func NewCache() Cache {
	return appengine_cache.NewCache()
}

func NewLocalCache() Cache {
	return local_cache.NewCache()
}
