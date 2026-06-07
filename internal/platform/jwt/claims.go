package jwt

import "github.com/golang-jwt/jwt/v5"

type jwtClaims struct {
	UserID string `json:"sub"`
	Role   string `json:"role"`

	jwt.RegisteredClaims
}
