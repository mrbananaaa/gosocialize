package post_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/platform/httpx"
	"github.com/mrbananaaa/gosocialize/internal/post"
	"github.com/mrbananaaa/gosocialize/internal/testutil"
	"github.com/mrbananaaa/gosocialize/store"
	"github.com/stretchr/testify/require"
)

func TestCreatePostHandler(t *testing.T) {
	t.Run("success - post created", func(t *testing.T) {
		env, handler := newPostHandler(t)

		user := env.CreateUser(t)

		body, err := json.Marshal(post.CreatePostRequest{
			Title:   "hello world",
			Content: "my first post",
		})
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodPost,
			"/posts",
			bytes.NewReader(body),
		)

		req = testutil.AuthenticatedRequest(
			req,
			user.ID,
		)

		rr := httptest.NewRecorder()

		handler.CreatePost(rr, req)

		require.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("fail - unauthorized", func(t *testing.T) {
		_, handler := newPostHandler(t)

		body, err := json.Marshal(post.CreatePostRequest{
			Title:   "hello",
			Content: "world",
		})
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodPost,
			"/posts",
			bytes.NewReader(body),
		)

		rr := httptest.NewRecorder()

		handler.CreatePost(rr, req)

		require.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("fail - invalid json", func(t *testing.T) {
		env, handler := newPostHandler(t)

		user := env.CreateUser(t)

		req := httptest.NewRequest(
			http.MethodPost,
			"/posts",
			bytes.NewBufferString("{invalid-json"),
		)

		req = testutil.AuthenticatedRequest(
			req,
			user.ID,
		)

		rr := httptest.NewRecorder()

		handler.CreatePost(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("fail - validation error", func(t *testing.T) {
		env, handler := newPostHandler(t)

		user := env.CreateUser(t)

		body, err := json.Marshal(post.CreatePostRequest{
			Title: "",
		})
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodPost,
			"/posts",
			bytes.NewReader(body),
		)

		req = testutil.AuthenticatedRequest(
			req,
			user.ID,
		)

		rr := httptest.NewRecorder()

		handler.CreatePost(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestListPostsHandler(t *testing.T) {
	t.Run("success - list posts", func(t *testing.T) {
		_, handler := newPostHandler(t)

		req := httptest.NewRequest(
			http.MethodGet,
			"/post",
			nil,
		)

		rr := httptest.NewRecorder()

		handler.ListPosts(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("success - response format", func(t *testing.T) {
		env, handler := newPostHandler(t)

		for range 11 {
			u := env.CreateUser(t)
			_ = env.CreatePost(t, u.ID)
		}

		req := httptest.NewRequest(
			http.MethodGet,
			"/post?limit=5",
			nil,
		)

		rr := httptest.NewRecorder()

		handler.ListPosts(rr, req)

		require.Equal(t, http.StatusOK, rr.Code, "status code should be 200")

		var resp httpx.Response
		err := json.Unmarshal(
			rr.Body.Bytes(),
			&resp,
		)
		require.NoError(t, err)

		require.True(t, resp.Success, "resp.Success should be true")
		require.Len(t, resp.Data, 5, "length of resp.Data should be 5")
		require.NotEmpty(t, resp.Meta, "resp.Meta should not be empty")
	})

	t.Run("fail - invalid cursor", func(t *testing.T) {
		_, handler := newPostHandler(t)

		req := httptest.NewRequest(
			http.MethodGet,
			"/post?limit=5&cursor=12389xhqieuhqwe8",
			nil,
		)

		rr := httptest.NewRecorder()

		handler.ListPosts(rr, req)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestGetPostByIDHandler(t *testing.T) {
	t.Run("success - get post", func(t *testing.T) {
		env, handler := newPostHandler(t)

		u := env.CreateUser(t)
		p := env.CreatePost(t, u.ID)

		router := chi.NewRouter()
		router.Get("/post/{postID}", handler.GetPostByID)

		req := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/post/%s", p.ID),
			nil,
		)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var resp httpx.Response
		err := json.Unmarshal(
			rr.Body.Bytes(),
			&resp,
		)
		require.NoError(t, err)

		require.True(t, resp.Success)
		require.NotEmpty(t, resp.Data)
	})

	t.Run("fail - invalid post id", func(t *testing.T) {
		_, handler := newPostHandler(t)

		router := chi.NewRouter()
		router.Get("/post/{postID}", handler.GetPostByID)

		req := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/post/%s", "myfaultgng"),
			nil,
		)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("fail - not found", func(t *testing.T) {
		env, handler := newPostHandler(t)

		u := env.CreateUser(t)
		_ = env.CreatePost(t, u.ID)

		router := chi.NewRouter()
		router.Get("/post/{postID}", handler.GetPostByID)

		req := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/post/%s", uuid.New().String()),
			nil,
		)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusNotFound, rr.Code)
	})
}

func TestUpdatePostHandler(t *testing.T) {
	t.Run("success - post updated", func(t *testing.T) {
		env, handler := newPostHandler(t)

		u := env.CreateUser(t)
		p := env.CreatePost(t, u.ID)

		router := chi.NewRouter()
		router.Patch("/post/{postID}", handler.UpdatePost)

		body, err := json.Marshal(post.UpdatePostRequest{
			Title:   "new post title",
			Content: "new post content",
		})
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodPatch,
			fmt.Sprintf("/post/%s", p.ID.String()),
			bytes.NewReader(body),
		)

		req = testutil.AuthenticatedRequest(
			req,
			u.ID,
		)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var resp httpx.Response
		err = json.Unmarshal(
			rr.Body.Bytes(),
			&resp,
		)
		require.NoError(t, err)

		require.True(t, resp.Success)
		require.Equal(t, "post updated!", resp.Message)
	})

	t.Run("fail - unauthorized", func(t *testing.T) {
		env, handler := newPostHandler(t)

		u := env.CreateUser(t)
		_ = env.CreatePost(t, u.ID)

		router := chi.NewRouter()
		router.Patch("/post/{postID}", handler.UpdatePost)

		body, err := json.Marshal(post.UpdatePostRequest{
			Title:   "new post title",
			Content: "new post content",
		})
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodPatch,
			fmt.Sprintf("/post/%s", uuid.New().String()),
			bytes.NewReader(body),
		)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("fail - invalid post id", func(t *testing.T) {
		env, handler := newPostHandler(t)

		u := env.CreateUser(t)

		router := chi.NewRouter()
		router.Patch("/post/{postID}", handler.UpdatePost)

		body, err := json.Marshal(post.UpdatePostRequest{
			Title:   "new post title",
			Content: "new post content",
		})
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodPatch,
			fmt.Sprintf("/post/%s", "awwthisisinvalidpostid"),
			bytes.NewReader(body),
		)

		req = testutil.AuthenticatedRequest(
			req,
			u.ID,
		)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestDeletePostHandler(t *testing.T) {
	t.Run("success - post deleted", func(t *testing.T) {
		env, handler := newPostHandler(t)

		u := env.CreateUser(t)
		p := env.CreatePost(t, u.ID)

		router := chi.NewRouter()
		router.Delete("/post/{postID}", handler.DeletePost)

		req := httptest.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/post/%s", p.ID),
			nil,
		)

		req = testutil.AuthenticatedRequest(
			req,
			u.ID,
		)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var resp httpx.Response
		err := json.Unmarshal(
			rr.Body.Bytes(),
			&resp,
		)
		require.NoError(t, err)

		require.Equal(t, resp.Message, "post deleted!")
	})

	t.Run("fail - unauthorized", func(t *testing.T) {
		env, handler := newPostHandler(t)

		u := env.CreateUser(t)
		p := env.CreatePost(t, u.ID)

		router := chi.NewRouter()
		router.Delete("/post/{postID}", handler.DeletePost)

		req := httptest.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/post/%s", p.ID),
			nil,
		)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("fail - invalid post id", func(t *testing.T) {
		env, handler := newPostHandler(t)

		u := env.CreateUser(t)

		router := chi.NewRouter()
		router.Delete("/post/{postID}", handler.DeletePost)

		req := httptest.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/post/%s", "awwwwwanjrot"),
			nil,
		)

		req = testutil.AuthenticatedRequest(
			req,
			u.ID,
		)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func newPostHandler(t *testing.T) (*testutil.Env, *post.Handler) {
	t.Helper()

	env := testutil.NewEnv(t)
	store := store.New(env.DB)
	svc := post.NewService(store)

	handler := post.NewHandler(svc)

	return env, handler
}
