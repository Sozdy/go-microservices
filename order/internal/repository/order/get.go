package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/Sozdy/go-microservices/order/internal/errs"
	"github.com/Sozdy/go-microservices/order/internal/model"
	"github.com/Sozdy/go-microservices/order/internal/repository/converter"
)

func (r *repo) Get(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rec, exists := r.orders[orderUUID]
	if !exists {
		return nil, errs.ErrOrderNotFound
	}

	order := converter.OrderFromRecord(rec)
	return &order, nil
}
