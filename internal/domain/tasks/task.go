package tasks

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	PENDING   TaskStatus = "PENDING"
	FAILED    TaskStatus = "FAILED"
	COMPLETED TaskStatus = "COMPLETED"
)

type ScheduledTask struct {
	Id        uuid.UUID
	OrderId   uuid.UUID
	ExecuteAt time.Time
	Status    TaskStatus
}
