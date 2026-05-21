package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Sozdy/go-microservices/payment/pkg/app"
)

const grpcAddress = "127.0.0.1:50052"

func main() {
	if err := run(); err != nil {
		slog.Error("не удалось запустить PaymentService", "error", err)
		os.Exit(1)
	}
}

func run() error {
	var lc net.ListenConfig
	listener, err := lc.Listen(context.Background(), "tcp", grpcAddress)
	if err != nil {
		return fmt.Errorf("создание listener: %w", err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer(app.Interceptors()...)
	app.RegisterServices(grpcServer)

	// Включаем reflection для postman/grpcurl
	reflection.Register(grpcServer)

	slog.Info("запуск PaymentService", "адрес", grpcAddress)

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
		return nil

	case err := <-serveErrCh:
		return fmt.Errorf("gRPC сервер завершился с ошибкой: %w", err)
	}
}
