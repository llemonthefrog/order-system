package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"order-system/internal/application"
	httpHandlers "order-system/internal/delivers/http"
	"order-system/internal/gateways"
	"order-system/internal/infrastructure/event_busses"
	"order-system/internal/infrastructure/postgres"
	"order-system/internal/workers"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	orderRepo := postgres.NewOrderRepository(db)
	taskRepo := postgres.NewTaskRepository(db)
	bus := event_busses.NewInMemEventBus()
	gateway := gateways.NewDummyPaymentGateway("https://best-payment")

	orderSvc := application.NewOrderService(orderRepo, taskRepo, bus, gateway)

	timeoutWorker := workers.NewTimeoutWorker(orderSvc)
	recoveryWorker := workers.NewRecoveryWorker(taskRepo, orderSvc)

	bus.SubscribeOrderCreated(timeoutWorker.Process)

	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go recoveryWorker.Start(rootCtx)

	mux := http.NewServeMux()
	handler := httpHandlers.NewOrderHandler(orderSvc)
	handler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Printf("server is starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-rootCtx.Done()
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server stopped")
}
