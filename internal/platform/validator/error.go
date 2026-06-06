package validator

import "github.com/go-playground/validator/v10"

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationError struct {
	Errors []FieldError `json:"errors"`
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

func NewValidationError(errs validator.ValidationErrors) *ValidationError {
	var result []FieldError

	for _, e := range errs {
		result = append(result, FieldError{
			Field:   e.Field(),
			Message: messageForTag(e),
		})
	}

	return &ValidationError{Errors: result}
}
