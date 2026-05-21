package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/Sozdy/go-microservices/order/pkg/app"
	inventoryv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/payment/v1"
)

const (
	httpAddress = "127.0.0.1:8080"

	// Таймауты для HTTP-сервера.
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second

	shutdownTimeout = 10 * time.Second

	inventoryServiceAddress = "localhost:50051"
	paymentServiceAddress   = "localhost:50052"
	keepaliveTime           = 10 * time.Second
	keepaliveTimeout        = 3 * time.Second
)

func main() {
	if err := run(); err != nil {
		slog.Error("не удалось запустить OrderService", "error", err)
		os.Exit(1)
	}
}

func run() error {
	paymentConn, err := grpc.NewClient(
		paymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                keepaliveTime,
			Timeout:             keepaliveTimeout,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return fmt.Errorf("подключение к PaymentService: %w", err)
	}
	defer paymentConn.Close()

	paymentClient := paymentv1.NewPaymentServiceClient(paymentConn)

	inventoryConn, err := grpc.NewClient(
		inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                keepaliveTime,
			Timeout:             keepaliveTimeout,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return fmt.Errorf("подключение к InventoryService: %w", err)
	}
	defer inventoryConn.Close()

	inventoryClient := inventoryv1.NewInventoryServiceClient(inventoryConn)

	// Создаём подключение к БД через pgxpool.
	ctx := context.Background()
	orderDSN := "postgres://order-service-user:order-service-password@localhost:5432/order-service?sslmode=disable" //nolint:gosec // учебный проект
	pool, err := pgxpool.New(ctx, orderDSN)
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

	// Создаём Transaction Manager для pgx
	txManager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		return fmt.Errorf("создание transaction manager: %w", err)
	}

	handler, err := app.NewHTTPHandler(pool, txManager, inventoryClient, paymentClient)
	if err != nil {
		return fmt.Errorf("создание HTTP handler: %w", err)
	}

	server := &http.Server{
		Addr:              httpAddress,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атаки
		ReadTimeout:       readTimeout,       // Лимит на чтение всего запроса
		WriteTimeout:      writeTimeout,      // Лимит на запись ответа
		IdleTimeout:       idleTimeout,       // Таймаут keep-alive соединений
	}

	var lc net.ListenConfig
	listener, err := lc.Listen(context.Background(), "tcp", httpAddress)
	if err != nil {
		return fmt.Errorf("создание listener: %w", err)
	}
	defer listener.Close()

	slog.Info("запуск OrderService", "адрес", httpAddress)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	serveErrCh := make(chan error, 1)
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErrCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		slog.Info("🛑 завершение работы сервера...")

		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancelShutdown()

		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("❌ ошибка при остановке сервера", "error", err)
		}

		slog.Info("✅ сервер остановлен")
		return nil

	case err := <-serveErrCh:
		return fmt.Errorf("HTTP сервер завершился с ошибкой: %w", err)
	}
}
