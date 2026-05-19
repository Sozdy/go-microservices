package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Sozdy/go-microservices/order/internal/errs"
	"github.com/Sozdy/go-microservices/order/internal/model"
	"github.com/Sozdy/go-microservices/order/internal/repository/converter"
)

const updateOrderQuery = `
          UPDATE orders                                                                                                                                                                                                                                                                      
          SET status           = $2,
              transaction_uuid = $3,
              payment_method   = $4,
              updated_at       = $5
          WHERE uuid = $1`

const deleteOrderItemsQuery = `DELETE FROM order_items WHERE order_uuid = $1`

func (r *repo) Update(ctx context.Context, order model.Order) error {
	return r.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := r.updateOrder(txCtx, order); err != nil {
			return err
		}
		if err := r.deleteOrderItems(txCtx, order.UUID); err != nil {
			return err
		}
		return r.insertOrderItems(txCtx, order)
	})
}

func (r *repo) updateOrder(ctx context.Context, order model.Order) error {
	rec := converter.OrderToRecord(order)

	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	tag, err := db.Exec(ctx, updateOrderQuery,
		rec.UUID,
		rec.Status,
		rec.TransactionUUID,
		rec.PaymentMethod,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("обновить заказ: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrOrderDoesNotExists
	}
	return nil
}

func (r *repo) deleteOrderItems(ctx context.Context, orderUUID uuid.UUID) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	_, err := db.Exec(ctx, deleteOrderItemsQuery, orderUUID)
	if err != nil {
		return fmt.Errorf("удалить позиции заказа: %w", err)
	}
	return nil
}
