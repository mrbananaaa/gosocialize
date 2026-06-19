package user

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/follow"
	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
	"github.com/mrbananaaa/gosocialize/internal/platform/httpx"
	"github.com/mrbananaaa/gosocialize/internal/platform/validator"
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

type FollowUserRequest struct {
	TargetUserID uuid.UUID `json:"target_user_id" validate:"required"`
}

func (h *Handler) FollowUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.Error("failed to parse uuid", logger.ErrorField(err))
		httpx.Error(w, apperr.InvalidUUIDErr(err, "post_id"))
		return
	}

	var req FollowUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to parse request body", logger.ErrorField(err))
		httpx.Error(w, apperr.DecodeBodyErr(err))
		return
	}

	if err := validator.ValidateStruct(req); err != nil {
		logger.Error("failed to valdiate request body", logger.ErrorField(err))
		httpx.Error(w, err)
		return
	}

	err = h.followSvc.Follow(r.Context(), userID, req.TargetUserID)
	if err != nil {
		logger.Error("failed to follow user", logger.ErrorField(err))
		httpx.Error(w, err)
		return
	}

	httpx.Message(w, http.StatusCreated, "user followed!")
}

type UnFollowUserRequest struct {
	TargetUserID uuid.UUID `json:"target_user_id" validate:"required"`
}

func (h *Handler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.Error("failed to parse uuid", logger.ErrorField(err))
		httpx.Error(w, apperr.InvalidUUIDErr(err, "post_id"))
		return
	}

	var req FollowUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to parse request body", logger.ErrorField(err))
		httpx.Error(w, apperr.DecodeBodyErr(err))
		return
	}

	if err := validator.ValidateStruct(req); err != nil {
		logger.Error("failed to valdiate request body", logger.ErrorField(err))
		httpx.Error(w, err)
		return
	}

	err = h.followSvc.Unfollow(r.Context(), userID, req.TargetUserID)
	if err != nil {
		logger.Error("failed to unfollow user", logger.ErrorField(err))
		httpx.Error(w, err)
		return
	}

	httpx.Message(w, http.StatusOK, "user unfollowed!")
}
