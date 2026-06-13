package post

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
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
		return nil, postgres.PgErrMapper(err)
	}

	return p, nil
}

func (s *Service) GetByID(
	ctx context.Context,
	postID uuid.UUID,
) (*Post, error) {
	post, err := s.q.FindPostByID(ctx, postID)
	if err != nil {
		return nil, postgres.PgErrMapper(err)
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
		// TODO: wrap with apperr
		return nil, err
	}

	paginationLimit := cursorQuery.Limit + 1

	var rows []sqlc.Post
	if cursor == nil {
		rows, err = s.q.ListPostsFirst(ctx, paginationLimit)
		if err != nil {
			return nil, postgres.PgErrMapper(err)
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
		if err != nil {
			return nil, postgres.PgErrMapper(err)
		}
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
			// TODO: wrap with apperr
			return nil, err
		}
	}

	return &ListPostsPayload{
		Posts:      p,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

type UpdatePostInput struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	Title   string
	Content string
}

func (s *Service) UpdatePost(
	ctx context.Context,
	input UpdatePostInput,
) error {
	post, err := s.q.FindPostByID(ctx, input.ID)
	if err != nil {
		err = postgres.PgErrMapper(err)
		return err
	}

	if post.UserID != input.UserID {
		return apperr.ErrForbidden
	}

	var title, content string
	if input.Title == "" {
		title = post.Title
	} else {
		title = input.Title
	}

	if input.Content == "" {
		content = post.Content
	} else {
		content = input.Content
	}

	err = s.q.UpdatePost(ctx, sqlc.UpdatePostParams{
		ID:      post.ID,
		Title:   title,
		Content: content,
	})
	if err != nil {
		return postgres.PgErrMapper(err)
	}

	return nil
}

type DeletePostInput struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (s *Service) DeletePost(
	ctx context.Context,
	input DeletePostInput,
) error {
	post, err := s.q.FindPostByID(ctx, input.ID)
	if err != nil {
		return postgres.PgErrMapper(err)
	}

	if post.UserID != input.UserID {
		return apperr.ErrForbidden
	}

	if err := s.q.DeletePost(ctx, post.ID); err != nil {
		return postgres.PgErrMapper(err)
	}

	return nil
}
