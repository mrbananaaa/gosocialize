package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres/sqlc"
	"github.com/mrbananaaa/gosocialize/internal/user"
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
) (*user.User, error) {
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

	return &user.User{
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

type LoginPayload struct {
	RefreshToken string
	AccessToken  string
}

func (s *Service) Login(
	ctx context.Context,
	input LoginInput,
) (*LoginPayload, error) {
	u, err := s.q.GetUserByUsername(ctx, input.Username)
	if err != nil {
		err = postgres.PgErrMapper(err)

		if errors.Is(err, apperr.ErrNotFound) {
			return nil, apperr.New(
				apperr.Code.Unauthorized,
				"invalid username/password",
			)
		}

		return nil, err
	}

	match, err := s.hasher.Compare(input.Password, u.Password)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, apperr.New(
			apperr.Code.Unauthorized,
			"invalid username/password",
		)
	}

	accessToken, err := s.token.GenerateAccess(u.ID.String(), "user")
	if err != nil {
		return nil, apperr.ErrInternal
	}

	refreshToken, err := s.token.GenerateRefresh(u.ID.String())
	if err != nil {
		return nil, apperr.ErrInternal
	}

	return &LoginPayload{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}

type RefreshInput struct {
	RefreshToken string
}

type RefreshOutput struct {
	AccessToken string
}

func (s *Service) Refresh(
	ctx context.Context,
	input RefreshInput,
) (*RefreshOutput, error) {
	claims, err := s.token.Verify(input.RefreshToken)
	if err != nil {
		return nil, apperr.New(
			apperr.Code.Unauthorized,
			"invalid refresh token",
		)
	}

	accessToken, err := s.token.GenerateAccess(claims.UserID.String(), "user")
	if err != nil {
		return nil, apperr.New(
			apperr.Code.Internal,
			"internal server error",
		)
	}

	return &RefreshOutput{
		AccessToken: accessToken,
	}, nil
}
