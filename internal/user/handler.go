package user

import (
	"net/http"

	"github.com/mrbananaaa/gosocialize/internal/auth"
	"github.com/mrbananaaa/gosocialize/internal/platform/httpx"
	"github.com/mrbananaaa/gosocialize/pkg/logger"
)

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

	// WARN: PASSWORD HASHER TESTING
	argonHasher := auth.NewArgon2Hasher()

	plainTextPassword := "panjangnyonipantek"
	hash, err := argonHasher.Hash(plainTextPassword)
	if err != nil {
		logger.Error("failed to hash password", logger.ErrorField(err))
		httpx.Error(w, err)
		return
	}

	logger.Info("hash length", logger.Int("length", len(hash)))
	logger.Info("password hashed", logger.String("original", "panjangnyonipantek"), logger.String("hash", TruncateWithEllipsis(hash, 12)))

	match, err := argonHasher.Compare(plainTextPassword, hash)
	if err != nil {
		logger.Error("failed to compare password", logger.ErrorField(err))
		httpx.Error(w, err)
		return
	}

	logger.Info("Password compare result", logger.Bool("match?", match))

	// TODO: dip some shi
	w.Write([]byte("RegisterUser handler"))
}

// WARN: REMOVE THIS LATER
func TruncateWithEllipsis(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "..."
}
