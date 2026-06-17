package post_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/mrbananaaa/gosocialize/internal/post"
	"github.com/mrbananaaa/gosocialize/internal/testutil"
	"github.com/mrbananaaa/gosocialize/store"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	env := testutil.NewEnv(t)
	store := store.New(env.DB)

	svc := post.NewService(store)

	t.Run("success - user can create post", func(t *testing.T) {
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
}
