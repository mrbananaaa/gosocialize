package post

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres/sqlc"
	"github.com/mrbananaaa/gosocialize/pkg/pagination"
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

	err := s.q.CreatePost(ctx, sqlc.CreatePostParams{
		ID:        p.ID,
		UserID:    p.UserID,
		Title:     p.Title,
		Content:   p.Content,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	})
	if err != nil {
		err = postgres.PgErrMapper(err)
		return nil, err
	}

	return p, nil
}

func (s *Service) GetByID(
	ctx context.Context,
	postID uuid.UUID,
) (*Post, error) {
	post, err := s.q.FindPostByID(ctx, postID)
	if err != nil {
		err = postgres.PgErrMapper(err)
		return nil, err
	}

	return &Post{
		ID:        post.ID,
		UserID:    post.UserID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}, nil
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

type ListPostsPayload struct {
	Posts      []Post
	NextCursor string
	HasMore    bool
}

func (s *Service) ListPosts(
	ctx context.Context,
	cursorQuery pagination.CursorQueryParam,
) (*ListPostsPayload, error) {
	cursor, err := pagination.DecodeCursor(cursorQuery.Cursor)
	if err != nil {
		return nil, err
	}

	paginationLimit := cursorQuery.Limit + 1

	var rows []sqlc.Post
	if cursor == nil {
		rows, err = s.q.ListPostsFirst(ctx, paginationLimit)
		if err != nil {
			return nil, err
		}
	} else {
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

	hasMore := len(rows) > int(cursorQuery.Limit)
	if hasMore {
		p = p[:cursorQuery.Limit]
	}

	var nextCursor string
	if hasMore && len(rows) > 0 {
		last := rows[len(rows)-1]

		c := pagination.Cursor{
			CreatedAt: last.CreatedAt,
			ID:        last.ID,
		}

		nextCursor, err = pagination.EncodeCursor(c)
		if err != nil {
			return nil, err
		}
	}

	return &ListPostsPayload{
		Posts:      p,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}
