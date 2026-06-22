package follow_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/follow"
	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
	"github.com/mrbananaaa/gosocialize/internal/testutil"
	"github.com/mrbananaaa/gosocialize/store"
	"github.com/stretchr/testify/require"
)

func TestFollow(t *testing.T) {
	env, svc := newFollowService(t)

	t.Run("success - user can follow other user", func(t *testing.T) {
		user1 := env.CreateUser(t)
		user2 := env.CreateUser(t)

		err := svc.Follow(t.Context(), user1.ID, user2.ID)
		require.NoError(t, err)
	})

	t.Run("fail - can't follow non-existent user", func(t *testing.T) {
		follower := env.CreateUser(t)

		err := svc.Follow(t.Context(), follower.ID, uuid.New())

		var appErr *apperr.Error
		require.ErrorAs(t, err, &appErr)
		require.Equal(t, "can't follow non-existent user", appErr.Message)
	})
}

func TestUnfollow(t *testing.T) {
	env, svc := newFollowService(t)

	t.Run("success - user can unfollow", func(t *testing.T) {
		user1 := env.CreateUser(t)
		user2 := env.CreateUser(t)

		err := svc.Follow(t.Context(), user1.ID, user2.ID)
		require.NoError(t, err)

		err = svc.Unfollow(t.Context(), user1.ID, user2.ID)
		require.NoError(t, err)

		count, err := env.Q.CountFollowers(t.Context(), user2.ID)
		require.NoError(t, err)
		require.Equal(t, 0, int(count))
	})

	t.Run("fail - can't unfollow unfollowed users", func(t *testing.T) {
		user := env.CreateUser(t)
		follower := env.CreateUser(t)

		err := svc.Unfollow(t.Context(), follower.ID, user.ID)

		var appErr *apperr.Error
		require.ErrorAs(t, err, &appErr)
		require.Equal(t, "you're not following this user", appErr.Message)
	})

	t.Run("fail - user can't unfollow non-existent user", func(t *testing.T) {
		follower := env.CreateUser(t)

		err := svc.Unfollow(t.Context(), follower.ID, uuid.New())

		var appErr *apperr.Error
		require.ErrorAs(t, err, &appErr)
		require.Equal(t, "can't found followed user", appErr.Message)
	})
}

func TestFollowerCount(t *testing.T) {
	env, svc := newFollowService(t)

	t.Run("success - followers counted", func(t *testing.T) {
		user := env.CreateUser(t)
		followerTotal := 15

		for range followerTotal {
			follower := env.CreateUser(t)
			err := svc.Follow(t.Context(), follower.ID, user.ID)
			require.NoError(t, err)
		}

		count, err := svc.FollowerCount(t.Context(), user.ID)
		require.NoError(t, err)
		require.Equal(t, followerTotal, count)
	})

	t.Run("fail - can't count follower of non-existent user", func(t *testing.T) {
		_, err := svc.FollowerCount(t.Context(), uuid.New())

		var appErr *apperr.Error
		require.ErrorAs(t, err, &appErr)
		require.Equal(t, "can't count follower of non-existent user", appErr.Message)
	})
}

func newFollowService(t *testing.T) (*testutil.Env, follow.Service) {
	t.Helper()

	env := testutil.NewEnv(t)
	store := store.New(env.DB)
	svc := follow.NewService(store)

	return env, svc
}
