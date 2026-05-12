package v1

import (
	paymentv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/payment/v1"
)

type client struct {
	grpcClient paymentv1.PaymentServiceClient
}

func NewClient(paymentClient paymentv1.PaymentServiceClient) *client {
	return &client{
		grpcClient: paymentClient,
	}
}
