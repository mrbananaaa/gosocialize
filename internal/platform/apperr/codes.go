package apperr

var Code = struct {
	Internal           string
	ValidationFailed   string
	UserNotFound       string
	InvalidCredentials string
	Forbidden          string
	Conflict           string
}{
	Internal:           "internal.server_error",
	ValidationFailed:   "validation.failed",
	UserNotFound:       "user.not_found",
	InvalidCredentials: "auth.invalid_credentials",
	Forbidden:          "auth.forbidden",
	Conflict:           "resource.conflict",
}
