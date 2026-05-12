package v1

import (
	"context"
	"testing"

	"github.com/Sozdy/go-microservices/inventory/internal/api/inventory/v1/mocks"
)

type apiFixture struct {
	ctx     context.Context
	api     *api
	service *mocks.InventoryService
}

func newApiFixture(t *testing.T) apiFixture {
	t.Helper()

	inventoryService := mocks.NewInventoryService(t)

	return apiFixture{
		ctx:     context.Background(),
		service: inventoryService,
		api:     NewApi(inventoryService),
	}
}
