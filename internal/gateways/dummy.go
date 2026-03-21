package gateways

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
)

type DummyPaymentGateway struct {
	baseURL string
	logger  *log.Logger
}

func NewDummyPaymentGateway(baseURL string) *DummyPaymentGateway {
	return &DummyPaymentGateway{
		baseURL: baseURL,
		logger:  log.New(os.Stdout, "[PaymentGateway] ", log.LstdFlags),
	}
}

func (g *DummyPaymentGateway) GeneratePayLink(ctx context.Context, orderId uuid.UUID, amount float64) (string, error) {
	payURL := fmt.Sprintf("%s/checkout/%s?amount=%.2f", g.baseURL, orderId.String(), amount)

	g.logger.Printf("generated URL for order %s on sum %.2f: %s", orderId, amount, payURL)

	return payURL, nil
}
