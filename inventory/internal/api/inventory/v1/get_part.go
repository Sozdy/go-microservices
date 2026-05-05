package v1

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Sozdy/go-microservices/inventory/internal/api/inventory"
	"github.com/Sozdy/go-microservices/inventory/internal/api/inventory/v1/converter"
	inventoryv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(
	ctx context.Context,
	req *inventoryv1.GetPartRequest,
) (*inventoryv1.GetPartResponse, error) {
	if req.Uuid == "" {
		return nil, status.Errorf(codes.InvalidArgument, "uuid не может быть пустым")
	}
	if _, err := uuid.Parse(req.Uuid); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "uuid не является валидным: %v", err)
	}

	part, err := a.InventoryService.GetPart(ctx, req.Uuid)
	if err != nil {
		return nil, handleGetPartError(err)
	}

	return &inventoryv1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}

func handleGetPartError(err error) error {
	inventoryError := inventory.FromError(err)
	if inventoryError.Log {
		slog.Error(inventoryError.Message, "err", err)
	}

	return status.Errorf(inventoryError.Code, "%s", inventoryError.Message)
}
