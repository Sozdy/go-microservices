package payment

import (
	"errors"

	"google.golang.org/grpc/codes"

	errs "github.com/Sozdy/go-microservices/payment/internal/errors"
)

type Error struct {
	Code    codes.Code
	Message string
	Log     bool
}

func FromError(err error) Error {
	switch {
	case errors.Is(err, errs.ErrOrderUUIDEmpty):
		return Error{
			Code:    codes.InvalidArgument,
			Message: "order_uuid не может быть пустым",
		}

	case errors.Is(err, errs.ErrInvalidOrderUUID):
		return Error{
			Code:    codes.InvalidArgument,
			Message: "неверный формат uuid заказа",
		}

	case errors.Is(err, errs.ErrPaymentMethodUnspecified):
		return Error{
			Code:    codes.InvalidArgument,
			Message: "payment_method не может быть UNSPECIFIED",
		}

	default:
		return Error{
			Code:    codes.Internal,
			Message: "внутренняя ошибка сервера",
			Log:     true,
		}
	}
}
