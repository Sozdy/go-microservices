package part

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/Sozdy/go-microservices/inventory/internal/errs"
	"github.com/Sozdy/go-microservices/inventory/internal/model"
	"github.com/Sozdy/go-microservices/inventory/internal/repository/converter"
	"github.com/Sozdy/go-microservices/inventory/internal/repository/record"
)

func (r *repo) ListParts(
	ctx context.Context,
	partUUIDs []string,
	partType model.PartType,
) ([]*model.Part, error) {
	query := r.builder.
		Select(allPartColumns...).
		From(partsTable)

	if len(partUUIDs) > 0 {
		query = query.Where(sq.Eq{colUUID: partUUIDs})
	} else {
		query = query.OrderBy(colName + " ASC")
		if partType != model.PartTypeUnspecified {
			query = query.Where(sq.Eq{colPartType: string(partType)})
		}
	}

	querySql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("построить запрос: %w", err)
	}

	rows, err := r.pool.Query(ctx, querySql, args...)
	if err != nil {
		return nil, fmt.Errorf("запрос деталей: %w", err)
	}

	records, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.Part])
	if err != nil {
		return nil, fmt.Errorf("маппинг деталей: %w", err)
	}

	if len(partUUIDs) > 0 {
		if len(records) < len(partUUIDs) {
			return nil, errs.ErrPartNotFound
		}
		records = reorderByUUIDs(records, partUUIDs)
	}

	result := make([]*model.Part, 0, len(records))
	for i := range records {
		result = append(result, converter.PartToModel(&records[i]))
	}
	return result, nil
}

func reorderByUUIDs(records []record.Part, partUUIDs []string) []record.Part {
	byUUID := make(map[string]record.Part, len(records))
	for _, rec := range records {
		byUUID[rec.UUID] = rec
	}

	ordered := make([]record.Part, 0, len(partUUIDs))
	for _, id := range partUUIDs {
		if rec, ok := byUUID[id]; ok {
			ordered = append(ordered, rec)
		}
	}
	return ordered
}
