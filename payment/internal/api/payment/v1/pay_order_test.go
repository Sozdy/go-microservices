package v1

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	paymentErrs "github.com/Sozdy/go-microservices/payment/internal/errors"
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
		wantCode        codes.Code
	}{
		{
			name: "пустой order_uuid",
			request: &paymentv1.PayOrderRequest{
				OrderUuid:     "",
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			wantModelMethod: model.PaymentMethodCard,
			serviceErr:      paymentErrs.ErrOrderUUIDEmpty,
			wantCode:        codes.InvalidArgument,
		},
		{
			name: "невалидный order_uuid",
			request: &paymentv1.PayOrderRequest{
				OrderUuid:     "not-a-uuid",
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			wantModelMethod: model.PaymentMethodCard,
			serviceErr:      paymentErrs.ErrInvalidOrderUUID,
			wantCode:        codes.InvalidArgument,
		},
		{
			name: "UNSPECIFIED метод оплаты",
			request: &paymentv1.PayOrderRequest{
				OrderUuid:     gofakeit.UUID(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
			},
			wantModelMethod: model.PaymentMethodUnspecified,
			serviceErr:      paymentErrs.ErrPaymentMethodUnspecified,
			wantCode:        codes.InvalidArgument,
		},
		{
			name: "внутренняя ошибка сервиса",
			request: &paymentv1.PayOrderRequest{
				OrderUuid:     gofakeit.UUID(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			wantModelMethod: model.PaymentMethodCard,
			serviceErr:      errors.New("что-то пошло не так в БД"),
			wantCode:        codes.Internal,
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

			grpcStatus, isGRPCStatus := status.FromError(err)
			require.True(t, isGRPCStatus, "api должен возвращать grpc status error")
			require.Equal(t, testCase.wantCode, grpcStatus.Code())
		})
	}
}
