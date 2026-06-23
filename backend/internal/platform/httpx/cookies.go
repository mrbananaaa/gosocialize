package httpx

import (
	"net/http"

	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
)

func SetRefreshToken(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/v1/auth/refresh",
	})
}

func CleanRefreshToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
		Path:     "/v1/auth/refresh",
	})
}

func GetRefreshToken(r *http.Request) (string, error) {
	token, err := r.Cookie("refresh_token")
	if err != nil {
		return "", apperr.New(
			apperr.Code.Unauthorized,
			"invalid refresh token",
		)
	}

	return token.Value, nil
}
