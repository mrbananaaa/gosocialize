package middlewares

import (
	"net/http"
	"strings"

	"github.com/mrbananaaa/gosocialize/internal/auth"
	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
	"github.com/mrbananaaa/gosocialize/internal/platform/httpx"
	"github.com/mrbananaaa/gosocialize/internal/platform/requestctx"
)

type AuthMiddleware struct {
	token auth.TokenService
}

func NewAuth(
	token auth.TokenService,
) *AuthMiddleware {
	return &AuthMiddleware{
		token: token,
	}
}

func (a *AuthMiddleware) WithAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			httpx.Error(w, apperr.New(
				apperr.Code.Unauthorized,
				"missing token",
			))
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := a.token.Verify(tokenStr)
		if err != nil {
			httpx.Error(w, apperr.New(
				apperr.Code.Unauthorized,
				"invalid token",
			))
			return
		}

		user := requestctx.UserFromClaims(claims)
		ctx := requestctx.WithUser(r.Context(), user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
