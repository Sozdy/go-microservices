package errs

import "errors"

var (
	ErrOrderUUIDEmpty           = errors.New("order_uuid не может быть пустым")
	ErrInvalidOrderUUID         = errors.New("неверный формат uuid заказа")
	ErrPaymentMethodUnspecified = errors.New("payment_method не может быть UNSPECIFIED")
)
