package post

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres/sqlc"
	"github.com/mrbananaaa/gosocialize/pkg/logger"
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

type ListPostsPayload struct {
	Posts      []Post
	NextCursor string
	HasMore    bool
}

func (s *Service) ListPosts(
	ctx context.Context,
	cursorStr string,
	limit int32,
) (*ListPostsPayload, error) {
	cursor, err := decodeCursor(cursorStr)
	if err != nil {
		return nil, err
	}

	paginationLimit := limit + 1

	var rows []sqlc.Post
	if cursor == nil {
		logger.Info("Post first fired")
		rows, err = s.q.ListPostsFirst(ctx, paginationLimit)
		if err != nil {
			return nil, err
		}
	} else {
		logger.Info("Post after fired")

		rows, err = s.q.ListPostsAfter(ctx, sqlc.ListPostsAfterParams{
			CursorCreatedAt: pgtype.Timestamptz{
				Time:  cursor.CreatedAt,
				Valid: true,
			},
			CursorID:        cursor.ID,
			PaginationLimit: paginationLimit,
		})
	}

	var p []Post
	for _, post := range rows {
		p = append(p, Post{
			ID:        post.ID,
			UserID:    post.UserID,
			Title:     post.Title,
			Content:   post.Content,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		})
	}

	hasMore := len(rows) > int(limit)
	if hasMore {
		p = p[:limit]
	}

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

	return &ListPostsPayload{
		Posts:      p,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
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
