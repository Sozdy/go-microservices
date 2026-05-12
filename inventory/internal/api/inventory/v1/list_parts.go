package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/Sozdy/go-microservices/inventory/internal/api/inventory/v1/converter"
	"github.com/Sozdy/go-microservices/inventory/internal/errs"
	inventoryv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(
	ctx context.Context,
	req *inventoryv1.ListPartsRequest,
) (*inventoryv1.ListPartsResponse, error) {
	for _, partUUID := range req.GetUuids() {
		if err := uuid.Validate(partUUID); err != nil {
			return nil, errs.ErrInvalidUUID
		}
	}

	parts, err := a.InventoryService.ListParts(
		ctx, req.GetUuids(),
		converter.PartTypeToModel(req.GetPartType()),
	)
	if err != nil {
		return nil, err
	}

	return converter.PartsToProto(parts), nil
}
