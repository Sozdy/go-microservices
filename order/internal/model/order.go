package model

import (
	"time"

	"github.com/google/uuid"
)

type (
	OrderStatus   string
	PaymentMethod string
)

type PartType string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"

	PaymentMethodUnspecified   PaymentMethod = "UNSPECIFIED"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

const (
	PartTypeHull   PartType = "hull"
	PartTypeEngine PartType = "engine"
	PartTypeShield PartType = "shield"
	PartTypeWeapon PartType = "weapon"
)

type Order struct {
	UUID            uuid.UUID
	OrderItems      []OrderItem
	TotalPrice      int64
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
}

type OrderItem struct {
	PartUUID uuid.UUID
	PartType PartType
	Price    int64
}
