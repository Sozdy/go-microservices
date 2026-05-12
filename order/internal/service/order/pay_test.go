package order

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Sozdy/go-microservices/order/internal/errs"
	"github.com/Sozdy/go-microservices/order/internal/model"
	paymentModel "github.com/Sozdy/go-microservices/order/internal/model/ports/payment"
)

func TestPayOrder_Success(t *testing.T) {
	t.Parallel()

	orderUUID := uuid.New()
	transactionUUID := uuid.New()
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
	fixture.paymentClient.EXPECT().
		PayOrder(fixture.ctx, mock.Anything).
		Return(&paymentModel.PayOrderResponse{TransactionUUID: transactionUUID}, nil).
		Once()
	fixture.orderRepository.EXPECT().
		Update(fixture.ctx, mock.MatchedBy(func(order model.Order) bool {
			return order.Status == model.OrderStatusPaid &&
				order.TransactionUUID != nil &&
				*order.TransactionUUID == transactionUUID
		})).
		Return(nil).
		Once()

	// === Act ===
	out, err := fixture.service.PayOrder(fixture.ctx, &PayOrderOrderIn{
		OrderUUID:     orderUUID,
		PaymentMethod: model.PaymentMethodCard,
	})

	// === Assert ===
	require.NoError(t, err)
	require.Equal(t, transactionUUID, out.TransactionUUID)
}

func TestPayOrder_StatusConflict(t *testing.T) {
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
	// (paymentClient и Update не должны вызываться)

	// === Act ===
	out, err := fixture.service.PayOrder(fixture.ctx, &PayOrderOrderIn{
		OrderUUID:     orderUUID,
		PaymentMethod: model.PaymentMethodCard,
	})

	// === Assert ===
	require.Error(t, err)
	require.Nil(t, out)
	require.ErrorIs(t, err, errs.ErrPayOrderStatusConflict)
}
