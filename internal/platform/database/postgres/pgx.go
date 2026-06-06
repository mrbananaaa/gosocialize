package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres/sqlc"
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
