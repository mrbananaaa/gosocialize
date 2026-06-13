package seeds

import (
	"context"
	"time"
	"unicode/utf8"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/auth"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres/sqlc"
)

func SeedUsers(
	ctx context.Context,
	q *sqlc.Queries,
	count int,
	password string,
) ([]sqlc.User, error) {
	var users []sqlc.User

	hasher := auth.NewArgon2Hasher()
	encodedHash, err := hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	for range count {
		now := time.Now()

		u := sqlc.CreateUserParams{
			ID:        uuid.New(),
			Email:     gofakeit.Email(),
			Username:  truncateString(gofakeit.Username(), 21),
			Password:  encodedHash,
			Name:      gofakeit.Name(),
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := q.CreateUser(ctx, u)
		if err != nil {
			return nil, err
		}

		users = append(users, sqlc.User(u))
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
