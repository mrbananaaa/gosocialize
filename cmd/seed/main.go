package main

import (
	"context"

	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres"
	"github.com/mrbananaaa/gosocialize/pkg/logger"
)

func main() {
	err := logger.Init(logger.Config{
		Development: true,
	})
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	db, err := postgres.New("postgres://postgres:postgres@localhost:5432/gosocialize?sslmode=disable")
	if err != nil {
		logger.Error("Failed to open db pool", logger.ErrorField(err))
		return
	}
	defer db.Close()

	if err := Seed(context.Background(), db); err != nil {
		logger.Error("Failed to seed the database", logger.ErrorField(err))
		return
	}
}
