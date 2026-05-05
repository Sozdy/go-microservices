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

// POST /api/v1/orders/{order_uuid}/pay.
func (a *api) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	if errDetails, ok := validatePayOrderRequest(req, params); !ok {
		return &orderv1.PayOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "ошибка валидации запроса",
			Details: errDetails,
		}, nil
	}

	out, err := a.orderService.PayOrder(ctx, converter.PayOrderInFromRequest(req, params))
	if err != nil {
		return handlePayOrderError(err), nil
	}

	return converter.PayOrderOutToResponse(out), nil
}

func validatePayOrderRequest(_ *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (errDetails []orderv1.ErrorDetail, ok bool) {
	if params.OrderUUID == uuid.Nil {
		errDetails = append(errDetails, orderv1.ErrorDetail{
			Message: "Поле order_uuid не может быть uuid null формата",
			Field:   "order_uuid",
		})
	}

	ok = len(errDetails) == 0

	return errDetails, ok
}

func handlePayOrderError(err error) orderv1.PayOrderRes {
	orderErr := apierr.FromError(err)
	if orderErr.Log {
		slog.Error(orderErr.Message, "err", err)
	}
	switch orderErr.Status {
	case http.StatusBadRequest:
		return &orderv1.PayOrderBadRequest{
			Code:    orderErr.Status,
			Message: orderErr.Message,
		}
	case http.StatusNotFound:
		return &orderv1.PayOrderNotFound{
			Code:    orderErr.Status,
			Message: orderErr.Message,
		}
	case http.StatusConflict:
		return &orderv1.PayOrderConflict{
			Code:    orderErr.Status,
			Message: orderErr.Message,
		}
	case http.StatusServiceUnavailable:
		return &orderv1.PayOrderInternalServerError{
			Code:    orderErr.Status,
			Message: orderErr.Message,
		}
	default:
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: orderErr.Message,
		}
	}
}
