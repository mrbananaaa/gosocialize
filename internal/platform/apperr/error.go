package apperr

type Error struct {
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
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
