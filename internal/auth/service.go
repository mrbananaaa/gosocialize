package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/domain"
	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres/sqlc"
)

type Service struct {
	token  TokenService
	hasher PasswordHasher
	q      *sqlc.Queries
}

func NewService(
	token TokenService,
	hasher PasswordHasher,
	q *sqlc.Queries,
) *Service {
	return &Service{
		token:  token,
		hasher: hasher,
		q:      q,
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
) (*domain.User, error) {
	userID := uuid.New()
	creationTime := time.Now()
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		// TODO: Wrap with internal error
		return nil, err
	}

	u := sqlc.CreateUserParams{
		ID:        userID,
		Email:     input.Email,
		Username:  input.Username,
		Password:  passwordHash,
		Name:      input.Name,
		CreatedAt: creationTime,
		UpdatedAt: creationTime,
	}

	err = s.q.CreateUser(ctx, u)
	if err != nil {
		err = postgres.PgErrMapper(err)
		return nil, err
	}

	return &domain.User{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		Password:  u.Password,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

type LoginInput struct {
	Username string
	Password string
}

func (s *Service) Login(
	ctx context.Context,
	input LoginInput,
) (string, error) {
	u, err := s.q.GetUserByUsername(ctx, input.Username)
	if err != nil {
		err = postgres.PgErrMapper(err)

		if errors.Is(err, apperr.ErrNotFound) {
			return "", apperr.New(
				apperr.Code.Unauthorized,
				"invalid username/password",
			)
		}

		return "", err
	}

	match, err := s.hasher.Compare(input.Password, u.Password)
	if err != nil {
		return "", err
	}

	if !match {
		return "", apperr.New(
			apperr.Code.Unauthorized,
			"invalid username/password",
		)
	}

	token, err := s.token.Generate(u.ID.String(), "user")
	if err != nil {
		return "", apperr.ErrInternal
	}

	return token, nil
}
