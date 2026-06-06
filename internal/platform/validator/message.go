package validator

import "github.com/go-playground/validator/v10"

func messageForTag(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return "is too short"
	case "max":
		return "is too long"
	case "gte":
		return "must be greater or equal to " + e.Param()
	}

	return "is invalid"
}
