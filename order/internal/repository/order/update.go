package order

import (
	"context"

	errs "github.com/Sozdy/go-microservices/order/internal/errors"
	"github.com/Sozdy/go-microservices/order/internal/model"
	"github.com/Sozdy/go-microservices/order/internal/repository/converter"
)

func (r *repo) Update(ctx context.Context, order model.Order) error {
	rec := converter.OrderToRecord(order)

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[rec.OrderUUID]; !exists {
		return errs.ErrOrderDoesNotExists
	}

	r.orders[rec.OrderUUID] = rec
	return nil
}
