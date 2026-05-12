package converter

import (
	"github.com/google/uuid"

	"github.com/Sozdy/go-microservices/order/internal/service/order"
	orderv1 "github.com/Sozdy/go-microservices/shared/pkg/openapi/order/v1"
)

func CreateOrderInFromRequest(req *orderv1.CreateOrderRequest) *order.CreateOrderIn {
	var shieldUUID *uuid.UUID
	if v, ok := req.ShieldUUID.Get(); ok {
		shieldUUID = &v
	}

	var weaponUUID *uuid.UUID
	if v, ok := req.WeaponUUID.Get(); ok {
		weaponUUID = &v
	}

	return &order.CreateOrderIn{
		HullUUID:   req.HullUUID,
		EngineUUID: req.EngineUUID,
		ShieldUUID: shieldUUID,
		WeaponUUID: weaponUUID,
	}
}

func CreateOrderOutToResponse(out *order.CreateOrderOut) orderv1.CreateOrderRes {
	return &orderv1.CreateOrderResponse{
		OrderUUID:  out.OrderUUID,
		TotalPrice: out.TotalPrice,
	}
}
