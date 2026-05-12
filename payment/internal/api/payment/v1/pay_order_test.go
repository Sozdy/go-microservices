package v1

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	"github.com/Sozdy/go-microservices/payment/internal/errs"
	"github.com/Sozdy/go-microservices/payment/internal/model"
	paymentv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/payment/v1"
)

func TestPayOrder_Success(t *testing.T) {
	t.Parallel()

	transactionUUID := gofakeit.UUID()

	testCases := []struct {
		name            string
		request         *paymentv1.PayOrderRequest
		wantModelMethod model.PaymentMethod
	}{
		{
			name: "оплата картой",
			request: &paymentv1.PayOrderRequest{
				OrderUuid:     gofakeit.UUID(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			wantModelMethod: model.PaymentMethodCard,
		},
		{
			name: "оплата СБП",
			request: &paymentv1.PayOrderRequest{
				OrderUuid:     gofakeit.UUID(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_SBP,
			},
			wantModelMethod: model.PaymentMethodSBP,
		},
		{
			name: "оплата кредитной картой",
			request: &paymentv1.PayOrderRequest{
				OrderUuid:     gofakeit.UUID(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
			},
			wantModelMethod: model.PaymentMethodCreditCard,
		},
		{
			name: "оплата деньгами инвестора",
			request: &paymentv1.PayOrderRequest{
				OrderUuid:     gofakeit.UUID(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
			},
			wantModelMethod: model.PaymentMethodInvestorMoney,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// === Arrange ===
			fixture := newAPIFixture(t)

			// === Expect ===
			fixture.paymentService.EXPECT().
				PayOrder(fixture.ctx, testCase.request.OrderUuid, testCase.wantModelMethod).
				Return(&model.Transaction{UUID: transactionUUID}, nil).
				Once()

			// === Act ===
			response, err := fixture.api.PayOrder(fixture.ctx, testCase.request)

			// === Assert ===
			require.NoError(t, err)
			require.Equal(t, &paymentv1.PayOrderResponse{TransactionUuid: transactionUUID}, response)
		})
	}
}

func TestPayOrder_ServiceError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		request         *paymentv1.PayOrderRequest
		wantModelMethod model.PaymentMethod
		serviceErr      error
		wantCode        errs.Code
	}{
		{
			name: "пустой order_uuid",
			request: &paymentv1.PayOrderRequest{
				OrderUuid:     "",
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			wantModelMethod: model.PaymentMethodCard,
			serviceErr:      errs.ErrOrderUUIDEmpty,
			wantCode:        errs.CodeInvalidArgument,
		},
		{
			name: "невалидный order_uuid",
			request: &paymentv1.PayOrderRequest{
				OrderUuid:     "not-a-uuid",
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			wantModelMethod: model.PaymentMethodCard,
			serviceErr:      errs.ErrInvalidOrderUUID,
			wantCode:        errs.CodeInvalidArgument,
		},
		{
			name: "UNSPECIFIED метод оплаты",
			request: &paymentv1.PayOrderRequest{
				OrderUuid:     gofakeit.UUID(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
			},
			wantModelMethod: model.PaymentMethodUnspecified,
			serviceErr:      errs.ErrPaymentMethodUnspecified,
			wantCode:        errs.CodeInvalidArgument,
		},
		{
			name: "внутренняя ошибка сервиса",
			request: &paymentv1.PayOrderRequest{
				OrderUuid:     gofakeit.UUID(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			wantModelMethod: model.PaymentMethodCard,
			serviceErr:      errors.New("что-то пошло не так в БД"),
			wantCode:        errs.CodeInternal,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// === Arrange ===
			fixture := newAPIFixture(t)

			// === Expect ===
			fixture.paymentService.EXPECT().
				PayOrder(fixture.ctx, testCase.request.OrderUuid, testCase.wantModelMethod).
				Return(nil, testCase.serviceErr).
				Once()

			// === Act ===
			response, err := fixture.api.PayOrder(fixture.ctx, testCase.request)

			// === Assert ===
			require.Error(t, err)
			require.Nil(t, response)
			require.Equal(t, testCase.wantCode, errs.CodeOf(err))
		})
	}
}
