package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Sozdy/go-microservices/order/internal/model"
)

func (s *orderService) GetOrder(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return nil, fmt.Errorf("получение заказа: %w", err)
	}

	return order, nil
}
