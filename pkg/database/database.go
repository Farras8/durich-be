package database

import (
	"durich-be/pkg/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	_ "github.com/lib/pq"
	"database/sql"
)

func Connect(cfg config.DatabaseConfig) *bun.DB {
	dsn := cfg.GetDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	bunDB := bun.NewDB(db, pgdialect.New())
	return bunDB
}