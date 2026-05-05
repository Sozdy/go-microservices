package payment

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/Sozdy/go-microservices/payment/internal/errors"
	"github.com/Sozdy/go-microservices/payment/internal/model"
)

func (s *paymentService) PayOrder(ctx context.Context, orderUUID string, method model.PaymentMethod) (*model.Transaction, error) {
	if orderUUID == "" {
		return nil, errs.ErrOrderUUIDEmpty
	}

	if err := uuid.Validate(orderUUID); err != nil {
		return nil, fmt.Errorf("валидация order_uuid: %w", errs.ErrInvalidOrderUUID)
	}

	if method == model.PaymentMethodUnspecified {
		return nil, errs.ErrPaymentMethodUnspecified
	}

	transactionUUID := uuid.New()

	slog.Info("оплата прошла успешно",
		"order_uuid", orderUUID,
		"transaction_uuid", transactionUUID,
	)

	return &model.Transaction{
		UUID: transactionUUID.String(),
	}, nil
}
