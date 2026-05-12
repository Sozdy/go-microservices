package errs

var (
	ErrOrderUUIDEmpty           = InvalidArgument("order_uuid не может быть пустым")
	ErrInvalidOrderUUID         = InvalidArgument("неверный формат uuid заказа")
	ErrPaymentMethodUnspecified = InvalidArgument("payment_method не может быть UNSPECIFIED")
)
