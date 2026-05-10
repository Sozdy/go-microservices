package payment

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"

	errs "github.com/Sozdy/go-microservices/payment/internal/errors"
	"github.com/Sozdy/go-microservices/payment/internal/model"
)

func TestPayOrderSuccess(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		uuid          string
		paymentMethod model.PaymentMethod
		expected      int
	}{
		{
			name:          "оплата картой",
			uuid:          gofakeit.UUID(),
			paymentMethod: model.PaymentMethodCard,
		},
		{
			name:          "оплата СБП",
			uuid:          gofakeit.UUID(),
			paymentMethod: model.PaymentMethodSBP,
		},
		{
			name:          "оплата кредитной картой",
			uuid:          gofakeit.UUID(),
			paymentMethod: model.PaymentMethodCreditCard,
		},
		{
			name:          "оплата деньгами инвестора",
			uuid:          gofakeit.UUID(),
			paymentMethod: model.PaymentMethodInvestorMoney,
		},
	}

	paymentService := NewPaymentService()
	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			transaction, err := paymentService.PayOrder(ctx, test.uuid, test.paymentMethod)
			assert.NoError(t, err)
			assert.NotNil(t, transaction)
		})
	}
}

func TestPayOrderError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		uuid          string
		paymentMethod model.PaymentMethod
		expectedError error
	}{
		{
			name:          "неизвестный метод оплаты",
			uuid:          gofakeit.UUID(),
			paymentMethod: model.PaymentMethodUnspecified,
			expectedError: errs.ErrPaymentMethodUnspecified,
		},
		{
			name:          "не валидный uuid",
			uuid:          "not valid uuid",
			paymentMethod: model.PaymentMethodCard,
			expectedError: errs.ErrInvalidOrderUUID,
		},
		{
			name:          "пустой uuid",
			uuid:          "",
			paymentMethod: model.PaymentMethodCard,
			expectedError: errs.ErrOrderUUIDEmpty,
		},
	}

	paymentService := NewPaymentService()
	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := paymentService.PayOrder(ctx, test.uuid, test.paymentMethod)
			if err == nil {
				t.Errorf("Test %s failed: expected error but got nil", test.name)
			} else if !errors.Is(err, test.expectedError) {
				t.Errorf("Test %s failed: expected error %v but got %v", test.name, test.expectedError, err)
			}
		})
	}
}
