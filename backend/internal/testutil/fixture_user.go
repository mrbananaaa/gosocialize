package testutil

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/platform/db"
	"github.com/stretchr/testify/require"
)

type UserOption func(*db.CreateUserParams)

func (e *Env) CreateUser(t *testing.T, opts ...UserOption) db.User {
	t.Helper()

	param := db.CreateUserParams{
		ID:        uuid.New(),
		Email:     gofakeit.Email(),
		Username:  gofakeit.Username(),
		Password:  gofakeit.Password(true, true, true, true, false, 8),
		Name:      gofakeit.Name(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, opt := range opts {
		opt(&param)
	}

	err := e.Q.CreateUser(t.Context(), param)
	require.NoError(t, err)

	return db.User(param)
}

func WithEmail(username string) UserOption {
	return func(p *db.CreateUserParams) {
		p.Username = username
	}
}

func WithUsername(email string) UserOption {
	return func(p *db.CreateUserParams) {
		p.Email = email
	}
}

func WithPassword(password string) UserOption {
	return func(p *db.CreateUserParams) {
		p.Password = password
	}
}
