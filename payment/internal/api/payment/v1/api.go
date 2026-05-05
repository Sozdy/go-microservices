package v1

import paymentv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/payment/v1"

type api struct {
	paymentv1.UnimplementedPaymentServiceServer
	PaymentService PaymentService
}

func NewApi(paymentService PaymentService) *api {
	return &api{
		PaymentService: paymentService,
	}
}
