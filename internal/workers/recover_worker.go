package workers

import (
	"context"
	"log"
	"order-system/internal/contracts/repositories"
	"order-system/internal/contracts/services"
	"order-system/internal/domain/tasks"
	"os"
	"time"
)

type RecoveryWorker struct {
	taskRepo     repositories.TaskRepository
	orderService services.OrderService
	logger       *log.Logger
}

func NewRecoveryWorker(tr repositories.TaskRepository, oserv services.OrderService) *RecoveryWorker {
	return &RecoveryWorker{
		taskRepo:     tr,
		orderService: oserv,
		logger:       log.New(os.Stdout, "[RecoveryWorker] ", log.LstdFlags),
	}
}

func (w *RecoveryWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.checkAndRecover(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (w *RecoveryWorker) checkAndRecover(ctx context.Context) {
	tasksToRecover, err := w.taskRepo.GetPendingTasks(ctx)
	if err != nil {
		w.logger.Println("error getting pending tasks: ", err)
		return
	}

	for _, task := range tasksToRecover {
		w.logger.Printf("attempting to recover task: %s", task.Id)
		_ = w.taskRepo.UpdateStatus(ctx, task.Id, tasks.PENDING)

		err := w.orderService.CancelOrder(ctx, task.OrderId)
		if err != nil {
			w.logger.Println("error canceling order: ", task.OrderId, err)
			_ = w.taskRepo.UpdateStatus(ctx, task.Id, tasks.FAILED)
		} else {
			_ = w.taskRepo.UpdateStatus(ctx, task.Id, tasks.COMPLETED)
		}
	}
}
