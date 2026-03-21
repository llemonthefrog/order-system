package workers

import (
	"context"
	"log"
	"order-system/internal/contracts/services"
	"os"
	"time"

	"github.com/google/uuid"
)

type TimeoutWorker struct {
	orderService services.OrderService
	logger       *log.Logger
}

func NewTimeoutWorker(svc services.OrderService) *TimeoutWorker {
	return &TimeoutWorker{
		orderService: svc,
		logger:       log.New(os.Stdout, "[TimeoutWorker] ", log.LstdFlags),
	}
}

func (w *TimeoutWorker) Process(ctx context.Context, orderID uuid.UUID) {
	w.logger.Println("processing order with id", orderID)

	go func() {
		<-time.After(15 * time.Minute)

		w.logger.Printf("timeout reached for %s. Executing cancellation...", orderID)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := w.orderService.CancelOrder(ctx, orderID); err != nil {
			w.logger.Printf("failed to cancel order %s: %v", orderID, err)
		}
	}()
}
