package store

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrbananaaa/gosocialize/internal/platform/db"
)

type Store struct {
	Pool *pgxpool.Pool
	Q    *db.Queries
}

func New(
	pool *pgxpool.Pool,
) *Store {
	return &Store{
		Pool: pool,
		Q:    db.New(pool),
	}
}
