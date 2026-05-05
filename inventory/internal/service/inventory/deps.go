package inventory

import (
	"context"

	"github.com/Sozdy/go-microservices/inventory/internal/model"
)

type PartRepository interface {
	GetPart(ctx context.Context, uuid string) (*model.Part, error)
	ListParts(ctx context.Context, partUUIDs []string, partType model.PartType) ([]*model.Part, error)
}
