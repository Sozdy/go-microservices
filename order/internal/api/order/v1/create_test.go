package v1

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Sozdy/go-microservices/order/internal/errs"
	"github.com/Sozdy/go-microservices/order/internal/service/order"
	orderv1 "github.com/Sozdy/go-microservices/shared/pkg/openapi/order/v1"
)

func TestCreateOrder_Success(t *testing.T) {
	t.Parallel()

	hullUUID := uuid.New()
	engineUUID := uuid.New()
	createdOrderUUID := uuid.New()

	// === Arrange ===
	fixture := newApiFixture(t)

	// === Expect ===
	fixture.orderService.EXPECT().
		CreateOrder(fixture.ctx, mock.Anything).
		Return(&order.CreateOrderOut{
			OrderUUID:  createdOrderUUID,
			TotalPrice: 200,
		}, nil).
		Once()

	// === Act ===
	response, err := fixture.api.CreateOrder(fixture.ctx, &orderv1.CreateOrderRequest{
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
	})

	// === Assert ===
	require.NoError(t, err)
	createOrderResp, ok := response.(*orderv1.CreateOrderResponse)
	require.True(t, ok, "ожидался *orderv1.CreateOrderResponse")
	require.Equal(t, createdOrderUUID, createOrderResp.OrderUUID)
	require.Equal(t, int64(200), createOrderResp.TotalPrice)
}

func TestCreateOrder_ValidationError(t *testing.T) {
	t.Parallel()

	// === Arrange ===
	fixture := newApiFixture(t)

	// === Expect ===
	// (валидация падает на uuid.Nil - сервис не вызывается)

	// === Act ===
	response, err := fixture.api.CreateOrder(fixture.ctx, &orderv1.CreateOrderRequest{
		HullUUID:   uuid.Nil,
		EngineUUID: uuid.Nil,
	})

	// === Assert ===
	require.NoError(t, err)
	require.IsType(t, &orderv1.CreateOrderBadRequest{}, response)
}

func TestCreateOrder_ServiceError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		serviceErr error
		wantType   orderv1.CreateOrderRes
	}{
		{
			name:       "деталь не найдена",
			serviceErr: errs.ErrPartNotFound,
			wantType:   &orderv1.CreateOrderNotFound{},
		},
		{
			name:       "деталь недоступна на складе",
			serviceErr: errs.ErrPartUnavailable,
			wantType:   &orderv1.CreateOrderConflict{},
		},
		{
			name:       "inventory недоступен",
			serviceErr: errs.ErrInventoryUnavailable,
			wantType:   &orderv1.CreateOrderInternalServerError{},
		},
		{
			name:       "внутренняя ошибка сервиса",
			serviceErr: errors.New("сбой БД"),
			wantType:   &orderv1.CreateOrderInternalServerError{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// === Arrange ===
			fixture := newApiFixture(t)

			// === Expect ===
			fixture.orderService.EXPECT().
				CreateOrder(fixture.ctx, mock.Anything).
				Return(nil, testCase.serviceErr).
				Once()

			// === Act ===
			response, err := fixture.api.CreateOrder(fixture.ctx, &orderv1.CreateOrderRequest{
				HullUUID:   uuid.New(),
				EngineUUID: uuid.New(),
			})

			// === Assert ===
			require.NoError(t, err)
			require.IsType(t, testCase.wantType, response)
		})
	}
}
