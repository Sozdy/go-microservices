package api

import (
	"errors"
	"net/http"

	errs "github.com/Sozdy/go-microservices/order/internal/errors"
)

type Error struct {
	Status  int
	Message string
	Log     bool
}

func FromError(err error) Error {
	switch {
	case errors.Is(err, errs.ErrPartNotFound):
		return Error{
			Status:  http.StatusNotFound,
			Message: "одна из деталей не найдена",
		}
	case errors.Is(err, errs.ErrOrderNotFound):
		return Error{
			Status:  http.StatusNotFound,
			Message: "заказ не найден",
		}
	case errors.Is(err, errs.ErrOrderDoesNotExists):
		return Error{
			Status:  http.StatusNotFound,
			Message: "заказ не существует",
		}
	case errors.Is(err, errs.ErrPayOrderNotFound):
		return Error{
			Status:  http.StatusNotFound,
			Message: "метод для оплаты не найден",
		}

	case errors.Is(err, errs.ErrInvalidPartFilter):
		return Error{
			Status:  http.StatusBadRequest,
			Message: "некорректные параметры запроса деталей",
		}
	case errors.Is(err, errs.ErrInvalidPayOrder):
		return Error{
			Status:  http.StatusBadRequest,
			Message: "неверные данные для оплаты заказа",
		}

	case errors.Is(err, errs.ErrPartUnavailable):
		return Error{
			Status:  http.StatusConflict,
			Message: "одна из деталей недоступна на складе",
		}
	case errors.Is(err, errs.ErrPayOrderStatusConflict):
		return Error{
			Status:  http.StatusConflict,
			Message: "статус заказа должен быть PENDING_PAYMENT для оплаты заказа",
		}
	case errors.Is(err, errs.ErrCancelOrderStatusConflict):
		return Error{
			Status:  http.StatusConflict,
			Message: "статус заказа должен быть PENDING_PAYMENT для отмены",
		}
	case errors.Is(err, errs.ErrOrderAlreadyExists):
		return Error{
			Status:  http.StatusConflict,
			Message: "коллизия UUID заказа",
		}

	case errors.Is(err, errs.ErrInventoryUnavailable):
		return Error{
			Status:  http.StatusServiceUnavailable,
			Message: "сервис inventory временно недоступен",
			Log:     true,
		}
	case errors.Is(err, errs.ErrPaymentUnavailable):
		return Error{
			Status:  http.StatusServiceUnavailable,
			Message: "сервис payment временно недоступен",
			Log:     true,
		}
	default:
		return Error{
			Status:  http.StatusInternalServerError,
			Message: "внутренняя ошибка сервера",
			Log:     true,
		}
	}
}
