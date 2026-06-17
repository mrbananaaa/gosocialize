package testutil

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/migrations"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	once sync.Once

	dbPool   *pgxpool.Pool
	setupErr error
)

func Setup(ctx context.Context) error {
	once.Do(func() {
		pgContainer, err := postgres.Run(ctx,
			"postgres:17-alpine",
			postgres.WithDatabase("gosocialize"),
			postgres.WithUsername("postgres"),
			postgres.WithPassword("postgres"),
			postgres.BasicWaitStrategies(),
		)
		if err != nil {
			setupErr = fmt.Errorf("couldn't start postgres container: %v", err)
			return
		}

		connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			setupErr = fmt.Errorf("couldn't get pg connection string: %v", err)
			return
		}

		pool, err := pgxpool.New(ctx, connStr)
		if err != nil {
			setupErr = fmt.Errorf("failed to create pgx pool: %v", err)
			return
		}

		dbPool = pool
	})

	return setupErr
}

func ResetDB(t *testing.T) error {
	t.Helper()
	ctx := t.Context()

	dbconn := stdlib.OpenDBFromPool(dbPool)
	defer dbconn.Close()

	if err := dbconn.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %v", err)
	}
	goose.SetBaseFS(migrations.FS)
	goose.SetLogger(goose.NopLogger())

	if err := goose.ResetContext(ctx, dbconn, "."); err != nil {
		if err = goose.UpContext(ctx, dbconn, "."); err != nil {
			return fmt.Errorf("failed to migrate database: %v", err)
		}

		return nil
	}

	if err := goose.UpContext(ctx, dbconn, "."); err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	t.Log("Database reseted!")

	return nil
}

func GetDB() *pgxpool.Pool {
	return dbPool
}
