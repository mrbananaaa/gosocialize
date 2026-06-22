package seeds

import (
	"context"
	"errors"
	"math/rand/v2"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/platform/db"
	"github.com/mrbananaaa/gosocialize/store"
)

func SeedPosts(
	ctx context.Context,
	s *store.Store,
	count int,
	users []db.User,
) ([]db.Post, error) {
	posts := make([]db.Post, 0, count)

	if len(users) == 0 {
		return nil, errors.New("Users ids are empty")
	}

	for range count {
		now := time.Now()
		randomUser := rand.IntN(len(users))
		userID := users[randomUser].ID

		p := db.CreatePostParams{
			ID:        uuid.New(),
			AuthorID:  userID,
			Title:     gofakeit.LoremIpsumSentence(16),
			Content:   gofakeit.LoremIpsumParagraph(10, 20, 15, " "),
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := s.Q.CreatePost(ctx, p)
		if err != nil {
			return nil, err
		}

		posts = append(posts, db.Post(p))
	}

	return posts, nil
}
