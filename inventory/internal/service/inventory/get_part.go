package inventory

import (
	"context"
	"fmt"

	"github.com/Sozdy/go-microservices/inventory/internal/model"
)

func (s *partService) GetPart(ctx context.Context, partUUID string) (*model.Part, error) {
	part, err := s.partRepository.GetPart(ctx, partUUID)
	if err != nil {
		return nil, fmt.Errorf("получение part: %w", err)
	}

	return part, nil
}
