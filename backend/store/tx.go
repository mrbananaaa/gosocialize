package store

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mrbananaaa/gosocialize/internal/platform/db"
)

func (s *Store) WithTx(
	ctx context.Context,
	fn func(q *db.Queries) error,
) error {
	tx, err := s.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.Q.WithTx(tx)

	if err := fn(qtx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
