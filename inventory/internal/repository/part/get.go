package part

import (
	"context"

	errs "github.com/Sozdy/go-microservices/inventory/internal/errors"
	"github.com/Sozdy/go-microservices/inventory/internal/model"
	"github.com/Sozdy/go-microservices/inventory/internal/repository/converter"
)

func (r *repo) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	part, ok := r.parts[uuid]
	if !ok {
		return nil, errs.ErrPartNotFound
	}

	return converter.PartToModel(&part), nil
}
