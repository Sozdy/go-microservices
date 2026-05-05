package v1

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Sozdy/go-microservices/order/internal/client/grpc/payment/v1/converter"
	errs "github.com/Sozdy/go-microservices/order/internal/errors"
	paymentModel "github.com/Sozdy/go-microservices/order/internal/model/ports/payment"
)

func (c *client) PayOrder(ctx context.Context, req *paymentModel.PayOrderRequest) (*paymentModel.PayOrderResponse, error) {
	{
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		payOrderResponse, err := c.grpcClient.PayOrder(ctx, converter.PayOrderRequestToProto(req))
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				switch st.Code() {
				case codes.NotFound:
					return nil, errs.ErrPayOrderNotFound
				case codes.InvalidArgument:
					return nil, errs.ErrInvalidPayOrder
				case codes.Unavailable, codes.DeadlineExceeded:
					return nil, errs.ErrPaymentUnavailable
				}
			}
			return nil, fmt.Errorf("сервис payment метод PayOrder: %w", err)
		}

		return converter.PayOrderResponseFromProto(payOrderResponse), nil
	}
}
