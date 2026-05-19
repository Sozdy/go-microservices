package v1

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Sozdy/go-microservices/order/internal/errs"
	"github.com/Sozdy/go-microservices/order/internal/model"
	orderv1 "github.com/Sozdy/go-microservices/shared/pkg/openapi/order/v1"
)

func TestGetOrder_Success(t *testing.T) {
	t.Parallel()

	orderUUID := uuid.New()
	hullUUID := uuid.New()
	engineUUID := uuid.New()

	storedOrder := &model.Order{
		UUID: orderUUID,
		OrderItems: []model.OrderItem{
			{PartUUID: hullUUID, PartType: model.PartTypeHull, Price: 30},
			{PartUUID: engineUUID, PartType: model.PartTypeEngine, Price: 70},
		},
		TotalPrice: 100,
		Status:     model.OrderStatusPendingPayment,
		CreatedAt:  time.Now().UTC(),
	}

	// === Arrange ===
	fixture := newApiFixture(t)

	// === Expect ===
	fixture.orderService.EXPECT().
		GetOrder(fixture.ctx, orderUUID).
		Return(storedOrder, nil).
		Once()

	// === Act ===
	response, err := fixture.api.GetOrder(fixture.ctx, orderv1.GetOrderParams{OrderUUID: orderUUID})

	// === Assert ===
	require.NoError(t, err)
	dto, ok := response.(*orderv1.OrderDto)
	require.True(t, ok, "ожидался *orderv1.OrderDto")
	require.Equal(t, orderUUID, dto.OrderUUID)
	require.Equal(t, hullUUID, dto.HullUUID)
}

func TestGetOrder_ValidationError(t *testing.T) {
	t.Parallel()

	// === Arrange ===
	fixture := newApiFixture(t)

	// === Expect ===
	// (валидация в api отрабатывает до обращения к сервису)

	// === Act ===
	response, err := fixture.api.GetOrder(fixture.ctx, orderv1.GetOrderParams{OrderUUID: uuid.Nil})

	// === Assert ===
	require.NoError(t, err)
	require.IsType(t, &orderv1.GetOrderBadRequest{}, response)
}

func TestGetOrder_ServiceError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		serviceErr error
		wantType   orderv1.GetOrderRes
	}{
		{
			name:       "заказ не найден",
			serviceErr: errs.ErrOrderNotFound,
			wantType:   &orderv1.GetOrderNotFound{},
		},
		{
			name:       "внутренняя ошибка сервиса",
			serviceErr: errors.New("что-то пошло не так в БД"),
			wantType:   &orderv1.GetOrderInternalServerError{},
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
				GetOrder(fixture.ctx, orderUUID).
				Return(nil, testCase.serviceErr).
				Once()

			// === Act ===
			response, err := fixture.api.GetOrder(fixture.ctx, orderv1.GetOrderParams{OrderUUID: orderUUID})

			// === Assert ===
			require.NoError(t, err)
			require.IsType(t, testCase.wantType, response)
		})
	}
}
