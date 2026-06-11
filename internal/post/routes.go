package post

import (
	"github.com/go-chi/chi/v5"
	"github.com/mrbananaaa/gosocialize/internal/middlewares"
)

func (h *Handler) Routes(
	authMiddleware *middlewares.AuthMiddleware,
) chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.ListPosts)
	r.Get("/{postID}", h.GetPostByID)

	// TODO: refactor and move create post with auth
	r.Post("/", h.CreatePost)

	r.Group(func(u chi.Router) {
		u.Use(authMiddleware.WithAccessToken)

		u.Patch("/{postID}", h.UpdatePost)
	})

	return r
}
