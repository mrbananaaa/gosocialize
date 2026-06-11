package apperr

import "fmt"

var (
	ErrConflict     = New(Code.Conflict, "")
	ErrNotFound     = New(Code.NotFound, "")
	ErrForbidden    = New(Code.Forbidden, "")
	ErrUnauthorized = New(Code.Unauthorized, "")
	ErrValidation   = New(Code.ValidationFailed, "")
	ErrInternal     = New(Code.Internal, "something went wrong")
)

type Error struct {
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}

	return e.Code
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func Wrap(err error, code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func DecodeBodyErr(err error) *Error {
	return Wrap(
		err,
		Code.BadRequest,
		"Bad request body",
	)
}

func InvalidUUIDErr(err error, field string) *Error {
	return Wrap(
		err,
		Code.BadRequest,
		fmt.Sprintf("invalid uuid field - %s", field),
	)
}
