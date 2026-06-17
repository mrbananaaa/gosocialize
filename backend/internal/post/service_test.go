package post_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
	"github.com/mrbananaaa/gosocialize/internal/post"
	"github.com/mrbananaaa/gosocialize/internal/testutil"
	"github.com/mrbananaaa/gosocialize/pkg/pagination"
	"github.com/mrbananaaa/gosocialize/store"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	t.Run("success - user can create post", func(t *testing.T) {
		env, svc := newPostService(t)

		user := env.CreateUser(t)

		input := post.CreateInput{
			AuthorID: user.ID,
			Title:    gofakeit.LoremIpsumSentence(16),
			Content:  gofakeit.LoremIpsumParagraph(10, 20, 15, " "),
			Tags:     nil,
		}

		post, err := svc.Create(t.Context(), input)

		require.NoError(t, err)
		require.Equal(t, input.Title, post.Title)
		require.Equal(t, input.Content, post.Content)
		require.Equal(t, user.ID, post.AuthorID)
		require.NotEmpty(t, post.ID)
		require.NotEmpty(t, post.CreatedAt)
		require.NotEmpty(t, post.UpdatedAt)
	})

	t.Run("fail - empty title not allowed", func(t *testing.T) {
		env, svc := newPostService(t)

		user := env.CreateUser(t)

		input := post.CreateInput{
			AuthorID: user.ID,
			Title:    "",
			Content:  gofakeit.LoremIpsumParagraph(10, 20, 15, " "),
			Tags:     nil,
		}

		_, err := svc.Create(t.Context(), input)

		require.Error(t, err)
	})

	t.Run("fail - user does not exists", func(t *testing.T) {
		_, svc := newPostService(t)

		input := post.CreateInput{
			AuthorID: uuid.New(),
			Title:    "Hello",
			Content:  "Gawd dayumn",
		}

		_, err := svc.Create(t.Context(), input)

		require.Error(t, err)
	})
}

func TestListPosts(t *testing.T) {
	t.Run("success - retrieve post list", func(t *testing.T) {
		env, svc := newPostService(t)

		user := env.CreateUser(t)
		_ = env.CreatePost(t, user.ID)

		p, err := svc.ListPosts(t.Context(), pagination.CursorQueryParam{
			Cursor: "",
			Limit:  10,
		})

		require.NoError(t, err)
		require.NotEmpty(t, p)
		require.Len(t, p.Posts, 1)
	})

	t.Run("success - valid cursor", func(t *testing.T) {
		env, svc := newPostService(t)

		for range 50 {
			user := env.CreateUser(t)
			_ = env.CreatePost(t, user.ID)
		}

		pf, err := svc.ListPosts(t.Context(), pagination.CursorQueryParam{
			Cursor: "",
			Limit:  10,
		})
		require.NoError(t, err)
		require.NotNil(t, pf.Posts)
		require.Len(t, pf.Posts, 10)
		require.True(t, pf.HasMore)

		pn, err := svc.ListPosts(t.Context(), pagination.CursorQueryParam{
			Cursor: pf.NextCursor,
			Limit:  10,
		})
		require.NoError(t, err)
		require.NotNil(t, pf.Posts)
		require.Len(t, pn.Posts, 10)
		require.True(t, pn.HasMore)
	})

	t.Run("fail - invalid cursor", func(t *testing.T) {
		_, svc := newPostService(t)

		p, err := svc.ListPosts(t.Context(), pagination.CursorQueryParam{
			Cursor: "xxaeuui213213xx",
			Limit:  10,
		})

		require.Error(t, err)
		require.Nil(t, p)
	})
}

func TestGetByID(t *testing.T) {
	t.Run("fail - no post with given id", func(t *testing.T) {
		_, svc := newPostService(t)

		p, err := svc.GetByID(t.Context(), uuid.New())

		require.Error(t, err)
		require.Nil(t, p)
	})

	t.Run("success - id matched", func(t *testing.T) {
		env, svc := newPostService(t)

		user := env.CreateUser(t)
		post := env.CreatePost(t, user.ID)

		p, err := svc.GetByID(t.Context(), post.ID)
		require.NoError(t, err)
		require.NotNil(t, p)
		require.Equal(t, post.ID, p.ID)
	})
}

