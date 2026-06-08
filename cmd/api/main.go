package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/mrbananaaa/gosocialize/internal/api"
	"github.com/mrbananaaa/gosocialize/internal/platform/config"
	"github.com/mrbananaaa/gosocialize/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	err = logger.Init(logger.Config{
		Development: cfg.App.Env == "development",
	})
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	s, err := api.NewServer(cfg)
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info(
			"App is up and running ✨",
			logger.String("service", cfg.App.Name),
			logger.String("env", cfg.App.Env),
		)
		if err := s.Run(); err != nil {
			logger.Fatal(
				"Failed to start the app",
				err,
			)
		}
	}()

	<-ctx.Done()
	logger.Warn("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Failed to shutdown http server", err)
	}

	logger.Warn("Server closed gracefully")
}
