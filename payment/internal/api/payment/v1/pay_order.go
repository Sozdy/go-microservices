package v1

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/status"

	"github.com/Sozdy/go-microservices/payment/internal/api/payment"
	"github.com/Sozdy/go-microservices/payment/internal/api/payment/v1/converter"
	paymentv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(
	ctx context.Context,
	req *paymentv1.PayOrderRequest,
) (*paymentv1.PayOrderResponse, error) {
	transaction, err := a.PaymentService.PayOrder(
		ctx,
		req.GetOrderUuid(),
		converter.PaymentMethodToModel(req.GetPaymentMethod()),
	)
	if err != nil {
		return nil, handlePayOrderError(err)
	}

	return converter.TransactionToProto(transaction), nil
}

func handlePayOrderError(err error) error {
	paymentError := payment.FromError(err)
	if paymentError.Log {
		slog.Error(paymentError.Message, "err", err)
	}

	return status.Errorf(paymentError.Code, "%s", paymentError.Message)
}
