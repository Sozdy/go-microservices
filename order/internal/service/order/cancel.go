package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/Sozdy/go-microservices/order/internal/errors"
	"github.com/Sozdy/go-microservices/order/internal/model"
)

type CancelOrderOut struct{}

func (s *orderService) CancelOrder(ctx context.Context, orderUUID uuid.UUID) (*CancelOrderOut, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return nil, fmt.Errorf("отмена заказа: %w", err)
	}

	if order.Status != model.OrderStatusPendingPayment {
		return nil, errs.ErrCancelOrderStatusConflict
	}

	// 3. Обновить статус на CANCELLED
	order.Status = model.OrderStatusCancelled
	if err := s.orderRepository.Update(ctx, *order); err != nil {
		return nil, fmt.Errorf("обновление статуса заказа: %w", err)
	}

	return &CancelOrderOut{}, nil
}
