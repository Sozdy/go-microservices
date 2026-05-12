package order

type orderService struct {
	inventoryClient InventoryClient
	paymentClient   PaymentClient
	orderRepository OrderRepository
}

func NewOrderService(
	inventoryClient InventoryClient,
	paymentClient PaymentClient,
	orderRepository OrderRepository,
) *orderService {
	return &orderService{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		orderRepository: orderRepository,
	}
}
