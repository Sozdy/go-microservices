package converter

import (
	"github.com/google/uuid"

	"github.com/Sozdy/go-microservices/order/internal/model/ports/inventory"
	inventoryv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/inventory/v1"
)

func PartTypeToProto(partType inventory.PartType) inventoryv1.PartType {
	switch partType {
	case inventory.PART_TYPE_HULL:
		return inventoryv1.PartType_PART_TYPE_HULL
	case inventory.PART_TYPE_ENGINE:
		return inventoryv1.PartType_PART_TYPE_ENGINE
	case inventory.PART_TYPE_SHIELD:
		return inventoryv1.PartType_PART_TYPE_SHIELD
	case inventory.PART_TYPE_WEAPON:
		return inventoryv1.PartType_PART_TYPE_WEAPON
	default:
		return inventoryv1.PartType_PART_TYPE_UNSPECIFIED
	}
}

func PartTypeFromProto(partType inventoryv1.PartType) inventory.PartType {
	switch partType {
	case inventoryv1.PartType_PART_TYPE_HULL:
		return inventory.PART_TYPE_HULL
	case inventoryv1.PartType_PART_TYPE_ENGINE:
		return inventory.PART_TYPE_ENGINE
	case inventoryv1.PartType_PART_TYPE_SHIELD:
		return inventory.PART_TYPE_SHIELD
	case inventoryv1.PartType_PART_TYPE_WEAPON:
		return inventory.PART_TYPE_WEAPON
	default:
		return inventory.PART_TYPE_UNSPECIFIED
	}
}

func ListPartsRequestToProto(req *inventory.ListPartsRequest) *inventoryv1.ListPartsRequest {
	return &inventoryv1.ListPartsRequest{
		PartType: PartTypeToProto(req.PartType),
		Uuids:    req.Uuids,
	}
}

func ListPartsResponseFromProto(resp *inventoryv1.ListPartsResponse) *inventory.ListPartsResponse {
	parts := make([]inventory.Part, 0, len(resp.Parts))
	for _, part := range resp.Parts {
		parts = append(parts, *PartFromProto(part))
	}

	return &inventory.ListPartsResponse{
		Parts: parts,
	}
}

func PartFromProto(p *inventoryv1.Part) *inventory.Part {
	return &inventory.Part{
		UUID:          uuid.MustParse(p.Uuid),
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		PartType:      PartTypeFromProto(p.PartType),
		StockQuantity: p.StockQuantity,
		CreatedAt:     p.CreatedAt.AsTime(),
	}
}
