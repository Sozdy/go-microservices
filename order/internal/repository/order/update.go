package order

import (
	"context"
	"fmt"
	"time"

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

func (r *repo) Update(ctx context.Context, order model.Order) error {
	return r.txManager.Do(ctx, func(txCtx context.Context) error {
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
	})
}
