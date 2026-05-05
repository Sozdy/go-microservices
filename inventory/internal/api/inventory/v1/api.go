package v1

import inventoryv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/inventory/v1"

type api struct {
	inventoryv1.UnimplementedInventoryServiceServer
	InventoryService InventoryService
}

func NewApi(inventoryService InventoryService) *api {
	return &api{
		InventoryService: inventoryService,
	}
}
