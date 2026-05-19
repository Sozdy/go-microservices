package part

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

const partsTable = "parts"

const (
	colUUID          = "uuid"
	colName          = "name"
	colDescription   = "description"
	colPartType      = "part_type"
	colPrice         = "price"
	colStockQuantity = "stock_quantity"
	colCreatedAt     = "created_at"
)

var allPartColumns = []string{
	colUUID,
	colName,
	colDescription,
	colPartType,
	colPrice,
	colStockQuantity,
	colCreatedAt,
}

type repo struct {
	pool    *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewRepository(pool *pgxpool.Pool) *repo {
	return &repo{
		pool:    pool,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
