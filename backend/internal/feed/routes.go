package feed

import (
	"github.com/go-chi/chi/v5"
	"github.com/mrbananaaa/gosocialize/internal/middlewares"
)

func (h *Handler) Routes(am *middlewares.AuthMiddleware) chi.Router {
	r := chi.NewRouter()

	r.Use(am.WithAccessToken)
	r.Get("/", h.GetUserFeeds)

	return r
}
