package service

import (
	"context"
	"sort"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/inventory/v1"
)

// Part представляет деталь космического корабля.
type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         int64 // В копейках
	PartType      inventoryv1.PartType
	StockQuantity int64
	CreatedAt     *timestamppb.Timestamp
}

// InventoryServer реализует gRPC сервис.
type InventoryServer struct {
	inventoryv1.UnimplementedInventoryServiceServer
	parts map[uuid.UUID]Part
}

// NewInventoryServer создаёт сервер с предзагруженными seed-данными.
func NewInventoryServer() *InventoryServer {
	now := timestamppb.Now()

	return &InventoryServer{
		parts: map[uuid.UUID]Part{
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440001",
				Name:          "Алюминиевый корпус",
				Description:   "Лёгкий корпус для небольших кораблей",
				Price:         500000, // 5000₽
				PartType:      inventoryv1.PartType_PART_TYPE_HULL,
				StockQuantity: 10,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440002",
				Name:          "Титановый корпус",
				Description:   "Прочный корпус для средних кораблей",
				Price:         1500000, // 15000₽
				PartType:      inventoryv1.PartType_PART_TYPE_HULL,
				StockQuantity: 5,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440003",
				Name:          "Ионный двигатель C",
				Description:   "Базовый ионный двигатель класса C",
				Price:         300000, // 3000₽
				PartType:      inventoryv1.PartType_PART_TYPE_ENGINE,
				StockQuantity: 8,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440004",
				Name:          "Ионный двигатель B",
				Description:   "Улучшенный ионный двигатель класса B",
				Price:         800000, // 8000₽
				PartType:      inventoryv1.PartType_PART_TYPE_ENGINE,
				StockQuantity: 3,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440005",
				Name:          "Энергетический щит",
				Description:   "Стандартный энергетический щит",
				Price:         400000, // 4000₽
				PartType:      inventoryv1.PartType_PART_TYPE_SHIELD,
				StockQuantity: 6,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440006"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440006",
				Name:          "Лазерная пушка",
				Description:   "Точная лазерная пушка",
				Price:         250000, // 2500₽
				PartType:      inventoryv1.PartType_PART_TYPE_WEAPON,
				StockQuantity: 7,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440007"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440007",
				Name:          "Плазменный корпус",
				Description:   "Экспериментальный корпус (нет на складе)",
				Price:         2000000, // 20000₽
				PartType:      inventoryv1.PartType_PART_TYPE_HULL,
				StockQuantity: 0,
				CreatedAt:     now,
			},
		},
	}
}

// GetPart возвращает деталь по UUID.
func (s *InventoryServer) GetPart(
	ctx context.Context,
	req *inventoryv1.GetPartRequest,
) (*inventoryv1.GetPartResponse, error) {
	partUUID := req.GetUuid()
	if partUUID == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid не может быть пустым")
	}

	parsedUUID, err := uuid.Parse(partUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "неверный формат uuid: "+err.Error())
	}

	part, ok := s.parts[parsedUUID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "деталь не найдена: %s", partUUID)
	}

	return &inventoryv1.GetPartResponse{
		Part: toProtoPart(part),
	}, nil
}

// ListParts возвращает список деталей с опциональной фильтрацией по типу.
func (s *InventoryServer) ListParts(
	ctx context.Context,
	req *inventoryv1.ListPartsRequest,
) (*inventoryv1.ListPartsResponse, error) {
	reqUUIDs := req.GetUuids()
	partType := req.GetPartType()

	result := make([]*inventoryv1.Part, 0)

	if len(reqUUIDs) > 0 {
		if partType != inventoryv1.PartType_PART_TYPE_UNSPECIFIED {
			return nil, status.Error(
				codes.InvalidArgument,
				"part_type не должен передаваться вместе с uuids",
			)
		}

		for _, partUUID := range reqUUIDs {
			parsedUUID, err := uuid.Parse(partUUID)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, "неверный формат uuid: "+err.Error())
			}

			part, ok := s.parts[parsedUUID]
			if !ok {
				return nil, status.Errorf(codes.NotFound, "деталь не найдена: %s", partUUID)
			}

			result = append(result, toProtoPart(part))
		}
	} else {
		for _, part := range s.parts {
			if partType != inventoryv1.PartType_PART_TYPE_UNSPECIFIED &&
				part.PartType != partType {
				continue
			}

			result = append(result, toProtoPart(part))
		}

		sort.Slice(result, func(i, j int) bool {
			return result[i].Name < result[j].Name
		})
	}

	return &inventoryv1.ListPartsResponse{
		Parts: result,
	}, nil
}

func toProtoPart(part Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      part.PartType,
		StockQuantity: part.StockQuantity,
		CreatedAt:     part.CreatedAt,
	}
}
