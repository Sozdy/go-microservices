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
				{UUID: hullUUID, PartType: inventoryModel.PART_TYPE_HULL, Price: 50, StockQuantity: 1},
				{UUID: engineUUID, PartType: inventoryModel.PART_TYPE_ENGINE, Price: 70, StockQuantity: 2},
			},
		}, nil).
		Once()
	fixture.orderRepository.EXPECT().
		Create(fixture.ctx, mock.MatchedBy(func(order model.Order) bool {
			return hasItem(order.OrderItems, hullUUID, model.PartTypeHull) &&
				hasItem(order.OrderItems, engineUUID, model.PartTypeEngine) &&
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

func hasItem(items []model.OrderItem, partUUID uuid.UUID, partType model.PartType) bool {
	for _, item := range items {
		if item.PartUUID == partUUID && item.PartType == partType {
			return true
		}
	}
	return false
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
