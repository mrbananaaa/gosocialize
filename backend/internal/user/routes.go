package user

import (
	"github.com/go-chi/chi/v5"
)

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	// follow
	r.Post("/{userID}/follow", h.FollowUser)
	r.Delete("/{userID}/follow", h.UnfollowUser)

	return r
}
