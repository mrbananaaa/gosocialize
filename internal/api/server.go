package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mrbananaaa/gosocialize/internal/auth"
	"github.com/mrbananaaa/gosocialize/internal/middlewares"
	"github.com/mrbananaaa/gosocialize/internal/platform/config"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres"
	"github.com/mrbananaaa/gosocialize/internal/platform/jwt"
	"github.com/mrbananaaa/gosocialize/internal/post"
	"github.com/mrbananaaa/gosocialize/internal/user"
	"github.com/mrbananaaa/gosocialize/pkg/logger"
)

type Server struct {
	httpServer *http.Server
	config     *config.Config
	db         *postgres.DB
}

func NewServer(cfg *config.Config) (*Server, error) {
	db, err := postgres.New(cfg.DB.URL)
	if err != nil {
		return nil, err
	}

	tokenService := jwt.New()
	argon2Hasher := auth.NewArgon2Hasher()

	// userRepo := user.NewRepository()

	authService := auth.NewService(tokenService, argon2Hasher, db.Q)
	userService := user.NewService(db.Q)
	postService := post.NewService(db.Q)

	authMiddleware := middlewares.NewAuth(tokenService)
	loggerMiddleware := middlewares.NewLogger()

	authHandler := auth.NewHandler(authService)
	userHandler := user.NewHandler(userService)
	postHandler := post.NewHandler(postService)

	handlers := Handlers{
		authHandler: authHandler,
		userHandler: userHandler,
		postHandler: postHandler,
	}

	middlewares := Middlewares{
		authMiddleware:   authMiddleware,
		loggerMiddleware: loggerMiddleware,
	}

	mux := NewRouter(handlers, middlewares)

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
		db:         db,
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
	defer s.db.Close()

	return s.httpServer.Shutdown(ctx)
}
