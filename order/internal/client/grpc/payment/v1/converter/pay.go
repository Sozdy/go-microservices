package converter

import (
	"github.com/google/uuid"

	"github.com/Sozdy/go-microservices/order/internal/model"
	"github.com/Sozdy/go-microservices/order/internal/model/ports/payment"
	paymentv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/payment/v1"
)

func PaymentMethodToProto(method model.PaymentMethod) paymentv1.PaymentMethod {
	switch method {
	case model.PaymentMethodCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCreditCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

func PayOrderRequestToProto(req *payment.PayOrderRequest) *paymentv1.PayOrderRequest {
	return &paymentv1.PayOrderRequest{
		OrderUuid:     req.OrderUuid.String(),
		PaymentMethod: PaymentMethodToProto(req.PaymentMethod),
	}
}

func PayOrderResponseFromProto(res *paymentv1.PayOrderResponse) *payment.PayOrderResponse {
	return &payment.PayOrderResponse{
		TransactionUUID: uuid.MustParse(res.TransactionUuid),
	}
}
