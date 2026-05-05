package inventory

import (
	"errors"

	"google.golang.org/grpc/codes"

	errs "github.com/Sozdy/go-microservices/inventory/internal/errors"
)

type Error struct {
	Code    codes.Code
	Message string
	Log     bool
}

func FromError(err error) Error {
	switch {
	case errors.Is(err, errs.ErrPartNotFound):
		return Error{
			Code:    codes.NotFound,
			Message: "одна из деталей не найдена",
		}

	case errors.Is(err, errs.ErrInvalidUUID):
		return Error{
			Code:    codes.InvalidArgument,
			Message: "неверный формат uuid",
		}

	default:
		return Error{
			Code:    codes.Internal,
			Message: "внутренняя ошибка сервера",
			Log:     true,
		}
	}
}
