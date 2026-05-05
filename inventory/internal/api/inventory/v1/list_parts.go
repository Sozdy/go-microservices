package v1

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/status"

	"github.com/Sozdy/go-microservices/inventory/internal/api/inventory"
	"github.com/Sozdy/go-microservices/inventory/internal/api/inventory/v1/converter"
	inventoryv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(
	ctx context.Context,
	req *inventoryv1.ListPartsRequest,
) (*inventoryv1.ListPartsResponse, error) {
	parts, err := a.InventoryService.ListParts(
		ctx, req.GetUuids(),
		converter.PartTypeToModel(req.GetPartType()),
	)
	if err != nil {
		return nil, handleListPartsError(err)
	}

	return converter.PartsToProto(parts), nil
}

func handleListPartsError(err error) error {
	inventoryError := inventory.FromError(err)
	if inventoryError.Log {
		slog.Error(inventoryError.Message, "err", err)
	}

	return status.Errorf(inventoryError.Code, "%s", inventoryError.Message)
}
