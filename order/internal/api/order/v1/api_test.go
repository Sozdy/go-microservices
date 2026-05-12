package v1

import (
	"context"
	"testing"

	"github.com/Sozdy/go-microservices/order/internal/api/order/v1/mocks"
)

type apiFixture struct {
	ctx          context.Context
	api          *api
	orderService *mocks.OrderService
}

func newApiFixture(t *testing.T) *apiFixture {
	t.Helper()

	orderService := mocks.NewOrderService(t)

	return &apiFixture{
		ctx:          context.Background(),
		api:          NewApi(orderService),
		orderService: orderService,
	}
}
