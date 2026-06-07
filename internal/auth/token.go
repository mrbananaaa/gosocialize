package auth

import "github.com/google/uuid"

type TokenService interface {
	Generate(userID string, role string) (string, error)
	Verify(token string) (*Claims, error)
}

type Claims struct {
	UserID uuid.UUID
	Role   string
}
