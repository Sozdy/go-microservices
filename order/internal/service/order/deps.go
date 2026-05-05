package order

import (
	"context"

	"github.com/google/uuid"

	internalModel "github.com/Sozdy/go-microservices/order/internal/model"
	inventoryModel "github.com/Sozdy/go-microservices/order/internal/model/ports/inventory"
	paymentModel "github.com/Sozdy/go-microservices/order/internal/model/ports/payment"
)

type InventoryClient interface {
	GetPart(ctx context.Context, partUUID uuid.UUID) (*inventoryModel.Part, error)
	ListParts(ctx context.Context, req *inventoryModel.ListPartsRequest) (*inventoryModel.ListPartsResponse, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, req *paymentModel.PayOrderRequest) (*paymentModel.PayOrderResponse, error)
}

type OrderRepository interface {
	Create(ctx context.Context, order internalModel.Order) error
	Get(ctx context.Context, orderUUID uuid.UUID) (*internalModel.Order, error)
	Update(ctx context.Context, order internalModel.Order) error
}
