package v1

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Sozdy/go-microservices/order/internal/client/grpc/inventory/v1/converter"
	errs "github.com/Sozdy/go-microservices/order/internal/errors"
	"github.com/Sozdy/go-microservices/order/internal/model/ports/inventory"
)

func (c *client) ListParts(ctx context.Context, req *inventory.ListPartsRequest) (*inventory.ListPartsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.grpcClient.ListParts(ctx, converter.ListPartsRequestToProto(req))
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, errs.ErrPartNotFound
			case codes.InvalidArgument:
				return nil, errs.ErrInvalidPartFilter
			case codes.Unavailable, codes.DeadlineExceeded:
				return nil, errs.ErrInventoryUnavailable
			}
		}
		return nil, fmt.Errorf("сервис inventory метод ListParts: %w", err)
	}

	return converter.ListPartsResponseFromProto(resp), nil
}
