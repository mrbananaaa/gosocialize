package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres/sqlc"
	"github.com/mrbananaaa/gosocialize/pkg/logger"
)

type DB struct {
	Pool *pgxpool.Pool
	Q    *sqlc.Queries
}

func New(conn string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), conn)
	if err != nil {
		return nil, err
	}

	go func() {
		var try int

		for {
			err := pool.Ping(context.Background())
			if err != nil {
				try++
				logger.Warn("Failed to reach database", logger.Int("try_count", try))
			}

			if try >= 5 {
				panic(err)
			}

			time.Sleep(2 * time.Second)
		}
	}()

	return &DB{
		Pool: pool,
		Q:    sqlc.New(pool),
	}, nil
}

func (db *DB) WithTx(
	ctx context.Context,
	fn func(*sqlc.Queries) error,
) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	q := db.Q.WithTx(tx)

	if err := fn(q); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
