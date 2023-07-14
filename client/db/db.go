package db

import (
	"context"

	"github.com/WiggidyW/weve-esi/client/db/sqlite_db"
)

type Db interface {
	GetRegionId(
		ctx context.Context,
		system_id uint64,
	) (int, error)
	GetStationSystemId(
		ctx context.Context,
		station_id uint64,
	) (uint32, error)
	GetSystemRegionId(
		ctx context.Context,
		system_id uint32,
	) (uint32, error)
}

func NewDb() Db {
	return sqlite_db.NewDb()
}
