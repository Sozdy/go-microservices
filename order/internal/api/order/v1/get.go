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

// GetOrder GET /api/v1/orders/{order_uuid}.
func (a *api) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	if errDetails, ok := validateGetOrderRequest(&params); !ok {
		return &orderv1.GetOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "неверные параметры запроса",
			Details: errDetails,
		}, nil
	}

	order, err := a.orderService.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		return handleGetOrderError(err), nil
	}

	return converter.OrderToDto(order), nil
}

func validateGetOrderRequest(params *orderv1.GetOrderParams) (errDetails []orderv1.ErrorDetail, ok bool) {
	if params.OrderUUID == uuid.Nil {
		errDetails = append(errDetails, orderv1.ErrorDetail{
			Message: "Параметр order_uuid не может быть uuid null формата",
			Field:   "order_uuid",
		})
	}

	ok = len(errDetails) == 0

	return errDetails, ok
}

func handleGetOrderError(err error) orderv1.GetOrderRes {
	orderErr := apierr.FromError(err)
	if orderErr.Log {
		slog.Error(orderErr.Message, "err", err)
	}
	switch orderErr.Status {
	case http.StatusNotFound:
		return &orderv1.GetOrderNotFound{
			Code:    orderErr.Status,
			Message: orderErr.Message,
		}
	default:
		return &orderv1.GetOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: orderErr.Message,
		}
	}
}
