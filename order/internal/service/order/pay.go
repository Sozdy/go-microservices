package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/Sozdy/go-microservices/order/internal/errors"
	"github.com/Sozdy/go-microservices/order/internal/model"
	paymentModel "github.com/Sozdy/go-microservices/order/internal/model/ports/payment"
)

type PayOrderOrderIn struct {
	OrderUUID     uuid.UUID
	PaymentMethod model.PaymentMethod
}

type PayOrderOrderOut struct {
	TransactionUUID uuid.UUID
}

func (s *orderService) PayOrder(ctx context.Context, in *PayOrderOrderIn) (*PayOrderOrderOut, error) {
	order, err := s.orderRepository.Get(ctx, in.OrderUUID)
	if err != nil {
		return nil, fmt.Errorf("оплата заказа: %w", err)
	}

	if order.Status != model.OrderStatusPendingPayment {
		return nil, fmt.Errorf("оплата заказа: %w", errs.ErrPayOrderStatusConflict)
	}

	payOrderResponse, err := s.paymentClient.PayOrder(ctx, &paymentModel.PayOrderRequest{
		OrderUuid:     order.OrderUUID,
		PaymentMethod: in.PaymentMethod,
	})
	if err != nil {
		return nil, fmt.Errorf("оплата заказа: %w", err)
	}

	order.Status = model.OrderStatusPaid
	order.TransactionUUID = &payOrderResponse.TransactionUUID
	order.PaymentMethod = &in.PaymentMethod

	if err = s.orderRepository.Update(ctx, *order); err != nil {
		return nil, fmt.Errorf("оплата заказа: %w", err)
	}

	return &PayOrderOrderOut{
		TransactionUUID: payOrderResponse.TransactionUUID,
	}, nil
}
