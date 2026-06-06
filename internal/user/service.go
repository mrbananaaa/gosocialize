package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres/sqlc"
)

type Service struct {
	q *sqlc.Queries
}

func NewService(q *sqlc.Queries) *Service {
	return &Service{
		q: q,
	}
}

func (s *Service) Register(
	ctx context.Context,
	email string,
	username string,
	password string,
	name string,
) (*User, error) {
	now := time.Now()
	user := &User{
		ID:        uuid.New(),
		Email:     email,
		Username:  username,
		Password:  password,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// if err := s.repo.Create(ctx, user); err != nil {
	// 	return nil, err
	// }

	return user, nil
}
