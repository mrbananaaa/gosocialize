package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

		body, err := json.Marshal(user.FollowUserRequest{
			TargetUserID: followee.ID,
		})
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/user/%s/follow", follower.ID.String()),
			bytes.NewReader(body),
		)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("fail - validation error", func(t *testing.T) {
		follower := env.CreateUser(t)

		body, err := json.Marshal(user.FollowUserRequest{
			TargetUserID: uuid.Nil,
		})
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/user/%s/follow", follower.ID.String()),
			bytes.NewReader(body),
		)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestUnfollowUserHandler(t *testing.T) {
	env, handler := newUserHandler(t)

	router := chi.NewRouter()
	router.Delete("/user/{userID}/unfollow", handler.UnfollowUser)

	t.Run("success - user unfollowed", func(t *testing.T) {
		followee := env.CreateUser(t)
		follower := env.CreateUser(t)

		body, err := json.Marshal(user.UnFollowUserRequest{
			TargetUserID: followee.ID,
		})
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/user/%s/follow", follower.ID.String()),
			bytes.NewReader(body),
		)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusCreated, rr.Code)
	})

}

func newUserHandler(t *testing.T) (*testutil.Env, *user.Handler) {
	t.Helper()

	env := testutil.NewEnv(t)
	store := store.New(env.DB)
	userSvc := user.NewService(store)
	followSvc := follow.NewService(store)

	handler := user.NewHandler(userSvc, followSvc)

	return env, handler
}
