package feed

import (
	"net/http"

	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
	"github.com/mrbananaaa/gosocialize/internal/platform/httpx"
	"github.com/mrbananaaa/gosocialize/internal/platform/requestctx"
	"github.com/mrbananaaa/gosocialize/pkg/logger"
)

type Handler struct {
	feedService *Service
}

func NewHandler(feedService *Service) *Handler {
	return &Handler{
		feedService: feedService,
	}
}

func (h *Handler) GetUserFeeds(w http.ResponseWriter, r *http.Request) {
	userCtx, exists := requestctx.UserFromContext(r.Context())
	if !exists {
		httpx.Error(w, apperr.ErrUnauthorized)
		return
	}

	feeds, err := h.feedService.GetFeeds(r.Context(), userCtx.ID)
	if err != nil {
		logger.Error("Couldn't get feeds for user", logger.ErrorField(err))
		httpx.Error(w, err)
		return
	}

	httpx.OK(w, feeds)
}
