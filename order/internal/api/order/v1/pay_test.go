package v1

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	orderErrs "github.com/Sozdy/go-microservices/order/internal/errors"
	"github.com/Sozdy/go-microservices/order/internal/service/order"
	orderv1 "github.com/Sozdy/go-microservices/shared/pkg/openapi/order/v1"
)

func TestPayOrder_Success(t *testing.T) {
	t.Parallel()

	orderUUID := uuid.New()
	transactionUUID := uuid.New()

	// === Arrange ===
	fixture := newApiFixture(t)

	// === Expect ===
	fixture.orderService.EXPECT().
		PayOrder(fixture.ctx, mock.Anything).
		Return(&order.PayOrderOrderOut{TransactionUUID: transactionUUID}, nil).
		Once()

	// === Act ===
	response, err := fixture.api.PayOrder(fixture.ctx, &orderv1.PayOrderRequest{
		PaymentMethod: orderv1.PaymentMethodCARD,
	}, orderv1.PayOrderParams{OrderUUID: orderUUID})

	// === Assert ===
	require.NoError(t, err)
	payResp, ok := response.(*orderv1.PayOrderResponse)
	require.True(t, ok, "ожидался *orderv1.PayOrderResponse")
	require.Equal(t, transactionUUID, payResp.TransactionUUID)
}

func TestPayOrder_ValidationError(t *testing.T) {
	t.Parallel()

	// === Arrange ===
	fixture := newApiFixture(t)

	// === Expect ===

	// === Act ===
	response, err := fixture.api.PayOrder(fixture.ctx, &orderv1.PayOrderRequest{
		PaymentMethod: orderv1.PaymentMethodCARD,
	}, orderv1.PayOrderParams{OrderUUID: uuid.Nil})

	// === Assert ===
	require.NoError(t, err)
	require.IsType(t, &orderv1.PayOrderBadRequest{}, response)
}

func TestPayOrder_ServiceError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		serviceErr error
		wantType   orderv1.PayOrderRes
	}{
		{
			name:       "заказ не найден",
			serviceErr: orderErrs.ErrOrderNotFound,
			wantType:   &orderv1.PayOrderNotFound{},
		},
		{
			name:       "конфликт статуса",
			serviceErr: orderErrs.ErrPayOrderStatusConflict,
			wantType:   &orderv1.PayOrderConflict{},
		},
		{
			name:       "payment недоступен",
			serviceErr: orderErrs.ErrPaymentUnavailable,
			wantType:   &orderv1.PayOrderInternalServerError{},
		},
		{
			name:       "внутренняя ошибка",
			serviceErr: errors.New("сбой БД"),
			wantType:   &orderv1.PayOrderInternalServerError{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// === Arrange ===
			fixture := newApiFixture(t)

			// === Expect ===
			fixture.orderService.EXPECT().
				PayOrder(fixture.ctx, mock.Anything).
				Return(nil, testCase.serviceErr).
				Once()

			// === Act ===
			response, err := fixture.api.PayOrder(fixture.ctx, &orderv1.PayOrderRequest{
				PaymentMethod: orderv1.PaymentMethodCARD,
			}, orderv1.PayOrderParams{OrderUUID: uuid.New()})

			// === Assert ===
			require.NoError(t, err)
			require.IsType(t, testCase.wantType, response)
		})
	}
}
