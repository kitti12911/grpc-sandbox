package apperror

import "fmt"

type Code int

const (
	CodeInternal Code = iota
	CodeNotFound
	CodeAlreadyExist
	CodeInvalidInput
	CodeUnauthorized
	CodeForbidden
)

type Error struct {
	code    Code
	message string
	err     error
}

func New(code Code, message string, err error) *Error {
	return &Error{code: code, message: message, err: err}
}

func (e *Error) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.message, e.err)
	}

	return e.message
}

func (e *Error) Code() Code      { return e.code }
func (e *Error) Message() string { return e.message }
func (e *Error) Unwrap() error   { return e.err }

func Internal(message string, err error) *Error {
	return New(CodeInternal, message, err)
}

func NotFound(message string, err error) *Error {
	return New(CodeNotFound, message, err)
}

func AlreadyExist(message string, err error) *Error {
	return New(CodeAlreadyExist, message, err)
}

func InvalidInput(message string, err error) *Error {
	return New(CodeInvalidInput, message, err)
}

func Unauthorized(message string, err error) *Error {
	return New(CodeUnauthorized, message, err)
}

func Forbidden(message string, err error) *Error {
	return New(CodeForbidden, message, err)
}
