package errs

import "errors"

var (
	ErrPartNotFound         = errors.New("деталь не найдена")
	ErrInvalidPartFilter    = errors.New("неверный фильтр поиска деталей")
	ErrInventoryUnavailable = errors.New("сервис inventory недоступен")

	ErrPayOrderNotFound   = errors.New("метод для оплаты не найден")
	ErrInvalidPayOrder    = errors.New("неверные данные для оплаты заказа")
	ErrPaymentUnavailable = errors.New("сервис payment недоступен")

	ErrPartUnavailable           = errors.New("деталь недоступна (нет на складе)")
	ErrCancelOrderStatusConflict = errors.New("статус заказа должен быть PENDING_PAYMENT для отмены")
	ErrPayOrderStatusConflict    = errors.New("статус заказа должен быть PENDING_PAYMENT для оплаты заказа")

	ErrOrderNotFound      = errors.New("заказ не найден")
	ErrOrderAlreadyExists = errors.New("заказ с таким UUID уже существует")
	ErrOrderDoesNotExists = errors.New("заказ не существует")
)
