package repositories

import (
	"context"
	"order-system/internal/domain/tasks"

	"github.com/google/uuid"
)

type TaskRepository interface {
	Save(ctx context.Context, task *tasks.ScheduledTask) error
	GetPendingTasks(ctx context.Context) ([]*tasks.ScheduledTask, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status tasks.TaskStatus) error
}
