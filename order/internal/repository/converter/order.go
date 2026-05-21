package converter

import (
	"github.com/google/uuid"

	"github.com/Sozdy/go-microservices/order/internal/model"
	"github.com/Sozdy/go-microservices/order/internal/repository/record"
)

func OrderToRecord(order model.Order) record.Order {
	var paymentMethod *string
	if order.PaymentMethod != nil {
		paymentMethod = new(string(*order.PaymentMethod))
	}

	return record.Order{
		UUID:            order.UUID,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          string(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}

func OrderItemsToRecord(order model.Order) []record.OrderItem {
	items := make([]record.OrderItem, 0, len(order.OrderItems))

	for _, orderItem := range order.OrderItems {
		items = append(items, record.OrderItem{
			UUID:      uuid.New(),
			OrderUUID: order.UUID,
			PartUUID:  orderItem.PartUUID,
			PartType:  record.PartType(orderItem.PartType),
			Price:     orderItem.Price,
		})
	}

	return items
}

func OrderFromRecord(order record.Order, orderItems []record.OrderItem) model.Order {
	items := make([]model.OrderItem, 0, len(orderItems))
	for _, orderItem := range orderItems {
		items = append(items, model.OrderItem{
			PartUUID: orderItem.PartUUID,
			PartType: model.PartType(orderItem.PartType),
			Price:    orderItem.Price,
		})
	}

	var TotalPrice int64
	for _, orderItem := range orderItems {
		TotalPrice += orderItem.Price
	}

	var paymentMethod *model.PaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = new(model.PaymentMethod(*order.PaymentMethod))
	}

	return model.Order{
		UUID:            order.UUID,
		OrderItems:      items,
		TotalPrice:      TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          model.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}
