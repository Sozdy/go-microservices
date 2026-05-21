package part

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/Sozdy/go-microservices/inventory/internal/errs"
	"github.com/Sozdy/go-microservices/inventory/internal/model"
	"github.com/Sozdy/go-microservices/inventory/internal/repository/converter"
	"github.com/Sozdy/go-microservices/inventory/internal/repository/record"
)

func (r *repo) GetPart(ctx context.Context, partUUID string) (*model.Part, error) {
	query := r.builder.
		Select(allPartColumns...).
		From(partsTable).
		Where(sq.Eq{colUUID: partUUID})

	querySql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("построить запрос: %w", err)
	}

	rows, err := r.pool.Query(ctx, querySql, args...)
	if err != nil {
		return nil, fmt.Errorf("запрос детали: %w", err)
	}

	rec, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[record.Part])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrPartNotFound
		}
		return nil, fmt.Errorf("маппинг детали: %w", err)
	}

	return converter.PartToModel(&rec), nil
}
