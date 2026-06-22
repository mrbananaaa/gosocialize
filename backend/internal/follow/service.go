package follow

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres"
	"github.com/mrbananaaa/gosocialize/internal/platform/db"
	"github.com/mrbananaaa/gosocialize/store"
)

type Service interface {
	Follow(ctx context.Context, followerID, followeeID uuid.UUID) error
	Unfollow(ctx context.Context, followerID, followeeID uuid.UUID) error
	FollowerCount(ctx context.Context, userID uuid.UUID) (int, error)
}

type followService struct {
	store *store.Store
}

func NewService(store *store.Store) Service {
	return &followService{
		store: store,
	}
}

func (s *followService) Follow(ctx context.Context, followerID, followeeID uuid.UUID) error {
	_, err := s.store.Q.FindUserByID(ctx, followeeID)
	if err != nil {
		err = postgres.PgErrMapper(err)

		if errors.Is(err, apperr.ErrNotFound) {
			return apperr.New(apperr.Code.NotFound, "can't follow non-existent user")
		}

		return err
	}

	err = s.store.Q.FollowUser(ctx, db.FollowUserParams{
		FollowerID: followerID,
		FolloweeID: followeeID,
	})

	return postgres.PgErrMapper(err)
}

func (s *followService) Unfollow(ctx context.Context, followerID, followeeID uuid.UUID) error {
	_, err := s.store.Q.FindUserByID(ctx, followeeID)
	if err != nil {
		err = postgres.PgErrMapper(err)

		if errors.Is(err, apperr.ErrNotFound) {
			return apperr.New(apperr.Code.NotFound, "can't found followed user")
		}

		return err
	}

	isFollowing, err := s.store.Q.IsFollowing(ctx, db.IsFollowingParams{
		FollowerID: followerID,
		FolloweeID: followeeID,
	})
	if err != nil {
		return postgres.PgErrMapper(err)
	}

	if !isFollowing {
		return apperr.New(apperr.Code.BadRequest, "you're not following this user")
	}

	err = s.store.Q.UnfollowUser(ctx, db.UnfollowUserParams{
		FollowerID: followerID,
		FolloweeID: followeeID,
	})

	return postgres.PgErrMapper(err)
}

func (s *followService) FollowerCount(ctx context.Context, userID uuid.UUID) (int, error) {
	_, err := s.store.Q.FindUserByID(ctx, userID)
	if err != nil {
		err = postgres.PgErrMapper(err)

		if errors.Is(err, apperr.ErrNotFound) {
			return 0, apperr.New(apperr.Code.NotFound, "can't count follower of non-existent user")
		}

		return 0, err
	}

	count, err := s.store.Q.CountFollowers(ctx, userID)
	return int(count), postgres.PgErrMapper(err)
}
