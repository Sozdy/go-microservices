package order

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Sozdy/go-microservices/order/internal/errs"
	"github.com/Sozdy/go-microservices/order/internal/model"
)

func TestCancelOrder_Success(t *testing.T) {
	t.Parallel()

	orderUUID := uuid.New()
	storedOrder := &model.Order{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusPendingPayment,
	}

	// === Arrange ===
	fixture := newServiceFixture(t)

	// === Expect ===
	fixture.orderRepository.EXPECT().
		Get(fixture.ctx, orderUUID).
		Return(storedOrder, nil).
		Once()
	fixture.orderRepository.EXPECT().
		Update(fixture.ctx, mock.MatchedBy(func(order model.Order) bool {
			return order.OrderUUID == orderUUID && order.Status == model.OrderStatusCancelled
		})).
		Return(nil).
		Once()

	// === Act ===
	out, err := fixture.service.CancelOrder(fixture.ctx, orderUUID)

	// === Assert ===
	require.NoError(t, err)
	require.NotNil(t, out)
}

func TestCancelOrder_StatusConflict(t *testing.T) {
	t.Parallel()

	orderUUID := uuid.New()
	storedOrder := &model.Order{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusPaid,
	}

	// === Arrange ===
	fixture := newServiceFixture(t)

	// === Expect ===
	fixture.orderRepository.EXPECT().
		Get(fixture.ctx, orderUUID).
		Return(storedOrder, nil).
		Once()
	// (Update не должен вызываться при конфликте статуса)

	// === Act ===
	out, err := fixture.service.CancelOrder(fixture.ctx, orderUUID)

	// === Assert ===
	require.Error(t, err)
	require.Nil(t, out)
	require.ErrorIs(t, err, errs.ErrCancelOrderStatusConflict)
}

func TestCancelOrder_UpdateError(t *testing.T) {
	t.Parallel()

	orderUUID := uuid.New()
	storedOrder := &model.Order{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusPendingPayment,
	}
	updateErr := errors.New("сбой БД при обновлении")

	// === Arrange ===
	fixture := newServiceFixture(t)

	// === Expect ===
	fixture.orderRepository.EXPECT().
		Get(fixture.ctx, orderUUID).
		Return(storedOrder, nil).
		Once()
	fixture.orderRepository.EXPECT().
		Update(fixture.ctx, mock.Anything).
		Return(updateErr).
		Once()

	// === Act ===
	out, err := fixture.service.CancelOrder(fixture.ctx, orderUUID)

	// === Assert ===
	require.Error(t, err)
	require.Nil(t, out)
	require.ErrorIs(t, err, updateErr)
}
