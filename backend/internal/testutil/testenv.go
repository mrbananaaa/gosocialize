package testutil

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrbananaaa/gosocialize/internal/platform/db"
	"github.com/stretchr/testify/require"
)

type Env struct {
	DB *pgxpool.Pool
	Q  *db.Queries
}

func NewEnv(t *testing.T) *Env {
	t.Helper()

	err := Setup(context.Background())
	require.NoError(t, err)

	err = ResetDB(t)
	require.NoError(t, err)

	return &Env{
		DB: GetDB(),
		Q:  db.New(GetDB()),
	}
}
