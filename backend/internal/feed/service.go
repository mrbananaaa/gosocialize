package feed

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres/sqlc"
	"github.com/mrbananaaa/gosocialize/internal/post"
)

type Service struct {
	q *sqlc.Queries
}

func NewService(q *sqlc.Queries) *Service {
	return &Service{
		q: q,
	}
}

func (s *Service) GetFeeds(
	ctx context.Context,
	userID uuid.UUID,
) ([]post.Post, error) {
	feeds, err := s.q.ListPostForFeeds(ctx, userID)
	if err != nil {
		return nil, postgres.PgErrMapper(err)
	}

	var f []post.Post
	for _, p := range feeds {
		f = append(f, post.Post{
			ID:        p.ID,
			AuthorID:  p.AuthorID,
			Title:     p.Title,
			Content:   p.Content,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}

	return f, nil
}
