package app

import (
	"net/http"

	inventoryv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/payment/v1"
)

func NewHTTPHandler(client inventoryv1.InventoryServiceClient, client2 paymentv1.PaymentServiceClient) (http.Handler, error) {
	return nil, nil
}
