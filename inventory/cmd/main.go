package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Sozdy/go-microservices/inventory/pkg/app"
)

const grpcAddress = "127.0.0.1:50051"

func main() {
	// Создаём подключение к БД через pgxpool.
	ctx := context.Background()
	inventoryDSN := "postgres://inventory-service-user:inventory-service-password@localhost:5433/inventory-service?sslmode=disable" //nolint:gosec // учебный проект

	pool, err := pgxpool.New(ctx, inventoryDSN)
	if err != nil {
		slog.Error("создание пула соединений", "error", err)
	}
	defer pool.Close()

	// Проверяем соединение
	err = pool.Ping(ctx)
	if err != nil {
		slog.Error("проверка соединения с БД", "error", err)
	}

	slog.Info("подключение к PostgreSQL установлено")

	var lc net.ListenConfig
	listener, err := lc.Listen(context.Background(), "tcp", grpcAddress)
	if err != nil {
		slog.Error("не удалось создать listener", "error", err)
		panic(err)
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

	case err := <-serveErrCh:
		slog.Error("🛑 gRPC сервер завершился с ошибкой", "error", err)
		return
	}
}
