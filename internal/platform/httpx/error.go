package httpx

import (
	"errors"
	"net/http"

	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Error(
	w http.ResponseWriter,
	err error,
) {
	var appErr *apperr.Error
	if !errors.As(err, &appErr) {
		writeJSON(
			w,
			http.StatusInternalServerError,
			ErrorResponse{
				Success: false,
				Code:    apperr.Code.Internal,
				Message: "something went wrong",
			},
		)
		return
	}

	message := appErr.Message
	if message == "" {
		message = "something went wrong"
	}

	writeJSON(
		w,
		statusFromCode(appErr.Code),
		ErrorResponse{
			Success: false,
			Code:    appErr.Code,
			Message: message,
		},
	)
}

func statusFromCode(code string) int {
	switch code {
	case apperr.Code.ValidationFailed:
		return http.StatusBadRequest

	case apperr.Code.InvalidCredentials:
		return http.StatusUnauthorized

	case apperr.Code.Forbidden:
		return http.StatusForbidden

	case apperr.Code.UserNotFound:
		return http.StatusNotFound

	case apperr.Code.Conflict:
		return http.StatusConflict

	default:
		// WARN: remove this log later
		// logger.Warn("[httpx] - Unknown Error, defaulting to Internal Server Error")
		return http.StatusInternalServerError
	}
}
