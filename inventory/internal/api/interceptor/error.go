package interceptor

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Sozdy/go-microservices/inventory/internal/errs"
)

func codeToGRPC(c errs.Code) codes.Code {
	switch c {
	case errs.CodeNotFound:
		return codes.NotFound
	case errs.CodeInvalidArgument:
		return codes.InvalidArgument
	case errs.CodeConflict:
		return codes.AlreadyExists
	case errs.CodeFailedPrecondition:
		return codes.FailedPrecondition
	case errs.CodeUnavailable:
		return codes.Unavailable
	default:
		return codes.Internal
	}
}

func UnaryErrorInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		if _, ok := status.FromError(err); ok {
			return nil, err
		}

		switch {
		case errors.Is(err, context.Canceled):
			return nil, status.Error(codes.Canceled, "запрос отменён")
		case errors.Is(err, context.DeadlineExceeded):
			return nil, status.Error(codes.DeadlineExceeded, "истёк таймаут запроса")
		}

		code := errs.CodeOf(err)
		if code == errs.CodeInternal {
			log.Error("внутренняя ошибка", "method", info.FullMethod, "err", err)
		}

		return nil, status.Error(codeToGRPC(code), errs.ClientMessage(err))
	}
}
