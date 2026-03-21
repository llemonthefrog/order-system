package contracts

import (
	"context"

	"github.com/google/uuid"
)

type PaymentGateway interface {
	GeneratePayLink(ctx context.Context, orderId uuid.UUID, amount float64) (string, error)
}
