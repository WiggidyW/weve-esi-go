package sqlite_db

import (
	_ "github.com/mattn/go-sqlite3"

	"context"
	"database/sql"
)

const (
	DB_DRIVER = "sqlite3"
	DB_PATH   = "db.sqlite"
)

type Db struct{ *sql.DB }

func NewDb() *Db {
	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		panic(err)
	}
	return &Db{db}
}

func (d *Db) GetRegionId(
	ctx context.Context,
	system_id uint64,
) (int, error) {
	var region_id int
	err := d.QueryRowContext(
		ctx,
		"SELECT region_id FROM systems WHERE system_id = ?",
		system_id,
	).Scan(&region_id)
	if err != nil {
		return 0, err
	}
	return region_id, nil
}

func (d *Db) GetStationSystemId(
	ctx context.Context,
	station_id uint64,
) (uint32, error) {
	var system_id uint32
	err := d.QueryRowContext(
		ctx,
		"SELECT system_id FROM v2_stations WHERE station_id = ?",
		station_id,
	).Scan(&system_id)
	if err != nil {
		return 0, err
	}
	return system_id, nil
}

func (d *Db) GetSystemRegionId(
	ctx context.Context,
	system_id uint32,
) (uint32, error) {
	var region_id uint32
	err := d.QueryRowContext(
		ctx,
		"SELECT region_id FROM v2_systems WHERE system_id = ?",
		system_id,
	).Scan(&region_id)
	if err != nil {
		return 0, err
	}
	return region_id, nil
}
