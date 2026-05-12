package v1

import (
	"context"

	"github.com/Sozdy/go-microservices/inventory/internal/model"
)

type InventoryService interface {
	GetPart(ctx context.Context, partUUID string) (*model.Part, error)
	ListParts(ctx context.Context, partUUIDs []string, partType model.PartType) ([]*model.Part, error)
}
