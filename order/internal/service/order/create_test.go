package order

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Sozdy/go-microservices/order/internal/errs"
	"github.com/Sozdy/go-microservices/order/internal/model"
	inventoryModel "github.com/Sozdy/go-microservices/order/internal/model/ports/inventory"
)

func TestCreateOrder_Success(t *testing.T) {
	t.Parallel()

	hullUUID := uuid.New()
	engineUUID := uuid.New()
	inputOrder := &CreateOrderIn{
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
	}

	// === Arrange ===
	fixture := newServiceFixture(t)

	// === Expect ===
	fixture.inventoryClient.EXPECT().
		ListParts(fixture.ctx, mock.Anything).
		Return(&inventoryModel.ListPartsResponse{
			Parts: []inventoryModel.Part{
				{UUID: hullUUID, Price: 50, StockQuantity: 1},
				{UUID: engineUUID, Price: 70, StockQuantity: 2},
			},
		}, nil).
		Once()
	fixture.orderRepository.EXPECT().
		Create(fixture.ctx, mock.MatchedBy(func(order model.Order) bool {
			return order.HullUUID == hullUUID &&
				order.EngineUUID == engineUUID &&
				order.TotalPrice == 120 &&
				order.Status == model.OrderStatusPendingPayment
		})).
		Return(nil).
		Once()

	// === Act ===
	out, err := fixture.service.CreateOrder(fixture.ctx, inputOrder)

	// === Assert ===
	require.NoError(t, err)
	require.NotNil(t, out)
	require.Equal(t, int64(120), out.TotalPrice)
}

func TestCreateOrder_PartUnavailable(t *testing.T) {
	t.Parallel()

	hullUUID := uuid.New()
	engineUUID := uuid.New()
	inputOrder := &CreateOrderIn{
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
	}

	// === Arrange ===
	fixture := newServiceFixture(t)

	// === Expect ===
	fixture.inventoryClient.EXPECT().
		ListParts(fixture.ctx, mock.Anything).
		Return(&inventoryModel.ListPartsResponse{
			Parts: []inventoryModel.Part{
				{UUID: hullUUID, Price: 50, StockQuantity: 0}, // нет на складе
			},
		}, nil).
		Once()
	// (Create не должен вызываться при недоступной детали)

	// === Act ===
	out, err := fixture.service.CreateOrder(fixture.ctx, inputOrder)

	// === Assert ===
	require.Error(t, err)
	require.Nil(t, out)
	require.ErrorIs(t, err, errs.ErrPartUnavailable)
}

func TestCreateOrder_InventoryClientError(t *testing.T) {
	t.Parallel()

	inputOrder := &CreateOrderIn{
		HullUUID:   uuid.New(),
		EngineUUID: uuid.New(),
	}

	// === Arrange ===
	fixture := newServiceFixture(t)

	// === Expect ===
	fixture.inventoryClient.EXPECT().
		ListParts(fixture.ctx, mock.Anything).
		Return(nil, errors.New("inventory недоступен")).
		Once()

	// === Act ===
	out, err := fixture.service.CreateOrder(fixture.ctx, inputOrder)

	// === Assert ===
	require.Error(t, err)
	require.Nil(t, out)
}
