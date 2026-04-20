package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/Sozdy/go-microservices/inventory/pkg/service"
	inventoryv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/inventory/v1"
)

const grpcAddress = "127.0.0.1:50051"

const (
	grpcMaxConnectionIdle     = 15 * time.Minute
	grpcMaxConnectionAge      = 30 * time.Minute
	grpcMaxConnectionAgeGrace = 5 * time.Second
	grpcKeepaliveTime         = 5 * time.Minute
	grpcKeepaliveTimeout      = 1 * time.Second
	grpcMinPingInterval       = 5 * time.Minute
)

func main() {
	var lc net.ListenConfig
	listener, err := lc.Listen(context.Background(), "tcp", grpcAddress)
	if err != nil {
		slog.Error("не удалось создать listener", "error", err)
		panic(err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer(
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
	)

	inventoryv1.RegisterInventoryServiceServer(grpcServer, service.NewInventoryServer())

	// Включаем reflection для postman/grpcurl
	reflection.Register(grpcServer)

	slog.Info("запуск InventoryService", "адрес", grpcAddress)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	serveErrCh := make(chan error, 1)
	go func() {
		if err := grpcServer.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			serveErrCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		slog.Info("🛑 остановка gRPC сервера")
		grpcServer.GracefulStop()
		slog.Info("✅ сервер остановлен")

	case err := <-serveErrCh:
		slog.Error("🛑 gRPC сервер завершился с ошибкой", "error", err)
		panic("🛑 gRPC сервер завершился с ошибкой: " + err.Error())
	}
}
