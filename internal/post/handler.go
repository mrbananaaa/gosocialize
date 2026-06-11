package post

import (
	"net/http"
	"strconv"

	"github.com/mrbananaaa/gosocialize/internal/platform/httpx"
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

func (h *Handler) GetAllPost(w http.ResponseWriter, r *http.Request) {
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

type GetPostQuery struct {
	Limit  int
	Cursor string // base64encoded
}

type GetPostResponse struct {
}
