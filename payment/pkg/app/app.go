package app

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	v1 "github.com/Sozdy/go-microservices/payment/internal/api/payment/v1"
	"github.com/Sozdy/go-microservices/payment/internal/service/payment"
	paymentv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/payment/v1"
)

const (
	grpcMaxConnectionIdle     = 15 * time.Minute
	grpcMaxConnectionAge      = 30 * time.Minute
	grpcMaxConnectionAgeGrace = 5 * time.Second
	grpcKeepaliveTime         = 5 * time.Minute
	grpcKeepaliveTimeout      = 1 * time.Second
	grpcMinPingInterval       = 5 * time.Minute
)

func Interceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     grpcMaxConnectionIdle,
			MaxConnectionAge:      grpcMaxConnectionAge,
			MaxConnectionAgeGrace: grpcMaxConnectionAgeGrace,
			Time:                  grpcKeepaliveTime,
			Timeout:               grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcMinPingInterval,
			PermitWithoutStream: true,
		}),
	}
}

func RegisterServices(server *grpc.Server) {
	paymentv1.RegisterPaymentServiceServer(
		server,
		v1.NewApi(
			payment.NewPaymentService(),
		),
	)
}
