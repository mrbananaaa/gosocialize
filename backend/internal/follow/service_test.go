package follow_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/follow"
	"github.com/mrbananaaa/gosocialize/internal/testutil"
	"github.com/mrbananaaa/gosocialize/store"
	"github.com/stretchr/testify/require"
)

func TestFollow(t *testing.T) {
	t.Run("success - user can follow other user", func(t *testing.T) {
		env, svc := newFollowService(t)

		user1 := env.CreateUser(t)
		user2 := env.CreateUser(t)

		err := svc.Follow(t.Context(), user1.ID, user2.ID)
		require.NoError(t, err)
	})

	t.Run("fail - cannot follow user that doesn't exits", func(t *testing.T) {
		env, svc := newFollowService(t)

		follower := env.CreateUser(t)

		err := svc.Follow(t.Context(), uuid.New(), follower.ID)
		require.Error(t, err)
	})
}

func TestUnfollow(t *testing.T) {
	t.Run("success - user can unfollow", func(t *testing.T) {

	})
}

func newFollowService(t *testing.T) (*testutil.Env, follow.Service) {
	t.Helper()

	env := testutil.NewEnv(t)
	store := store.New(env.DB)
	svc := follow.NewService(store)

	return env, svc
}
