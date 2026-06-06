package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres/sqlc"
)

type Service struct {
	ph PasswordHasher
	q  *sqlc.Queries
}

func NewService(
	passwordHasher PasswordHasher,
	q *sqlc.Queries,
) *Service {
	return &Service{
		ph: passwordHasher,
		q:  q,
	}
}

type RegisterInput struct {
	Email    string
	Username string
	Password string
	Name     string
}

func (s *Service) Register(
	ctx context.Context,
	input RegisterInput,
) error {
	userID := uuid.New()
	creationTime := time.Now()
	passwordHash, err := s.ph.Hash(input.Password)
	if err != nil {
		// TODO: Wrap with internal error
		return err
	}

	err = s.q.CreateUser(ctx, sqlc.CreateUserParams{
		ID:        userID,
		Email:     input.Email,
		Username:  input.Username,
		Password:  passwordHash,
		Name:      input.Name,
		CreatedAt: creationTime,
		UpdatedAt: creationTime,
	})
	if err != nil {
		// TODO: Wrap with internal error
		return err
	}

	return nil
}
