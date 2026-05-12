package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

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
		slog.Error("не удалось подключиться к PaymentService", "error", err)
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
		slog.Error("не удалось подключиться к InventoryService", "error", err)
	}
	defer inventoryConn.Close()

	inventoryClient := inventoryv1.NewInventoryServiceClient(inventoryConn)

	handler, err := app.NewHTTPHandler(inventoryClient, paymentClient)
	if err != nil {
		panic(err)
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
		slog.Error("не удалось создать listener", "error", err)
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

	case err := <-serveErrCh:
		slog.Error("🛑 HTTP сервер завершился с ошибкой", "error", err)
		return
	}
}
