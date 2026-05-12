package errs

var (
	ErrPartNotFound       = NotFound("одна из деталей не найдена")
	ErrOrderNotFound      = NotFound("заказ не найден")
	ErrOrderDoesNotExists = NotFound("заказ не существует")
	ErrPayOrderNotFound   = NotFound("метод для оплаты не найден")

	ErrInvalidPartFilter = InvalidArgument("некорректные параметры запроса деталей")
	ErrInvalidPayOrder   = InvalidArgument("неверные данные для оплаты заказа")

	ErrPartUnavailable           = Conflict("одна из деталей недоступна на складе")
	ErrOrderAlreadyExists        = Conflict("коллизия UUID заказа")
	ErrPayOrderStatusConflict    = FailedPrecondition("статус заказа должен быть PENDING_PAYMENT для оплаты заказа")
	ErrCancelOrderStatusConflict = FailedPrecondition("статус заказа должен быть PENDING_PAYMENT для отмены")

	ErrInventoryUnavailable = Unavailable("сервис inventory временно недоступен")
	ErrPaymentUnavailable   = Unavailable("сервис payment временно недоступен")
)
