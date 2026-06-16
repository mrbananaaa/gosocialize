package post

import (
	"time"

	"github.com/google/uuid"
)

type CreatePostRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=255"`
	Content string `json:"content" validate:"required"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" validate:"max=255"`
	Content string `json:"content"`
}

type PostResponse struct {
	ID        uuid.UUID `json:"id"`
	AuthorID  uuid.UUID `json:"author_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toPostResponse(p *Post) PostResponse {
	return PostResponse{
		ID:        p.ID,
		AuthorID:  p.AuthorID,
		Title:     p.Title,
		Content:   p.Content,
		Tags:      p.Tags,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func toPostListResponse(p []Post) []PostResponse {
	r := make([]PostResponse, 0, len(p))
	for _, post := range p {
		r = append(r, toPostResponse(&post))
	}
	return r
}
