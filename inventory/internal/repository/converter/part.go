package converter

import (
	"github.com/Sozdy/go-microservices/inventory/internal/model"
	"github.com/Sozdy/go-microservices/inventory/internal/repository/record"
)

func PartToRecord(part *model.Part) *record.Part {
	return &record.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		PartType:      string(part.PartType),
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		CreatedAt:     part.CreatedAt,
	}
}

func PartToModel(part *record.Part) *model.Part {
	return &model.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		PartType:      model.PartType(part.PartType),
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		CreatedAt:     part.CreatedAt,
	}
}
