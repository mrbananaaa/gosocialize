package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres"
	"github.com/mrbananaaa/gosocialize/internal/platform/db"
	"github.com/mrbananaaa/gosocialize/internal/token"
	"github.com/mrbananaaa/gosocialize/internal/user"
	"github.com/mrbananaaa/gosocialize/store"
)

type Service struct {
	token  token.TokenService
	hasher PasswordHasher
	store  *store.Store
}

func NewService(
	token token.TokenService,
	hasher PasswordHasher,
	s *store.Store,
) *Service {
	return &Service{
		token:  token,
		hasher: hasher,
		store:  s,
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

	u := db.CreateUserParams{
		ID:        userID,
		Email:     input.Email,
		Username:  input.Username,
		Password:  passwordHash,
		Name:      input.Name,
		CreatedAt: creationTime,
		UpdatedAt: creationTime,
	}

	err = s.store.Q.CreateUser(ctx, u)
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
	u, err := s.store.Q.GetUserByUsername(ctx, input.Username)
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

type RefreshPayload struct {
	AccessToken string
}

func (s *Service) Refresh(
	ctx context.Context,
	input RefreshInput,
) (*RefreshPayload, error) {
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

	return &RefreshPayload{
		AccessToken: accessToken,
	}, nil
}
