package converter

import (
	"github.com/Sozdy/go-microservices/order/internal/model"
	"github.com/Sozdy/go-microservices/order/internal/service/order"
	orderv1 "github.com/Sozdy/go-microservices/shared/pkg/openapi/order/v1"
)

func ToPaymentMethod(method orderv1.PaymentMethod) model.PaymentMethod {
	switch method {
	case orderv1.PaymentMethodCARD:
		return model.PaymentMethodCard
	case orderv1.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case orderv1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCreditCard
	case orderv1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}

func PayOrderInFromRequest(req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) *order.PayOrderOrderIn {
	return &order.PayOrderOrderIn{
		OrderUUID:     params.OrderUUID,
		PaymentMethod: ToPaymentMethod(req.PaymentMethod),
	}
}

func PayOrderOutToResponse(out *order.PayOrderOrderOut) orderv1.PayOrderRes {
	return &orderv1.PayOrderResponse{
		TransactionUUID: out.TransactionUUID,
	}
}
