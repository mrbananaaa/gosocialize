package user_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/mrbananaaa/gosocialize/internal/feed"
	"github.com/mrbananaaa/gosocialize/internal/follow"
	"github.com/mrbananaaa/gosocialize/internal/testutil"
	"github.com/mrbananaaa/gosocialize/internal/user"
	"github.com/mrbananaaa/gosocialize/store"
	"github.com/stretchr/testify/require"
)

func TestFollowUserHandler(t *testing.T) {
	env, handler := newUserHandler(t)

	router := chi.NewRouter()
	router.Get("/user/{userID}/follow", handler.FollowUser)

	t.Run("success - user followed", func(t *testing.T) {
		followee := env.CreateUser(t)
		follower := env.CreateUser(t)

		req := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/user/%s/follow", followee.ID.String()),
			nil,
		)

		req = testutil.AuthenticatedRequest(
			req,
			follower.ID,
		)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusCreated, rr.Code)
	})
}

func TestUnfollowUserHandler(t *testing.T) {
	env, handler := newUserHandler(t)

	router := chi.NewRouter()
	router.Delete("/user/{userID}/unfollow", handler.UnfollowUser)

	t.Run("success - user unfollowed", func(t *testing.T) {
		followee := env.CreateUser(t)
		follower := env.CreateUser(t)
		env.CreateFollow(t, follower.ID, followee.ID)

		req := httptest.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/user/%s/unfollow", followee.ID.String()),
			nil,
		)

		req = testutil.AuthenticatedRequest(
			req,
			follower.ID,
		)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestUserFeedsHandler(t *testing.T) {
	env, handler := newUserHandler(t)

	router := chi.NewRouter()
	router.Get("/user/feed", handler.UserFeeds)

	t.Run("success - get user feeds", func(t *testing.T) {
		currentUser := env.CreateUser(t)
		postsCount := 20

		for range postsCount {
			u := env.CreateUser(t)
			_ = env.CreatePost(t, u.ID)
			env.CreateFollow(t, currentUser.ID, u.ID)
		}

		req := httptest.NewRequest(
			http.MethodGet,
			"/user/feed",
			nil,
		)

		req = testutil.AuthenticatedRequest(
			req,
			currentUser.ID,
		)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("fail - unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodGet,
			"/user/feed",
			nil,
		)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}

func newUserHandler(t *testing.T) (*testutil.Env, *user.Handler) {
	t.Helper()

	env := testutil.NewEnv(t)
	store := store.New(env.DB)
	userSvc := user.NewService(store)
	followSvc := follow.NewService(store)
	feedSvc := feed.NewService(store)

	handler := user.NewHandler(userSvc, followSvc, feedSvc)

	return env, handler
}
