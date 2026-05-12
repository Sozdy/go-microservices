package app

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Sozdy/go-microservices/order/internal/api"
	v1 "github.com/Sozdy/go-microservices/order/internal/api/order/v1"
	inventoryclient "github.com/Sozdy/go-microservices/order/internal/client/grpc/inventory/v1"
	paymentclient "github.com/Sozdy/go-microservices/order/internal/client/grpc/payment/v1"
	orderRepo "github.com/Sozdy/go-microservices/order/internal/repository/order"
	"github.com/Sozdy/go-microservices/order/internal/service/order"
	orderv1 "github.com/Sozdy/go-microservices/shared/pkg/openapi/order/v1"
	protoinventoryv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/inventory/v1"
	protopaymentv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/payment/v1"
)

func NewHTTPHandler(inventoryGRPC protoinventoryv1.InventoryServiceClient, paymentGRPC protopaymentv1.PaymentServiceClient) (http.Handler, error) {
	orderApi := v1.NewApi(
		order.NewOrderService(
			inventoryclient.NewClient(inventoryGRPC),
			paymentclient.NewClient(paymentGRPC),
			orderRepo.NewRepository(),
		),
	)

	orderHandler, err := orderv1.NewServer(orderApi, orderv1.WithErrorHandler(api.OgenErrorHandler))
	if err != nil {
		slog.Error("ошибка создания сервера OpenAPI", "error", err)
		return nil, fmt.Errorf("ошибка создания сервера OpenAPI: %w", err)
	}

	return orderHandler, nil
}
