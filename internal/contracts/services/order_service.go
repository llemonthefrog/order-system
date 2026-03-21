package services

import (
	"context"
	"order-system/internal/domain/orders"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *orders.Order) (string, error)
	ConfirmOrder(ctx context.Context, orderId uuid.UUID) error
	CancelOrder(ctx context.Context, orderId uuid.UUID) error
}
