package v1

import (
	"context"

	"github.com/Sozdy/go-microservices/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, orderUUID string, method model.PaymentMethod) (*model.Transaction, error)
}
