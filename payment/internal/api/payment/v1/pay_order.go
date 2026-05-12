package v1

import (
	"context"

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
		return nil, err
	}

	return converter.TransactionToProto(transaction), nil
}
