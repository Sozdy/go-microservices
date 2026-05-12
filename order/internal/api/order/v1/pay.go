package v1

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/Sozdy/go-microservices/order/internal/api/order/v1/converter"
	"github.com/Sozdy/go-microservices/order/internal/errs"
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
	switch errs.CodeOf(err) {
	case errs.CodeNotFound:
		return &orderv1.PayOrderNotFound{Code: http.StatusNotFound, Message: errs.ClientMessage(err)}
	case errs.CodeInvalidArgument:
		return &orderv1.PayOrderBadRequest{Code: http.StatusBadRequest, Message: errs.ClientMessage(err)}
	case errs.CodeConflict, errs.CodeFailedPrecondition:
		return &orderv1.PayOrderConflict{Code: http.StatusConflict, Message: errs.ClientMessage(err)}
	default:
		slog.Error("внутренняя ошибка", "err", err)
		return &orderv1.PayOrderInternalServerError{Code: http.StatusInternalServerError, Message: errs.ClientMessage(err)}
	}
}
