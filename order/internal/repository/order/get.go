package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Sozdy/go-microservices/order/internal/errs"
	"github.com/Sozdy/go-microservices/order/internal/model"
	"github.com/Sozdy/go-microservices/order/internal/repository/converter"
	"github.com/Sozdy/go-microservices/order/internal/repository/record"
)

const getOrderQuery = `SELECT uuid, status, transaction_uuid, payment_method, created_at, updated_at
          FROM orders
          WHERE uuid = $1`

const getOrderItemsQuery = `
          SELECT uuid, order_uuid, part_uuid, part_type, price
          FROM order_items
          WHERE order_uuid = $1`

func (r *repo) Get(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	orderRow, err := db.Query(ctx, getOrderQuery, orderUUID)
	if err != nil {
		return nil, fmt.Errorf("запрос заказа: %w", err)
	}
	defer orderRow.Close()

	orderRec, err := pgx.CollectOneRow(orderRow, pgx.RowToStructByName[record.Order])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrOrderNotFound
		}
		return nil, fmt.Errorf("маппинг заказа: %w", err)
	}

	itemRows, err := db.Query(ctx, getOrderItemsQuery, orderUUID)
	if err != nil {
		return nil, fmt.Errorf("запрос позиций заказа: %w", err)
	}
	defer itemRows.Close()

	orderItems, err := pgx.CollectRows(itemRows, pgx.RowToStructByName[record.OrderItem])
	if err != nil {
		return nil, fmt.Errorf("маппинг позиций заказа: %w", err)
	}

	order := converter.OrderFromRecord(orderRec, orderItems)
	return &order, nil
}
