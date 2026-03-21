package postgres

import (
	"context"
	"database/sql"

	"order-system/internal/domain/tasks"

	"github.com/google/uuid"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Save(ctx context.Context, task *tasks.ScheduledTask) error {
	query := `INSERT INTO scheduled_tasks (id, order_id, execute_at) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, task.Id, task.OrderId, task.ExecuteAt)
	return err
}

func (r *TaskRepository) GetPendingTasks(ctx context.Context) ([]*tasks.ScheduledTask, error) {
	query := `SELECT id, order_id, execute_at FROM scheduled_tasks 
              WHERE status = 'pending' AND execute_at <= NOW()`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*tasks.ScheduledTask
	for rows.Next() {
		t := &tasks.ScheduledTask{}
		if err := rows.Scan(&t.Id, &t.OrderId, &t.ExecuteAt); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (r *TaskRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status tasks.TaskStatus) error {
	query := `UPDATE scheduled_tasks SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}
