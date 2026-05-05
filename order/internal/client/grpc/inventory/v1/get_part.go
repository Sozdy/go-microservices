package v1

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Sozdy/go-microservices/order/internal/client/grpc/inventory/v1/converter"
	errs "github.com/Sozdy/go-microservices/order/internal/errors"
	"github.com/Sozdy/go-microservices/order/internal/model/ports/inventory"
	inventoryv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/inventory/v1"
)

func (c *client) GetPart(ctx context.Context, partUUID uuid.UUID) (*inventory.Part, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.grpcClient.GetPart(ctx, &inventoryv1.GetPartRequest{
		Uuid: partUUID.String(),
	})
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
		return nil, fmt.Errorf("сервис inventory метод GetPart: %w", err)
	}

	return converter.PartFromProto(resp.Part), nil
}
