package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

// TODO: implementations goes here, e.g., PostgresRepository

type PostgresRepository struct {
}

func NewRepository() *PostgresRepository {
	return &PostgresRepository{}
}

func (r *PostgresRepository) Create(ctx context.Context, user *User) error {
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return nil, nil
}

func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}
