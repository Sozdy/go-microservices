package part

import (
	"context"
	"sort"

	"github.com/google/uuid"

	errs "github.com/Sozdy/go-microservices/inventory/internal/errors"
	"github.com/Sozdy/go-microservices/inventory/internal/model"
	"github.com/Sozdy/go-microservices/inventory/internal/repository/converter"
)

func (r *repo) ListParts(ctx context.Context, partUUIDs []string, partType model.PartType) ([]*model.Part, error) {
	result := make([]*model.Part, 0)

	r.mu.Lock()
	defer r.mu.Unlock()

	if len(partUUIDs) > 0 {
		for _, partUUID := range partUUIDs {
			if err := uuid.Validate(partUUID); err != nil {
				return nil, errs.ErrInvalidUUID
			}

			part, ok := r.parts[partUUID]
			if !ok {
				return nil, errs.ErrPartNotFound
			}

			result = append(result, converter.PartToModel(&part))
		}
	} else {
		for _, part := range r.parts {
			if partType != model.PartTypeUnspecified &&
				part.PartType != string(partType) {
				continue
			}

			result = append(result, converter.PartToModel(&part))
		}

		sort.Slice(result, func(i, j int) bool {
			return result[i].Name < result[j].Name
		})
	}

	return result, nil
}
