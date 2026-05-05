package converter

import (
	"time"

	"github.com/Sozdy/go-microservices/inventory/internal/model"
	"github.com/Sozdy/go-microservices/inventory/internal/repository/record"
)

func PartToRecord(part *model.Part) *record.Part {
	return &record.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      (string)(part.PartType),
		StockQuantity: part.StockQuantity,
		CreatedAt:     part.CreatedAt.Unix(),
	}
}

func PartToModel(part *record.Part) *model.Part {
	return &model.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      (model.PartType)(part.PartType),
		StockQuantity: part.StockQuantity,
		CreatedAt:     time.Unix(part.CreatedAt, 0),
	}
}
