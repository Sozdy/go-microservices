package order

import (
	"context"
	"testing"

	"github.com/Sozdy/go-microservices/order/internal/service/order/mocks"
)

type serviceFixture struct {
	ctx             context.Context
	service         *orderService
	orderRepository *mocks.OrderRepository
	inventoryClient *mocks.InventoryClient
	paymentClient   *mocks.PaymentClient
}

func newServiceFixture(t *testing.T) *serviceFixture {
	t.Helper()

	orderRepository := mocks.NewOrderRepository(t)
	inventoryClient := mocks.NewInventoryClient(t)
	paymentClient := mocks.NewPaymentClient(t)

	return &serviceFixture{
		ctx:             context.Background(),
		service:         NewOrderService(inventoryClient, paymentClient, orderRepository),
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
