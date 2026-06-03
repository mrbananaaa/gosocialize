package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mrbananaaa/gosocialize/pkg/config"
	"github.com/mrbananaaa/gosocialize/pkg/logger"
)

type Server struct {
	httpServer *http.Server
	config     *config.Config
}

func NewServer(cfg *config.Config) (*Server, error) {
	mux := NewRouter()

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	return &Server{
		httpServer: httpServer,
		config:     cfg,
	}, nil
}

func (s *Server) Run() error {
	logger.Info(
		"http server listening 👂",
		logger.String("PORT", s.config.Server.Port),
	)
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
