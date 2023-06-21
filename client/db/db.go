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
}

func NewDb() Db {
	return sqlite_db.NewDb()
}
