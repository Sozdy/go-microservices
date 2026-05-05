package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	errs "github.com/Sozdy/go-microservices/order/internal/errors"
	"github.com/Sozdy/go-microservices/order/internal/model"
	inventoryModel "github.com/Sozdy/go-microservices/order/internal/model/ports/inventory"
)

type CreateOrderIn struct {
	HullUUID   uuid.UUID
	EngineUUID uuid.UUID
	ShieldUUID *uuid.UUID
	WeaponUUID *uuid.UUID
}

func (in *CreateOrderIn) collectPartUUIDs() []string {
	var partUUIDs []string
	partUUIDs = append(partUUIDs, in.HullUUID.String())
	partUUIDs = append(partUUIDs, in.EngineUUID.String())

	if in.ShieldUUID != nil && *in.ShieldUUID != uuid.Nil {
		partUUIDs = append(partUUIDs, in.ShieldUUID.String())
	}

	if in.WeaponUUID != nil && *in.WeaponUUID != uuid.Nil {
		partUUIDs = append(partUUIDs, in.WeaponUUID.String())
	}

	return partUUIDs
}

type CreateOrderOut struct {
	OrderUUID  uuid.UUID
	TotalPrice int64
}

func (s *orderService) CreateOrder(ctx context.Context, in *CreateOrderIn) (*CreateOrderOut, error) {
	partUUIDs := in.collectPartUUIDs()

	listPartsRes, err := s.inventoryClient.ListParts(ctx, &inventoryModel.ListPartsRequest{
		PartType: inventoryModel.PART_TYPE_UNSPECIFIED,
		Uuids:    partUUIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("получить детали из inventory: %w", err)
	}

	totalPrice, err := calculateTotalPrice(listPartsRes.Parts)
	if err != nil {
		return nil, fmt.Errorf("рассчитать стоимость заказа: %w", err)
	}

	order := model.Order{
		OrderUUID:       uuid.New(),
		HullUUID:        in.HullUUID,
		EngineUUID:      in.EngineUUID,
		ShieldUUID:      in.ShieldUUID,
		WeaponUUID:      in.WeaponUUID,
		TotalPrice:      totalPrice,
		TransactionUUID: nil,
		PaymentMethod:   nil,
		Status:          model.OrderStatusPendingPayment,
		CreatedAt:       time.Now(),
	}

	if err := s.orderRepository.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("сохранить заказ: %w", err)
	}

	return &CreateOrderOut{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}

func calculateTotalPrice(listParts []inventoryModel.Part) (int64, error) {
	var totalPrice int64
	for _, part := range listParts {
		if part.StockQuantity <= 0 {
			return 0, fmt.Errorf("деталь %s: %w", part.UUID, errs.ErrPartUnavailable)
		}
		totalPrice += part.Price
	}

	return totalPrice, nil
}
