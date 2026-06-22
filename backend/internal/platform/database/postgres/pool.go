package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrbananaaa/gosocialize/pkg/logger"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB(connStr string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	go func() {
		var try int

		for {
			err := pool.Ping(context.Background())
			if err == nil {
				logger.Warn("database connected 🔌")
				return
			}

			try++
			if try >= 5 {
				panic(err)
			}
			logger.Warn("Failed to reach database 💣", logger.Int("try_count", try))
			time.Sleep(2 * time.Second)
		}
	}()

	return &DB{
		Pool: pool,
	}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
	logger.Warn("database pool closed 🔌")
}
