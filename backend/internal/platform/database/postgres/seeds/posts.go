package seeds

import (
	"context"
	"errors"
	"math/rand/v2"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres/sqlc"
)

func SeedPosts(
	ctx context.Context,
	q *sqlc.Queries,
	count int,
	users []sqlc.User,
) ([]sqlc.Post, error) {
	var posts []sqlc.Post

	if len(users) == 0 {
		return nil, errors.New("Users ids are empty")
	}

	for range count {
		now := time.Now()
		randomUser := rand.IntN(len(users))
		userID := users[randomUser].ID

		p := sqlc.CreatePostParams{
			ID:        uuid.New(),
			AuthorID:  userID,
			Title:     gofakeit.LoremIpsumSentence(16),
			Content:   gofakeit.LoremIpsumParagraph(10, 20, 15, " "),
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := q.CreatePost(ctx, p)
		if err != nil {
			return nil, err
		}

		posts = append(posts, sqlc.Post(p))
	}

	return posts, nil
}
