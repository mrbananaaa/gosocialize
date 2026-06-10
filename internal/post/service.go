package post

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

type PostCursor struct {
	CreatedAt time.Time
	ID        uuid.UUID
}

func (s *Service) ListPosts(
	ctx context.Context,
	cursorStr string,
	limit int32,
) ([]Post, string, error) {
	cursor, err := decodeCursor(cursorStr)
	if err != nil {
		return nil, "", err
	}

	params := sqlc.ListPostsParams{
		PaginationLimit: limit + 1,
	}

	if cursor != nil {
		params.CursorCreatedAt = pgtype.Timestamptz{
			Time:  cursor.CreatedAt,
			Valid: true,
		}
		params.CursorID = cursor.ID
	}

	rows, err := s.q.ListPosts(ctx, params)
	if err != nil {
		return nil, "", err
	}

	hasMore := len(rows) > int(limit)

	var nextCursor string

	if hasMore && len(rows) > 0 {
		last := rows[len(rows)-1]

		c := PostCursor{
			CreatedAt: last.CreatedAt,
			ID:        last.ID,
		}

		b, _ := json.Marshal(c)
		nextCursor = base64.StdEncoding.EncodeToString(b)
	}

	var posts []Post
	for _, p := range rows {
		posts = append(posts, Post{
			ID:        p.ID,
			UserID:    p.UserID,
			Title:     p.Title,
			Content:   p.Content,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}

	return posts, nextCursor, nil
}

func decodeCursor(s string) (*PostCursor, error) {
	if s == "" {
		return nil, nil
	}

	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	var c PostCursor
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
