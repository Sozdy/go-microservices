package payment

import (
	"github.com/google/uuid"

	"github.com/Sozdy/go-microservices/order/internal/model"
)

type PayOrderRequest struct {
	OrderUuid     uuid.UUID
	PaymentMethod model.PaymentMethod
}

type PayOrderResponse struct {
	TransactionUUID uuid.UUID
}
