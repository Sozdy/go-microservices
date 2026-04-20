package handler

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	orderv1 "github.com/Sozdy/go-microservices/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/payment/v1"
)

// Order представляет заказ на постройку космического корабля.
type Order struct {
	OrderUUID       uuid.UUID
	HullUUID        uuid.UUID
	EngineUUID      uuid.UUID
	ShieldUUID      *uuid.UUID // опциональный
	WeaponUUID      *uuid.UUID // опциональный
	TotalPrice      int64      // в копейках
	TransactionUUID *uuid.UUID
	PaymentMethod   *string
	Status          string // PENDING_PAYMENT, PAID, CANCELLED
	CreatedAt       time.Time
}

// OrderStore — хранилище заказов (in-memory).
type OrderStore struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]Order
}

// NewOrderStore создаёт новое пустое хранилище заказов.
func NewOrderStore() *OrderStore {
	return &OrderStore{
		orders: make(map[uuid.UUID]Order),
	}
}

// OrderHandler реализует интерфейс orderv1. Handler, сгенерированный ogen.
type OrderHandler struct {
	orderv1.UnimplementedHandler
	inventoryClient inventoryv1.InventoryServiceClient
	paymentClient   paymentv1.PaymentServiceClient
	store           *OrderStore
}

// NewOrderHandler создаёт новый обработчик заказов.
func NewOrderHandler(
	inventoryClient inventoryv1.InventoryServiceClient,
	paymentClient paymentv1.PaymentServiceClient,
	store *OrderStore,
) *OrderHandler {
	return &OrderHandler{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		store:           store,
	}
}

// SetupServer создаёт OpenAPI сервер на основе обработчика.
func SetupServer(h *OrderHandler) (*orderv1.Server, error) {
	return orderv1.NewServer(h)
}

// GetOrder реализует операцию getOrder (пример реализации).
// GET /api/v1/orders/{order_uuid}.
func (h *OrderHandler) GetOrder(_ context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	// 1. Найти заказ в store (с блокировкой для thread-safety)
	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()

	// 2. Если не найден — вернуть 404
	if !ok {
		return &orderv1.GetOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	// 3. Преобразовать в DTO и вернуть
	var shieldUUID orderv1.OptNilUUID
	if order.ShieldUUID != nil {
		shieldUUID = orderv1.NewOptNilUUID(*order.ShieldUUID)
	}

	var weaponUUID orderv1.OptNilUUID
	if order.WeaponUUID != nil {
		weaponUUID = orderv1.NewOptNilUUID(*order.WeaponUUID)
	}

	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
	}

	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethod(*order.PaymentMethod))
	}

	return &orderv1.OrderDto{
		OrderUUID:       order.OrderUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderv1.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}, nil
}

// CreateOrder реализует операцию createOrder
// POST /api/v1/orders.
func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	// 1. Валидация: hull_uuid и engine_uuid обязательны
	if req.GetHullUUID() == uuid.Nil {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "hull_uuid не валиден",
		}, nil
	}
	if req.GetEngineUUID() == uuid.Nil {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "engine_uuid не валиден",
		}, nil
	}

	// 2. Получить детали через InventoryService
	var partUUIDs []string
	partUUIDs = append(partUUIDs, req.GetHullUUID().String())
	partUUIDs = append(partUUIDs, req.GetEngineUUID().String())
	if req.GetShieldUUID().Set && req.GetShieldUUID().Value != uuid.Nil {
		partUUIDs = append(partUUIDs, req.GetShieldUUID().Value.String())
	}
	if req.GetWeaponUUID().Set && req.GetWeaponUUID().Value != uuid.Nil {
		partUUIDs = append(partUUIDs, req.GetWeaponUUID().Value.String())
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	partsResponse, err := h.inventoryClient.ListParts(ctx, &inventoryv1.ListPartsRequest{
		PartType: inventoryv1.PartType_PART_TYPE_UNSPECIFIED,
		Uuids:    partUUIDs,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return &orderv1.CreateOrderNotFound{Code: http.StatusNotFound, Message: err.Error()}, nil
			case codes.InvalidArgument:
				return &orderv1.CreateOrderBadRequest{Code: http.StatusBadRequest, Message: err.Error()}, nil
			default:
				return &orderv1.CreateOrderInternalServerError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
			}
		}
	}

	// 3. Проверить stock_quantity > 0
	// 4. Вычислить total_price
	var totalPrice int64
	for _, part := range partsResponse.Parts {
		if part.StockQuantity <= 0 {
			return &orderv1.CreateOrderConflict{
				Code:    http.StatusConflict,
				Message: "деталь " + part.Uuid + " недоступна (stock_quantity <= 0)",
			}, nil
		}
		totalPrice += part.Price
	}

	// 5. Сгенерировать order_uuid (UUID v4)
	orderUUID := uuid.New() // 6. Создать заказ со статусом PENDING_PAYMENT
	order := Order{
		OrderUUID:       orderUUID,
		EngineUUID:      req.GetEngineUUID(),
		HullUUID:        req.GetHullUUID(),
		ShieldUUID:      new(req.GetShieldUUID().Value),
		WeaponUUID:      new(req.GetWeaponUUID().Value),
		TotalPrice:      totalPrice,
		TransactionUUID: nil,
		PaymentMethod:   nil,
		Status:          string(orderv1.OrderStatusPENDINGPAYMENT),
		CreatedAt:       time.Now(),
	}

	// 7. Сохранить в store
	h.store.mu.Lock()
	h.store.orders[orderUUID] = order
	h.store.mu.Unlock()

	// 8. Вернуть order_uuid и total_price
	return &orderv1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}

