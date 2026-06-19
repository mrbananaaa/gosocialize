package follow

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/platform/db"
	"github.com/mrbananaaa/gosocialize/store"
)

type Service interface {
	Follow(ctx context.Context, followeeID, followerID uuid.UUID) error
	Unfollow(ctx context.Context, followeeID, followerID uuid.UUID) error
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

func (s *followService) Follow(ctx context.Context, followeeID, followerID uuid.UUID) error {
	return s.store.Q.FollowUser(ctx, db.FollowUserParams{
		FolloweeID: followeeID,
		FollowerID: followerID,
	})
}

func (s *followService) Unfollow(ctx context.Context, followeeID, followerID uuid.UUID) error {
	return s.store.Q.UnfollowUser(ctx, db.UnfollowUserParams{
		FolloweeID: followeeID,
		FollowerID: followerID,
	})
}

func (s *followService) FollowerCount(ctx context.Context, userID uuid.UUID) (int, error) {
	count, err := s.store.Q.FollowersCount(ctx, userID)
	return int(count), err
}
