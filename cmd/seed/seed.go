package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/migrations"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres/seeds"
	"github.com/mrbananaaa/gosocialize/pkg/logger"
	"github.com/pressly/goose/v3"
)

func Seed(
	ctx context.Context,
	db *postgres.DB,
) error {
	logger.Info("Initialize seeding")

	if err := ResetDB(ctx, db); err != nil {
		return err
	}

	logger.Info("Seeding users...")
	users, err := seeds.SeedUsers(ctx, db.Q, 20, "seedingpassword")
	if err != nil {
		return fmt.Errorf("Failed to seed users: %v", err)
	}

	logger.Info("Seeding posts...")
	_, err = seeds.SeedPosts(ctx, db.Q, 100, users)
	if err != nil {
		return fmt.Errorf("Failed to seed users: %v", err)
	}

	logger.Info("Seeding done")
	return nil
}

func ResetDB(
	ctx context.Context,
	db *postgres.DB,
) error {
	// collecting connection from pool
	dbconn := stdlib.OpenDBFromPool(db.Pool)
	defer dbconn.Close()

	if err := dbconn.Ping(); err != nil {
		return fmt.Errorf("Failed to reach database: %v", err)
	}

	// configuring goose
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("Failed to select dialect: %v", err)
	}
	goose.SetBaseFS(migrations.FS)
	goose.SetLogger(goose.NopLogger())

	logger.Warn("Cleaning database...")
	if err := goose.ResetContext(ctx, dbconn, "."); err != nil {
		return fmt.Errorf("Failed to reset goose migrations: %v", err)
	}
	logger.Info("Database cleaned!")

	logger.Info("Migrating database...")
	if err := goose.UpContext(ctx, dbconn, "."); err != nil {
		return fmt.Errorf("Failed to migrate database: %v", err)
	}
	logger.Info("Migration successful")

	logger.Info("Database reseted!")

	return nil
}
