package feed

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres"
	"github.com/mrbananaaa/gosocialize/internal/platform/db"
	"github.com/mrbananaaa/gosocialize/internal/post"
	"github.com/mrbananaaa/gosocialize/pkg/pagination"
	"github.com/mrbananaaa/gosocialize/store"
)

type Service interface {
	GetUserFeeds(ctx context.Context, userID uuid.UUID, cursorQuery pagination.CursorQueryParam) (*GetUserFeedsPayload, error)
}

type feedService struct {
	store *store.Store
}

func NewService(
	s *store.Store,
) Service {
	return &feedService{
		store: s,
	}
}

type GetUserFeedsPayload struct {
	Posts      []post.Post
	NextCursor string
	HasMore    bool
}

func (s *feedService) GetUserFeeds(
	ctx context.Context,
	userID uuid.UUID,
	cursorQuery pagination.CursorQueryParam,
) (*GetUserFeedsPayload, error) {
	cursor, err := pagination.DecodeCursor(cursorQuery.Cursor)
	if err != nil {
		return nil, apperr.New(apperr.Code.BadRequest, "invalid cursor")
	}

	paginationLimit := cursorQuery.Limit + 1

	var rows []db.Post
	if cursor == nil {
		rows, err = s.store.Q.UserFeedsFirst(ctx, db.UserFeedsFirstParams{
			UserID:          userID,
			PaginationLimit: paginationLimit,
		})
		if err != nil {
			return nil, postgres.PgErrMapper(err)
		}
	} else {
		rows, err = s.store.Q.UserFeedsAfter(ctx, db.UserFeedsAfterParams{
			UserID: userID,
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

	ps := make([]post.Post, 0, len(rows))
	for _, p := range rows {
		ps = append(ps, post.Post{
			ID:        p.ID,
			AuthorID:  p.AuthorID,
			Title:     p.Title,
			Content:   p.Content,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}

	hasMore := len(rows) > int(cursorQuery.Limit)
	if hasMore {
		ps = ps[:cursorQuery.Limit]
	}

	var nextCursor string
	if hasMore && len(rows) > 0 {
		lastPost := rows[len(rows)-1]
		c := pagination.Cursor{
			CreatedAt: lastPost.CreatedAt,
			ID:        lastPost.ID,
		}

		nextCursor, err = pagination.EncodeCursor(c)
		if err != nil {
			return nil, apperr.ErrInternal
		}
	}

	return &GetUserFeedsPayload{
		Posts:      ps,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}
