package errs

import "errors"

type Code string

const (
	CodeInternal           Code = "internal"
	CodeNotFound           Code = "not_found"
	CodeInvalidArgument    Code = "invalid_argument"
	CodeConflict           Code = "conflict"
	CodeFailedPrecondition Code = "failed_precondition"
	CodeUnavailable        Code = "unavailable"
)

const internalClientMessage = "внутренняя ошибка сервера"

type Error struct {
	Code    Code
	Message string
}

func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return string(e.Code)
}

func CodeOf(err error) Code {
	if err == nil {
		return ""
	}

	var e *Error
	if errors.As(err, &e) {
		return e.Code
	}

	return CodeInternal
}

func ClientMessage(err error) string {
	var e *Error
	if errors.As(err, &e) && e.Code != CodeInternal && e.Message != "" {
		return e.Message
	}

	return internalClientMessage
}

func NotFound(message string) error { return &Error{Code: CodeNotFound, Message: message} }
func InvalidArgument(message string) error {
	return &Error{Code: CodeInvalidArgument, Message: message}
}
func Conflict(message string) error { return &Error{Code: CodeConflict, Message: message} }
func FailedPrecondition(message string) error {
	return &Error{Code: CodeFailedPrecondition, Message: message}
}
func Unavailable(message string) error { return &Error{Code: CodeUnavailable, Message: message} }
