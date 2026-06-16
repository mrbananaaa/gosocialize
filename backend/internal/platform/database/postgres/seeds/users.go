package seeds

import (
	"context"
	"time"
	"unicode/utf8"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/auth"
	"github.com/mrbananaaa/gosocialize/internal/platform/db"
	"github.com/mrbananaaa/gosocialize/store"
)

func SeedUsers(
	ctx context.Context,
	s *store.Store,
	count int,
	password string,
) ([]db.User, error) {
	users := make([]db.User, 0, count)

	hasher := auth.NewArgon2Hasher()
	encodedHash, err := hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	for range count {
		now := time.Now()

		u := db.CreateUserParams{
			ID:        uuid.New(),
			Email:     gofakeit.Email(),
			Username:  truncateString(gofakeit.Username(), 21),
			Password:  encodedHash,
			Name:      gofakeit.Name(),
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := s.Q.CreateUser(ctx, u)
		if err != nil {
			return nil, err
		}

		users = append(users, db.User(u))
	}

	return users, nil
}

func truncateString(s string, max int) string {
	if utf8.RuneCountInString(s) <= max {
		return s
	}

	runes := []rune(s)
	return string(runes[:max])
}
