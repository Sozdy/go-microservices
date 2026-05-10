package v1

import (
	"context"
	"testing"

	"github.com/Sozdy/go-microservices/payment/internal/api/payment/v1/mocks"
)

type apiFixture struct {
	ctx            context.Context
	api            *api
	paymentService *mocks.PaymentService
}

func newAPIFixture(t *testing.T) *apiFixture {
	t.Helper()

	paymentService := mocks.NewPaymentService(t)

	return &apiFixture{
		ctx:            context.Background(),
		paymentService: paymentService,
		api:            NewApi(paymentService),
	}
}
