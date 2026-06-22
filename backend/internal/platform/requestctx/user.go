package requestctx

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/token"
)

type User struct {
	ID   uuid.UUID
	Role string
}

func UserFromClaims(c *token.Claims) User {
	return User{
		ID:   c.UserID,
		Role: c.Role,
	}
}

func WithUser(ctx context.Context, u User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

func UserFromContext(ctx context.Context) (User, bool) {
	u, ok := ctx.Value(userKey).(User)
	return u, ok
}