func TestDeletePost(t *testing.T) {
	t.Run("success - post deleted", func(t *testing.T) {
		env, svc := newPostService(t)

		user := env.CreateUser(t)
		p := env.CreatePost(t, user.ID)

		err := svc.DeletePost(t.Context(), post.DeletePostInput{
			ID:       p.ID,
			AuthorID: p.AuthorID,
		})
		require.NoError(t, err)

		_, err = svc.GetByID(t.Context(), p.ID)
		require.Error(t, err)
	})

	t.Run("fail - no post with given id", func(t *testing.T) {
		_, svc := newPostService(t)

		err := svc.DeletePost(context.Background(), post.DeletePostInput{
			ID:       uuid.New(),
			AuthorID: uuid.New(),
		})
		require.Error(t, err)
	})

	t.Run("fail - cannot delete other user post", func(t *testing.T) {
		env, svc := newPostService(t)

		user := env.CreateUser(t)
		p := env.CreatePost(t, user.ID)

		err := svc.DeletePost(t.Context(), post.DeletePostInput{
			ID:       p.ID,
			AuthorID: uuid.New(),
		})
		require.Error(t, err)
		require.ErrorIs(t, err, apperr.ErrForbidden)
	})
}

func TestUpdatePost(t *testing.T) {
	t.Run("success - post updated", func(t *testing.T) {
		env, svc := newPostService(t)

		u := env.CreateUser(t)
		p := env.CreatePost(t, u.ID)

		input := post.UpdatePostInput{
			ID:      p.ID,
			UserID:  u.ID,
			Title:   "this is the new title",
			Content: "this is the new content",
		}

		err := svc.UpdatePost(t.Context(), input)
		require.NoError(t, err)

		up, err := svc.GetByID(t.Context(), p.ID)
		require.NoError(t, err)
		require.Equal(t, input.Title, up.Title)
		require.Equal(t, input.Content, up.Content)
		require.NotEqual(t, p.UpdatedAt, up.UpdatedAt)
	})

	t.Run("success - updated only title field", func(t *testing.T) {
		env, svc := newPostService(t)

		u := env.CreateUser(t)
		p := env.CreatePost(t, u.ID)

		input := post.UpdatePostInput{
			ID:      p.ID,
			UserID:  u.ID,
			Title:   "new updated title",
			Content: "",
		}

		err := svc.UpdatePost(t.Context(), input)
		require.NoError(t, err)

		up, err := svc.GetByID(t.Context(), p.ID)
		require.NoError(t, err)
		require.Equal(t, input.Title, up.Title)
		require.Equal(t, p.Content, up.Content)
	})

	t.Run("fail - no post with given id", func(t *testing.T) {
		_, svc := newPostService(t)

		input := post.UpdatePostInput{
			ID:      uuid.New(),
			UserID:  uuid.New(),
			Title:   "updated title",
			Content: "updated content",
		}

		err := svc.UpdatePost(t.Context(), input)
		require.Error(t, err)
	})

	t.Run("fail - cannot update other user post", func(t *testing.T) {
		env, svc := newPostService(t)

		u := env.CreateUser(t)
		p := env.CreatePost(t, u.ID)

		input := post.UpdatePostInput{
			ID:      p.ID,
			UserID:  uuid.New(),
			Title:   "newly updated title btw",
			Content: "",
		}

		err := svc.UpdatePost(t.Context(), input)
		require.Error(t, err)
		require.ErrorIs(t, err, apperr.ErrForbidden)
	})
}

func newPostService(t *testing.T) (*testutil.Env, *post.Service) {
	t.Helper()

	env := testutil.NewEnv(t)
	store := store.New(env.DB)
	svc := post.NewService(store)

	return env, svc
}
