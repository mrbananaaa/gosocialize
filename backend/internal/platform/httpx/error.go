package httpx

import (
	"errors"
	"net/http"

	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
	"github.com/mrbananaaa/gosocialize/internal/platform/validator"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func Error(
	w http.ResponseWriter,
	err error,
) {
	var validationError *validator.ValidationError

	if errors.As(err, &validationError) {
		writeJSON(
			w,
			http.StatusBadRequest,
			ErrorResponse{
				Success: false,
				Code:    apperr.Code.ValidationFailed,
				Message: "validation failed",
				Details: validationError.Errors,
			},
		)
		return
	}

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

	writeJSON(
		w,
		resolveStatus(appErr.Code),
		ErrorResponse{
			Success: false,
			Code:    appErr.Code,
			Message: resolveMessage(appErr),
		},
	)
}
