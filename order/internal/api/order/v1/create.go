package v1

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	apierr "github.com/Sozdy/go-microservices/order/internal/api"
	"github.com/Sozdy/go-microservices/order/internal/api/order/v1/converter"
	orderv1 "github.com/Sozdy/go-microservices/shared/pkg/openapi/order/v1"
)

// CreateOrder POST /api/v1/orders.
func (a *api) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	if errDetails, ok := validateCreateOrderRequest(req); !ok {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "ошибка валидации запроса",
			Details: errDetails,
		}, nil
	}

	out, err := a.orderService.CreateOrder(ctx, converter.CreateOrderInFromRequest(req))
	if err != nil {
		return handleCreateOrderError(err), nil
	}

	return converter.CreateOrderOutToResponse(out), nil
}

func validateCreateOrderRequest(req *orderv1.CreateOrderRequest) (errDetails []orderv1.ErrorDetail, ok bool) {
	if req.GetHullUUID() == uuid.Nil {
		errDetails = append(errDetails, orderv1.ErrorDetail{
			Message: "Поле hull_uuid не может быть uuid null формата",
			Field:   "hull_uuid",
		})
	}
	if req.GetEngineUUID() == uuid.Nil {
		errDetails = append(errDetails, orderv1.ErrorDetail{
			Message: "Поле engine_uuid не может быть uuid null формата",
			Field:   "engine_uuid",
		})
	}

	ok = len(errDetails) == 0

	return errDetails, ok
}

func handleCreateOrderError(err error) orderv1.CreateOrderRes {
	orderErr := apierr.FromError(err)
	if orderErr.Log {
		slog.Error(orderErr.Message, "err", err)
	}
	switch orderErr.Status {
	case http.StatusBadRequest:
		return &orderv1.CreateOrderBadRequest{
			Code:    orderErr.Status,
			Message: orderErr.Message,
		}
	case http.StatusNotFound:
		return &orderv1.CreateOrderNotFound{
			Code:    orderErr.Status,
			Message: orderErr.Message,
		}
	case http.StatusConflict:
		return &orderv1.CreateOrderConflict{
			Code:    orderErr.Status,
			Message: orderErr.Message,
		}
	case http.StatusServiceUnavailable:
		return &orderv1.CreateOrderInternalServerError{
			Code:    orderErr.Status,
			Message: orderErr.Message,
		}
	default:
		return &orderv1.CreateOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: orderErr.Message,
		}
	}
}
