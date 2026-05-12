package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/Sozdy/go-microservices/order/internal/model"
	"github.com/Sozdy/go-microservices/order/internal/service/order"
)

type OrderService interface {
	GetOrder(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error)
	CreateOrder(ctx context.Context, in *order.CreateOrderIn) (*order.CreateOrderOut, error)
	PayOrder(ctx context.Context, in *order.PayOrderOrderIn) (*order.PayOrderOrderOut, error)
	CancelOrder(ctx context.Context, orderUUID uuid.UUID) (*order.CancelOrderOut, error)
}
