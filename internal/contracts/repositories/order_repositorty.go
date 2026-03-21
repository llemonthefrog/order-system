package repositories

import (
	"context"
	"order-system/internal/domain/orders"

	"github.com/google/uuid"
)

type OrderRepository interface {
	Save(ctx context.Context, order *orders.Order) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status orders.OrderState) error
	GetById(ctx context.Context, id uuid.UUID) (*orders.Order, error)
}
