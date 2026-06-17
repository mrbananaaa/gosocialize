package testutil

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/platform/db"
	"github.com/stretchr/testify/require"
)

func (e *Env) CreatePost(
	t *testing.T,
	userID uuid.UUID,
) db.Post {
	t.Helper()

	param := db.CreatePostParams{
		ID:        uuid.New(),
		AuthorID:  userID,
		Title:     gofakeit.LoremIpsumSentence(16),
		Content:   gofakeit.LoremIpsumParagraph(10, 20, 15, " "),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := e.Q.CreatePost(t.Context(), param)
	require.NoError(t, err)

	return db.Post(param)
}
