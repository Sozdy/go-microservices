// Package errs — доменные ошибки сервиса inventory: тип Error с транспортно-независимой
// категорией Code и конструкторы под каждую категорию.
//
// Идея (см. Ben Johnson, "Failure Is Your Domain"): доменный/сервисный слой присваивает
// ошибке категорию, delivery-слой (gRPC-интерцептор) один раз маппит её в транспортный код
// и логирует. Технические/неклассифицированные ошибки в errs.Error НЕ оборачивают —
// CodeOf автоматически считает их CodeInternal. Зависимостей на grpc/http тут нет.
package errs

import "errors"

// Code — транспортно-независимая категория ошибки.
type Code string

const (
	CodeInternal           Code = "internal"
	CodeNotFound           Code = "not_found"
	CodeInvalidArgument    Code = "invalid_argument"
	CodeConflict           Code = "conflict"
	CodeFailedPrecondition Code = "failed_precondition"
	CodeUnavailable        Code = "unavailable"
)

// internalClientMessage отдаётся клиенту вместо настоящего текста, чтобы наружу
// не утекли детали БД/стека.
const internalClientMessage = "внутренняя ошибка сервера"

// Error — доменная ошибка с категорией. Message предназначен клиенту.
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

// CodeOf достаёт категорию ошибки (раскручивая %w-обёртки через errors.As).
// Не доменная ошибка → CodeInternal; nil → "".
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

// ClientMessage возвращает текст, безопасный для клиента: Message доменной ошибки,
// а для CodeInternal и не-доменных ошибок — фиксированную строку.
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
