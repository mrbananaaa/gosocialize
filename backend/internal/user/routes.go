package user

import (
	"github.com/go-chi/chi/v5"
	"github.com/mrbananaaa/gosocialize/internal/middlewares"
)

func (h *Handler) Routes(authMiddleware *middlewares.AuthMiddleware) chi.Router {
	r := chi.NewRouter()

	r.Group(func(u chi.Router) {
		u.Use(authMiddleware.WithAccessToken)

		u.Get("/feed", h.UserFeeds)

		// follow
		u.Post("/{userID}/follow", h.FollowUser)
		u.Delete("/{userID}/follow", h.UnfollowUser)
	})

	return r
}
