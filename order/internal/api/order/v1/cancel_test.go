package v1

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	orderErrs "github.com/Sozdy/go-microservices/order/internal/errors"
	"github.com/Sozdy/go-microservices/order/internal/service/order"
	orderv1 "github.com/Sozdy/go-microservices/shared/pkg/openapi/order/v1"
)

func TestCancelOrder_Success(t *testing.T) {
	t.Parallel()

	orderUUID := uuid.New()

	// === Arrange ===
	fixture := newApiFixture(t)

	// === Expect ===
	fixture.orderService.EXPECT().
		CancelOrder(fixture.ctx, orderUUID).
		Return(&order.CancelOrderOut{}, nil).
		Once()

	// === Act ===
	response, err := fixture.api.CancelOrder(fixture.ctx, orderv1.CancelOrderParams{OrderUUID: orderUUID})

	// === Assert ===
	require.NoError(t, err)
	require.IsType(t, &orderv1.CancelOrderResponse{}, response)
}

func TestCancelOrder_ValidationError(t *testing.T) {
	t.Parallel()

	// === Arrange ===
	fixture := newApiFixture(t)

	// === Expect ===

	// === Act ===
	response, err := fixture.api.CancelOrder(fixture.ctx, orderv1.CancelOrderParams{OrderUUID: uuid.Nil})

	// === Assert ===
	require.NoError(t, err)
	require.IsType(t, &orderv1.CancelOrderBadRequest{}, response)
}

func TestCancelOrder_ServiceError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		serviceErr error
		wantType   orderv1.CancelOrderRes
	}{
		{
			name:       "заказ не найден",
			serviceErr: orderErrs.ErrOrderNotFound,
			wantType:   &orderv1.CancelOrderNotFound{},
		},
		{
			name:       "конфликт статуса",
			serviceErr: orderErrs.ErrCancelOrderStatusConflict,
			wantType:   &orderv1.CancelOrderConflict{},
		},
		{
			name:       "внутренняя ошибка",
			serviceErr: errors.New("сбой БД"),
			wantType:   &orderv1.CancelOrderInternalServerError{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			orderUUID := uuid.New()

			// === Arrange ===
			fixture := newApiFixture(t)

			// === Expect ===
			fixture.orderService.EXPECT().
				CancelOrder(fixture.ctx, orderUUID).
				Return(nil, testCase.serviceErr).
				Once()

			// === Act ===
			response, err := fixture.api.CancelOrder(fixture.ctx, orderv1.CancelOrderParams{OrderUUID: orderUUID})

			// === Assert ===
			require.NoError(t, err)
			require.IsType(t, testCase.wantType, response)
		})
	}
}
