package post

import "github.com/go-chi/chi/v5"

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.GetAllPost)
	r.Post("/", h.CreatePost)

	return r
}
