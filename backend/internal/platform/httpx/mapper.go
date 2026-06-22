package httpx

import (
	"net/http"

	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
)

func resolveMessage(err *apperr.Error) string {
	if err.Code == apperr.Code.Internal {
		return "something went wrong"
	}

	if err.Message == "" {
		return "something went wrong"
	}

	return err.Message
}

func resolveStatus(code string) int {
	switch code {
	case apperr.Code.ValidationFailed:
		return http.StatusBadRequest
	case apperr.Code.BadRequest:
		return http.StatusBadRequest

	case apperr.Code.UserNotFound:
		return http.StatusNotFound
	case apperr.Code.NotFound:
		return http.StatusNotFound

	case apperr.Code.Conflict:
		return http.StatusConflict

	case apperr.Code.InvalidCredentials:
		return http.StatusUnauthorized
	case apperr.Code.Unauthorized:
		return http.StatusUnauthorized
	case apperr.Code.Forbidden:
		return http.StatusForbidden

	default:
		return http.StatusInternalServerError
	}
}
