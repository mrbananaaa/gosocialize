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

	r.Group(func(u chi.Router) {
		u.Use(authMiddleware.WithAccessToken)

		u.Post("/", h.CreatePost)
		u.Patch("/{postID}", h.UpdatePost)
		u.Delete("/{postID}", h.DeletePost)
	})

	return r
}
