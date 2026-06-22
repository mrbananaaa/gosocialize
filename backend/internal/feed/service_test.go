package feed_test

import (
	"testing"

	"github.com/mrbananaaa/gosocialize/internal/feed"
	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
	"github.com/mrbananaaa/gosocialize/internal/testutil"
	"github.com/mrbananaaa/gosocialize/pkg/pagination"
	"github.com/mrbananaaa/gosocialize/store"
	"github.com/stretchr/testify/require"
)

func TestGetFeeds(t *testing.T) {
	env, svc := newFeedService(t)

	t.Run("success - get feeds", func(t *testing.T) {
		user := env.CreateUser(t)
		_ = env.CreatePost(t, user.ID)

		payload, err := svc.GetUserFeeds(t.Context(), user.ID, pagination.CursorQueryParam{
			Cursor: "",
			Limit:  10,
		})
		require.NoError(t, err)
		require.Len(t, payload.Posts, 1)
	})

	t.Run("success - get cursor", func(t *testing.T) {
		user := env.CreateUser(t)
		postTotal := 30

		for range postTotal {
			u := env.CreateUser(t)
			_ = env.CreatePost(t, u.ID)
			env.CreateFollow(t, user.ID, u.ID)
		}

		payload, err := svc.GetUserFeeds(t.Context(), user.ID, pagination.CursorQueryParam{
			Cursor: "",
			Limit:  20,
		})
		require.NoError(t, err)
		require.Len(t, payload.Posts, 20)
		require.NotEmpty(t, payload.NextCursor)
		require.True(t, payload.HasMore)
	})

	t.Run("success - got next feeds", func(t *testing.T) {
		user := env.CreateUser(t)
		postTotal := 30

		for range postTotal {
			u := env.CreateUser(t)
			_ = env.CreatePost(t, u.ID)
			env.CreateFollow(t, user.ID, u.ID)
		}

		first, err := svc.GetUserFeeds(t.Context(), user.ID, pagination.CursorQueryParam{
			Cursor: "",
			Limit:  10,
		})
		require.NoError(t, err)
		require.Len(t, first.Posts, 10)
		require.NotEmpty(t, first.NextCursor)
		require.True(t, first.HasMore)

		next, err := svc.GetUserFeeds(t.Context(), user.ID, pagination.CursorQueryParam{
			Cursor: first.NextCursor,
			Limit:  5,
		})
		require.NoError(t, err)
		require.Len(t, next.Posts, 5)
		require.NotEmpty(t, next.NextCursor)
		require.True(t, next.HasMore)
	})

	t.Run("fail - invalid cursor", func(t *testing.T) {
		user := env.CreateUser(t)
		postTotal := 30

		for range postTotal {
			u := env.CreateUser(t)
			_ = env.CreatePost(t, u.ID)
			env.CreateFollow(t, user.ID, u.ID)
		}

		payload, err := svc.GetUserFeeds(t.Context(), user.ID, pagination.CursorQueryParam{
			Cursor: "xxoi21312300-0asdx",
			Limit:  10,
		})

		var appErr *apperr.Error
		require.ErrorAs(t, err, &appErr)
		require.Equal(t, "invalid cursor", appErr.Message)
		require.Nil(t, payload)
	})
}

func newFeedService(t *testing.T) (*testutil.Env, feed.Service) {
	t.Helper()

	env := testutil.NewEnv(t)
	store := store.New(env.DB)
	svc := feed.NewService(store)

	return env, svc
}
