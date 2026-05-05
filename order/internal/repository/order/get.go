package order

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/Sozdy/go-microservices/order/internal/errors"
	"github.com/Sozdy/go-microservices/order/internal/model"
	"github.com/Sozdy/go-microservices/order/internal/repository/converter"
)

func (r *repo) Get(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rec, exists := r.orders[orderUUID]
	if !exists {
		return nil, errs.ErrOrderNotFound
	}

	order := converter.OrderFromRecord(rec)
	return &order, nil
}
