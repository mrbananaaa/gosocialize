package user

import (
	"github.com/go-chi/chi/v5"
)

// TODO: mount this route using r.Mount(pattern, Routes())
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	// r.Post("/signup", h.RegisterUser)

	return r
}
