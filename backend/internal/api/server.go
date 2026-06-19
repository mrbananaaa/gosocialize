package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mrbananaaa/gosocialize/internal/auth"
	"github.com/mrbananaaa/gosocialize/internal/feed"
	"github.com/mrbananaaa/gosocialize/internal/follow"
	"github.com/mrbananaaa/gosocialize/internal/health"
	"github.com/mrbananaaa/gosocialize/internal/middlewares"
	"github.com/mrbananaaa/gosocialize/internal/platform/config"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres"
	"github.com/mrbananaaa/gosocialize/internal/platform/jwt"
	"github.com/mrbananaaa/gosocialize/internal/post"
	"github.com/mrbananaaa/gosocialize/internal/user"
	"github.com/mrbananaaa/gosocialize/pkg/logger"
	"github.com/mrbananaaa/gosocialize/store"
)

type Server struct {
	httpServer *http.Server
	config     *config.Config
	db         *postgres.DB
}

func NewServer(cfg *config.Config) (*Server, error) {
	db, err := postgres.NewDB(cfg.DB.URL)
	if err != nil {
		return nil, err
	}

	store := store.New(db.Pool)

	tokenService := jwt.New()
	argon2Hasher := auth.NewArgon2Hasher()

	authService := auth.NewService(tokenService, argon2Hasher, store)
	userService := user.NewService(store)
	postService := post.NewService(store)
	feedService := feed.NewService(store)
	followService := follow.NewService(store)

	authMiddleware := middlewares.NewAuth(tokenService)
	loggerMiddleware := middlewares.NewLogger()

	authHandler := auth.NewHandler(authService)
	userHandler := user.NewHandler(userService, followService)
	postHandler := post.NewHandler(postService)
	feedHandler := feed.NewHandler(feedService)
	healthHandler := health.NewHandler()

	handlers := Handlers{
		authHandler:   authHandler,
		userHandler:   userHandler,
		postHandler:   postHandler,
		feedHandler:   feedHandler,
		healthHandler: healthHandler,
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
