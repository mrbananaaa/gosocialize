package httpx

import (
	"net/http"
)

func SetRefreshToken(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		// Path:     "/v1/auth/refresh",
	})
}

func CleanRefreshToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})
}

func GetRefreshToken(r *http.Request) (string, error) {
	token, err := r.Cookie("refresh_token")
	if err != nil {
		return "", err
	}

	return token.Value, nil
}