// PayOrder реализует операцию payOrder
// POST /api/v1/orders/{order_uuid}/pay.
func (h *OrderHandler) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	// 1. Найти заказ в store
	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()
	if !ok {
		return &orderv1.PayOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	// 2. Проверить статус == PENDING_PAYMENT
	if order.Status != string(orderv1.OrderStatusPENDINGPAYMENT) {
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: "статус заказа должен быть PENDING_PAYMENT для оплаты",
		}, nil
	}

	var paymentMethod paymentv1.PaymentMethod
	switch req.PaymentMethod {
	case orderv1.PaymentMethodCARD:
		paymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderv1.PaymentMethodSBP:
		paymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderv1.PaymentMethodCREDITCARD:
		paymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderv1.PaymentMethodINVESTORMONEY:
		paymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return &orderv1.PayOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "неизвестный payment_method",
		}, nil
	}

	// 3. Вызвать h.paymentClient.PayOrder для обработки платежа
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	payOrderResponse, err := h.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     order.OrderUUID.String(),
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return &orderv1.PayOrderNotFound{Code: http.StatusNotFound, Message: err.Error()}, nil
			case codes.InvalidArgument:
				return &orderv1.PayOrderBadRequest{Code: http.StatusBadRequest, Message: err.Error()}, nil
			default:
				return &orderv1.PayOrderInternalServerError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
			}
		}
	}

	// 4. Обновить статус на PAID и сохранить transaction_uuid
	transactionUUID, err := uuid.Parse(payOrderResponse.TransactionUuid)
	if err != nil {
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "неверный формат transaction_uuid от PaymentService",
		}, nil
	}
	orderStatus := string(orderv1.OrderStatusPAID)
	order.Status = orderStatus
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = new(string(req.PaymentMethod))

	h.store.mu.Lock()
	h.store.orders[order.OrderUUID] = order
	h.store.mu.Unlock()

	// 5. Вернуть transaction_uuid
	return &orderv1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}

// CancelOrder реализует операцию cancelOrder
// POST /api/v1/orders/{order_uuid}/cancel.
func (h *OrderHandler) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	// 1. Найти заказ в store
	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()
	if !ok {
		return &orderv1.CancelOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	// 2. Проверить статус == PENDING_PAYMENT
	if order.Status != string(orderv1.OrderStatusPENDINGPAYMENT) {
		return &orderv1.CancelOrderConflict{
			Code:    http.StatusConflict,
			Message: "статус заказа должен быть PENDING_PAYMENT для отмены",
		}, nil
	}

	// 3. Обновить статус на CANCELLED
	order.Status = string(orderv1.OrderStatusCANCELLED)

	h.store.mu.Lock()
	h.store.orders[order.OrderUUID] = order
	h.store.mu.Unlock()

	// 4. Вернуть success
	return &orderv1.CancelOrderResponse{}, nil
}
