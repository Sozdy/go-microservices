package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/Sozdy/go-microservices/inventory/internal/api/inventory/v1/converter"
	"github.com/Sozdy/go-microservices/inventory/internal/errs"
	inventoryv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(
	ctx context.Context,
	req *inventoryv1.GetPartRequest,
) (*inventoryv1.GetPartResponse, error) {
	if err := uuid.Validate(req.GetUuid()); err != nil {
		return nil, errs.ErrInvalidUUID
	}

	part, err := a.InventoryService.GetPart(ctx, req.GetUuid())
	if err != nil {
		return nil, err
	}

	return &inventoryv1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
