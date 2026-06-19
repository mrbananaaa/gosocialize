package user

import (
	"github.com/go-chi/chi/v5"
	"github.com/mrbananaaa/gosocialize/internal/middlewares"
)

func (h *Handler) Routes(authMiddleware *middlewares.AuthMiddleware) chi.Router {
	r := chi.NewRouter()

	r.Group(func(u chi.Router) {
		u.Use(authMiddleware.WithAccessToken)

		// follow
		r.Post("/{userID}/follow", h.FollowUser)
		r.Delete("/{userID}/follow", h.UnfollowUser)
	})

	return r
}
