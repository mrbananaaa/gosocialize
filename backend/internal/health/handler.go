package health

import (
	"net/http"

	"github.com/mrbananaaa/gosocialize/internal/platform/httpx"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	httpx.Message(w, http.StatusOK, "Service is up and running!")
}
