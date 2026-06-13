package post

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
	"github.com/mrbananaaa/gosocialize/internal/platform/httpx"
	"github.com/mrbananaaa/gosocialize/internal/platform/requestctx"
	"github.com/mrbananaaa/gosocialize/internal/platform/validator"
	"github.com/mrbananaaa/gosocialize/pkg/logger"
	"github.com/mrbananaaa/gosocialize/pkg/pagination"
)

type Handler struct {
	postService *Service
}

func NewHandler(postService *Service) *Handler {
	return &Handler{
		postService: postService,
	}
}

func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	cursorStr := r.URL.Query().Get("cursor")

	limit := int32(10)
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil {
			limit = int32(v)
		}
	}

	p, err := h.postService.ListPosts(r.Context(), pagination.CursorQueryParam{
		Cursor: cursorStr,
		Limit:  limit,
	})
	if err != nil {
		httpx.Error(w, err)
		return
	}

	httpx.OK(w, p.Posts, httpx.WithPaginationMeta(pagination.PaginationMeta{
		NextCursor: p.NextCursor,
		HasMore:    p.HasMore,
	}))
}

type CreatePostRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=255"`
	Content string `json:"content" validate:"required"`
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userCtx, exists := requestctx.UserFromContext(r.Context())
	if !exists {
		logger.Error("couldn't get user from context")
		httpx.Error(w, apperr.ErrUnauthorized)
		return
	}

	var req CreatePostRequest

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

	post, err := h.postService.Create(r.Context(), CreateInput{
		AuthorID: userCtx.ID,
		Title:    req.Title,
		Content:  req.Content,
	})
	if err != nil {
		logger.Error("failed to create post", logger.ErrorField(err))
		httpx.Error(w, err)
		return
	}

	httpx.Created(w, post)
}

func (h *Handler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	postIDstr := chi.URLParam(r, "postID")
	postID, err := uuid.Parse(postIDstr)
	if err != nil {
		logger.Error("failed to parse uuid", logger.ErrorField(err))
		httpx.Error(w, apperr.InvalidUUIDErr(err, "post_id"))
		return
	}

	post, err := h.postService.GetByID(r.Context(), postID)
	if err != nil {
		logger.Error("failed to get post", logger.ErrorField(err))
		httpx.Error(w, err)
		return
	}

	httpx.OK(w, post)
}

type UpdatePostRequest struct {
	Title   string `json:"title" validate:"max=255"`
	Content string `json:"content"`
}

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "postID")
	postID, err := uuid.Parse(postIDStr)
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

	var req UpdatePostRequest

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

	logger.Debug(
		"update post request body",
		logger.String("title", req.Title),
		logger.String("content", req.Content),
	)

	err = h.postService.UpdatePost(r.Context(), UpdatePostInput{
		ID:      postID,
		UserID:  userCtx.ID,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		logger.Error("failed to update post", logger.ErrorField(err))
		httpx.Error(w, err)
		return
	}

	httpx.Message(w, http.StatusOK, "post updated!")
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "postID")
	postID, err := uuid.Parse(postIDStr)
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

	err = h.postService.DeletePost(r.Context(), DeletePostInput{
		ID:       postID,
		AuthorID: userCtx.ID,
	})
	if err != nil {
		logger.Error("couldn't delete post", logger.ErrorField(err))
		httpx.Error(w, err)
		return
	}

	httpx.Message(w, http.StatusOK, "post deleted!")
}
