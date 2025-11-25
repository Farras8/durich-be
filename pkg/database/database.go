package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"durich-be/pkg/config"
	"durich-be/pkg/logger"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type Database struct {
	*bun.DB
}

var (
	dbInstance *Database
	once       sync.Once
)

func InitDB(c config.DatabaseConfig) {
	once.Do(func() {
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)

		sqlDB, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Fatal("Failed to connect to database", err)
		}

		instance := bun.NewDB(sqlDB, pgdialect.New())
		dbInstance = &Database{DB: instance}
		dbInstance.SetMaxIdleConns(10)
		dbInstance.SetMaxOpenConns(25)
		dbInstance.SetConnMaxIdleTime(5 * time.Minute)

		if err := dbInstance.Ping(); err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}

		dbInstance.AddQueryHook(NewQueryHook(logger.Log, 200*time.Millisecond))
	})
}

func GetDB() *Database {
	return dbInstance
}

func (db *Database) InitQuery(ctx interface{}) *bun.DB {
	return db.DB
}

func Connect(cfg config.DatabaseConfig) *bun.DB {
	dsn := cfg.GetDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	bunDB := bun.NewDB(db, pgdialect.New())
	return bunDB
}
