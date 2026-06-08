package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/mrbananaaa/gosocialize/internal/auth"
	"github.com/mrbananaaa/gosocialize/internal/middlewares"
	"github.com/mrbananaaa/gosocialize/internal/platform/httpx"
	"github.com/mrbananaaa/gosocialize/internal/user"
)

type Handlers struct {
	authHandler *auth.Handler
	userHandler *user.Handler
}

type Middlewares struct {
	authMiddleware *middlewares.AuthMiddleware
}

func NewRouter(h Handlers, m Middlewares) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromHeader("X-Real-IP"))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		// TODO: change cors origin
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/v1", func(u chi.Router) {
		u.Mount("/auth", h.authHandler.Routes())
		u.Mount("/user", h.userHandler.Routes())

		// TODO: move this to separate packages
		u.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			d, err := json.Marshal(struct {
				Message string `json:"message"`
			}{
				Message: "OK",
			})
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(d)
		})

		// auth test
		u.Route("/priv", func(x chi.Router) {
			x.Use(m.authMiddleware.WithAccessToken)

			x.Get("/", func(w http.ResponseWriter, r *http.Request) {
				httpx.Message(w, http.StatusOK, "this is private route")
			})
		})
	})

	return r
}
