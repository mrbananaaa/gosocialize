package auth

import (
	"net/http"

	"github.com/mrbananaaa/gosocialize/internal/platform/httpx"
	"github.com/mrbananaaa/gosocialize/pkg/logger"
)

type Handler struct {
	authService *Service
}

func NewHandler(authService *Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	err := h.authService.Register(r.Context(), RegisterInput{
		Email:    "test@mail.com",
		Username: "hardcodedusername",
		Password: "wellit'snotsafetogaveurpassword",
		Name:     "John Doe",
	})
	if err != nil {
		logger.Error("Failed to register user", logger.ErrorField(err))
		httpx.Error(w, err)
		return
	}

	httpx.OK(w, nil)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	httpx.Message(w, http.StatusOK, "auth/signin handler")
}
