package feed

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
	"github.com/mrbananaaa/gosocialize/internal/platform/httpx"
	"github.com/mrbananaaa/gosocialize/internal/platform/requestctx"
	"github.com/mrbananaaa/gosocialize/internal/post"
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

	httpx.OK(w, toFeedsResponse(feeds))
}

type FeedResponse struct {
	ID        uuid.UUID `json:"id"`
	AuthorID  uuid.UUID `json:"author_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toPostResponse(p *post.Post) FeedResponse {
	return FeedResponse{
		ID:        p.ID,
		AuthorID:  p.AuthorID,
		Title:     p.Title,
		Content:   p.Content,
		Tags:      p.Tags,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func toFeedsResponse(p []post.Post) []FeedResponse {
	r := make([]FeedResponse, 0, len(p))
	for _, post := range p {
		r = append(r, toPostResponse(&post))
	}
	return r
}
