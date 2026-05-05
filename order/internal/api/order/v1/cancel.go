package v1

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	apierr "github.com/Sozdy/go-microservices/order/internal/api"
	orderv1 "github.com/Sozdy/go-microservices/shared/pkg/openapi/order/v1"
)

// POST /api/v1/orders/{order_uuid}/cancel.
func (a *api) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	if errDetails, ok := validateCancelOrderRequest(params); !ok {
		return &orderv1.CancelOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "ошибка валидации запроса",
			Details: errDetails,
		}, nil
	}

	if _, err := a.orderService.CancelOrder(ctx, params.OrderUUID); err != nil {
		return handleCancelOrderError(err), nil
	}

	return &orderv1.CancelOrderResponse{}, nil
}

func validateCancelOrderRequest(params orderv1.CancelOrderParams) (errDetails []orderv1.ErrorDetail, ok bool) {
	if params.OrderUUID == uuid.Nil {
		errDetails = append(errDetails, orderv1.ErrorDetail{
			Message: "Поле order_uuid не может быть uuid null формата",
			Field:   "order_uuid",
		})
	}

	ok = len(errDetails) == 0

	return errDetails, ok
}

func handleCancelOrderError(err error) orderv1.CancelOrderRes {
	orderErr := apierr.FromError(err)
	if orderErr.Log {
		slog.Error(orderErr.Message, "err", err)
	}
	switch orderErr.Status {
	case http.StatusNotFound:
		return &orderv1.CancelOrderNotFound{
			Code:    orderErr.Status,
			Message: orderErr.Message,
		}
	case http.StatusConflict:
		return &orderv1.CancelOrderConflict{
			Code:    orderErr.Status,
			Message: orderErr.Message,
		}
	default:
		return &orderv1.CancelOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: orderErr.Message,
		}
	}
}
