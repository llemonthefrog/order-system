package contracts

import (
	"context"

	"github.com/google/uuid"
)

type EventBus interface {
	PublishOrderCreated(ctx context.Context, orderID uuid.UUID) error
	SubscribeOrderCreated(handler func(ctx context.Context, orderID uuid.UUID))
}
