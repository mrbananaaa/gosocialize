package auth

import "github.com/go-chi/chi/v5"

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/signup", h.SignUp)
	r.Post("/signin", h.SignIn)

	return r
}
