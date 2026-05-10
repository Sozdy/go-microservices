package order

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Sozdy/go-microservices/order/internal/model"
)

func TestGetOrder_Success(t *testing.T) {
	t.Parallel()

	orderUUID := uuid.New()
	storedOrder := &model.Order{
		OrderUUID:  orderUUID,
		HullUUID:   uuid.New(),
		EngineUUID: uuid.New(),
		Status:     model.OrderStatusPendingPayment,
	}

	// === Arrange ===
	fixture := newServiceFixture(t)

	// === Expect ===
	fixture.orderRepository.EXPECT().
		Get(fixture.ctx, orderUUID).
		Return(storedOrder, nil).
		Once()

	// === Act ===
	got, err := fixture.service.GetOrder(fixture.ctx, orderUUID)

	// === Assert ===
	require.NoError(t, err)
	require.Equal(t, storedOrder, got)
}

func TestGetOrder_RepositoryError(t *testing.T) {
	t.Parallel()

	orderUUID := uuid.New()
	repoErr := errors.New("сбой БД")

	// === Arrange ===
	fixture := newServiceFixture(t)

	// === Expect ===
	fixture.orderRepository.EXPECT().
		Get(fixture.ctx, orderUUID).
		Return(nil, repoErr).
		Once()

	// === Act ===
	got, err := fixture.service.GetOrder(fixture.ctx, orderUUID)

	// === Assert ===
	require.Error(t, err)
	require.Nil(t, got)
	require.ErrorIs(t, err, repoErr)
}
