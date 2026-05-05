package converter

import (
	"github.com/google/uuid"

	"github.com/Sozdy/go-microservices/order/internal/model"
	"github.com/Sozdy/go-microservices/shared/pkg/openapi/order/v1"
)

// OrderToDto преобразует model.Order в openapi-DTO для ответа клиенту.
func OrderToDto(o *model.Order) *orderv1.OrderDto {
	return &orderv1.OrderDto{
		OrderUUID:       o.OrderUUID,
		HullUUID:        o.HullUUID,
		EngineUUID:      o.EngineUUID,
		ShieldUUID:      uuidPtrToOptNil(o.ShieldUUID),
		WeaponUUID:      uuidPtrToOptNil(o.WeaponUUID),
		TotalPrice:      o.TotalPrice,
		TransactionUUID: uuidPtrToOptNil(o.TransactionUUID),
		PaymentMethod:   paymentMethodPtrToOptNil(o.PaymentMethod),
		Status:          orderStatusToDto(o.Status),
		CreatedAt:       o.CreatedAt,
	}
}

// uuidPtrToOptNil превращает *uuid.UUID в OptNilUUID.
func uuidPtrToOptNil(p *uuid.UUID) orderv1.OptNilUUID {
	if p == nil {
		var v orderv1.OptNilUUID
		v.SetToNull()
		return v
	}
	return orderv1.NewOptNilUUID(*p)
}

// paymentMethodPtrToOptNil превращает *string в OptNilPaymentMethod.
func paymentMethodPtrToOptNil(p *model.PaymentMethod) orderv1.OptNilPaymentMethod {
	if p == nil {
		var v orderv1.OptNilPaymentMethod
		v.SetToNull()
		return v
	}
	return orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethod(*p))
}

// orderStatusToDto переводит доменный enum в openapi enum.
func orderStatusToDto(s model.OrderStatus) orderv1.OrderStatus {
	switch s {
	case model.OrderStatusPendingPayment:
		return orderv1.OrderStatusPENDINGPAYMENT
	case model.OrderStatusPaid:
		return orderv1.OrderStatusPAID
	case model.OrderStatusCancelled:
		return orderv1.OrderStatusCANCELLED
	default:
		return ""
	}
}
