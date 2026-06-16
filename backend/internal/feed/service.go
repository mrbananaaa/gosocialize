package feed

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres"
	"github.com/mrbananaaa/gosocialize/internal/post"
	"github.com/mrbananaaa/gosocialize/store"
)

type Service struct {
	store *store.Store
}

func NewService(
	s *store.Store,
) *Service {
	return &Service{
		store: s,
	}
}

func (s *Service) GetFeeds(
	ctx context.Context,
	userID uuid.UUID,
) ([]post.Post, error) {
	feeds, err := s.store.Q.ListPostForFeeds(ctx, userID)
	if err != nil {
		return nil, postgres.PgErrMapper(err)
	}

	f := make([]post.Post, 0, len(feeds))
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
