package post

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres/sqlc"
)

type Service struct {
	q *sqlc.Queries
}

func NewService(q *sqlc.Queries) *Service {
	return &Service{
		q: q,
	}
}

type CreateInput struct {
	UserID  uuid.UUID
	Title   string
	Content string
	Tags    []string
}

func (s *Service) Create(
	ctx context.Context,
	input CreateInput,
) (*Post, error) {

	// check if the user id exists

	// create post

	id := uuid.New()
	now := time.Now()

	p := &Post{
		ID:        id,
		UserID:    input.UserID,
		Title:     input.Title,
		Content:   input.Content,
		Tags:      input.Tags,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return p, nil
}

func (s *Service) GetPosts(
	ctx context.Context,
) ([]*Post, error) {
	var posts []*Post

	p, err := s.q.GetPost(ctx)
	if err != nil {
		err = postgres.PgErrMapper(err)
		return nil, err
	}

	for _, post := range p {
		posts = append(posts, &Post{
			ID:        post.ID,
			UserID:    post.UserID,
			Title:     post.Title,
			Content:   post.Content,
			Tags:      []string{},
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		})
	}

	return posts, nil
}

func (s *Service) GetPost() {}

type DeleteInput struct {
	ID uuid.UUID
}

func (s *Service) Delete(
	ctx context.Context,
	input DeleteInput,
) error {
	return nil
}
