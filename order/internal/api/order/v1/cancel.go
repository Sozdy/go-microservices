package v1

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/Sozdy/go-microservices/order/internal/errs"
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
	switch errs.CodeOf(err) {
	case errs.CodeNotFound:
		return &orderv1.CancelOrderNotFound{Code: http.StatusNotFound, Message: errs.ClientMessage(err)}
	case errs.CodeInvalidArgument:
		return &orderv1.CancelOrderBadRequest{Code: http.StatusBadRequest, Message: errs.ClientMessage(err)}
	case errs.CodeConflict, errs.CodeFailedPrecondition:
		return &orderv1.CancelOrderConflict{Code: http.StatusConflict, Message: errs.ClientMessage(err)}
	default:
		slog.Error("внутренняя ошибка", "err", err)
		return &orderv1.CancelOrderInternalServerError{Code: http.StatusInternalServerError, Message: errs.ClientMessage(err)}
	}
}
