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

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Sozdy/go-microservices/inventory/pkg/app"
)

const grpcAddress = "127.0.0.1:50051"

func main() {
	if err := run(); err != nil {
		slog.Error("не удалось запустить InventoryService", "error", err)
		os.Exit(1)
	}
}

func run() error {
	// Создаём подключение к БД через pgxpool.
	ctx := context.Background()
	inventoryDSN := "postgres://inventory-service-user:inventory-service-password@localhost:5433/inventory-service?sslmode=disable" //nolint:gosec // учебный проект

	pool, err := pgxpool.New(ctx, inventoryDSN)
	if err != nil {
		return fmt.Errorf("создание пула соединений: %w", err)
	}
	defer pool.Close()

	// Проверяем соединение
	err = pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("проверка соединения с БД: %w", err)
	}

	slog.Info("подключение к PostgreSQL установлено")

	var lc net.ListenConfig
	listener, err := lc.Listen(context.Background(), "tcp", grpcAddress)
	if err != nil {
		return fmt.Errorf("создание listener: %w", err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer(app.Interceptors()...)
	app.RegisterServices(grpcServer, pool)

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
		return nil

	case err := <-serveErrCh:
		return fmt.Errorf("gRPC сервер завершился с ошибкой: %w", err)
	}
}
