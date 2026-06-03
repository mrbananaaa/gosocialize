package user

import "net/http"

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type Handler struct {
	service *Service
}

func NewHandler(
	service *Service,
) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// var req RegisterRequest

	// TODO: do some decode and validation shi

	// TODO: use service shi

	// TODO: dip some shi
	w.Write([]byte("RegisterUser handler"))
}
