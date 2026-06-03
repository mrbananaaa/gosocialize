package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/mrbananaaa/gosocialize/pkg/config"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config) (*Server, error) {
	mux := NewRouter()

	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	return &Server{
		httpServer: httpServer,
	}, nil
}

func (s *Server) Run() error {
	// TODO: Use environment variable from Config
	log.Println("Server listening on :8080")
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
