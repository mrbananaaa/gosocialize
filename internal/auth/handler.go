package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
	"github.com/mrbananaaa/gosocialize/internal/platform/httpx"
	"github.com/mrbananaaa/gosocialize/internal/platform/validator"
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

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
}

type SignUpResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to parse request body", logger.ErrorField(err))
		httpx.Error(w, apperr.New("BAD REQUEST", "bad request body"))
		return
	}

	if err := validator.ValidateStruct(req); err != nil {
		logger.Error("failed to valdiate request body", logger.ErrorField(err))
		httpx.Error(w, err)
		return
	}

	u, err := h.authService.Register(r.Context(), RegisterInput{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		logger.Error("Failed to register user", logger.ErrorField(err))
		httpx.Error(w, err)
		return
	}

	httpx.Created(w, SignUpResponse{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	})
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	httpx.Message(w, http.StatusOK, "auth/signin handler")
}
