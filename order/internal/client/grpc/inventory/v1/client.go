package v1

import (
	inventoryv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/inventory/v1"
)

type client struct {
	grpcClient inventoryv1.InventoryServiceClient
}

func NewClient(inventoryClient inventoryv1.InventoryServiceClient) *client {
	return &client{
		grpcClient: inventoryClient,
	}
}
