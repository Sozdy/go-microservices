package inventory

import (
	"context"
	"fmt"

	"github.com/Sozdy/go-microservices/inventory/internal/model"
)

func (s *partService) ListParts(ctx context.Context, partUUIDs []string, partType model.PartType) ([]*model.Part, error) {
	parts, err := s.partRepository.ListParts(ctx, partUUIDs, partType)
	if err != nil {
		return nil, fmt.Errorf("получения списка деталей: %w", err)
	}
	return parts, nil
}
