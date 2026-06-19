package token

import "github.com/google/uuid"

type TokenService interface {
	GenerateAccess(userID string, role string) (string, error)
	GenerateRefresh(userID string) (string, error)
	Verify(token string) (*Claims, error)
}

type Claims struct {
	UserID uuid.UUID
	Role   string
}
