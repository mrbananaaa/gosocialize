package post

import (
	"net/http"

	"github.com/mrbananaaa/gosocialize/internal/platform/httpx"
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
	posts, err := h.postService.GetPosts(r.Context())
	if err != nil {
		httpx.Error(w, err)
		return
	}

	httpx.OK(w, posts)
}
