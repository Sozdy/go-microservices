package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Sozdy/go-microservices/order/internal/model"
	"github.com/Sozdy/go-microservices/order/internal/repository/converter"
)

const queryCreateOrder = `
        INSERT INTO orders (uuid, status, created_at)
        VALUES ($1, $2, $3);`

const queryCreateOrderItems = `
      INSERT INTO order_items (uuid, order_uuid, part_uuid, part_type, price)
      SELECT * FROM unnest($1::uuid[], $2::uuid[], $3::uuid[], $4::text[], $5::bigint[]);`

func (r *repo) Create(ctx context.Context, order model.Order) error {
	return r.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := r.insertOrder(txCtx, order); err != nil {
			return err
		}
		return r.insertOrderItems(txCtx, order)
	})
}

func (r *repo) insertOrder(ctx context.Context, order model.Order) error {
	record := converter.OrderToRecord(order)

	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	_, err := db.Exec(ctx, queryCreateOrder,
		record.UUID,
		record.Status,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("создать заказ: %w", err)
	}

	return nil
}

func (r *repo) insertOrderItems(ctx context.Context, order model.Order) error {
	records := converter.OrderItemsToRecord(order)

	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	UUIDs := make([]uuid.UUID, len(records))
	orderUUIDs := make([]uuid.UUID, len(records))
	partUUIDs := make([]uuid.UUID, len(records))
	partTypes := make([]string, len(records))
	prices := make([]int64, len(records))

	for i, record := range records {
		UUIDs[i] = record.UUID
		orderUUIDs[i] = record.OrderUUID
		partUUIDs[i] = record.PartUUID
		partTypes[i] = string(record.PartType)
		prices[i] = record.Price
	}

	_, err := db.Exec(ctx, queryCreateOrderItems, UUIDs, orderUUIDs, partUUIDs, partTypes, prices)
	if err != nil {
		return fmt.Errorf("создать элементы заказа: %w", err)
	}

	return nil
}
