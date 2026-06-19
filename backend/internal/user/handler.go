package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/follow"
	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
	"github.com/mrbananaaa/gosocialize/internal/platform/httpx"
	"github.com/mrbananaaa/gosocialize/internal/platform/requestctx"
	"github.com/mrbananaaa/gosocialize/pkg/logger"
)

type Handler struct {
	userSvc   *Service
	followSvc follow.Service
}

func NewHandler(
	userSvc *Service,
	followSvc follow.Service,
) *Handler {
	return &Handler{
		userSvc:   userSvc,
		followSvc: followSvc,
	}
}

func (h *Handler) FollowUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	followeeID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.Error("failed to parse uuid", logger.ErrorField(err))
		httpx.Error(w, apperr.InvalidUUIDErr(err, "post_id"))
		return
	}

	userCtx, exists := requestctx.UserFromContext(r.Context())
	if !exists {
		logger.Error("couldn't get user from context")
		httpx.Error(w, apperr.ErrUnauthorized)
		return
	}

	err = h.followSvc.Follow(r.Context(), userCtx.ID, followeeID)
	if err != nil {
		logger.Error("failed to follow user", logger.ErrorField(err))
		httpx.Error(w, err)
		return
	}

	httpx.Message(w, http.StatusCreated, "user followed!")
}

func (h *Handler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	followeeID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.Error("failed to parse uuid", logger.ErrorField(err))
		httpx.Error(w, apperr.InvalidUUIDErr(err, "post_id"))
		return
	}

	userCtx, exists := requestctx.UserFromContext(r.Context())
	if !exists {
		logger.Error("couldn't get user from context")
		httpx.Error(w, apperr.ErrUnauthorized)
		return
	}

	err = h.followSvc.Unfollow(r.Context(), userCtx.ID, followeeID)
	if err != nil {
		logger.Error("failed to unfollow user", logger.ErrorField(err))
		httpx.Error(w, err)
		return
	}

	httpx.Message(w, http.StatusOK, "user unfollowed!")
}
