package errors

import "fmt"

type Code string

const (
	Invalid  Code = "invalid_argument"
	NotFound Code = "not_found"
	Conflict Code = "conflict"
	Timeout  Code = "timeout"
	Internal Code = "internal"
)

type Error struct {
	Code    Code
	Message string
	Cause   error
}

func (e *Error) Error() string { return string(e.Code) + ": " + e.Message }
func (e *Error) Unwrap() error { return e.Cause }
func New(code Code, message string, cause error) error {
	return &Error{Code: code, Message: message, Cause: cause}
}
func Wrap(code Code, cause error) error {
	if cause == nil {
		return nil
	}
	return &Error{Code: code, Message: cause.Error(), Cause: cause}
}
func Ensure(err error, code Code) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", code, err)
}
