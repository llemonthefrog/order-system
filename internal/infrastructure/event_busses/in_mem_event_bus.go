package event_busses

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type InMemEventBus struct {
	mu          sync.RWMutex
	subscribers []func(context.Context, uuid.UUID)
}

func NewInMemEventBus() *InMemEventBus {
	return &InMemEventBus{
		subscribers: make([]func(context.Context, uuid.UUID), 0),
	}
}

func (b *InMemEventBus) SubscribeOrderCreated(handler func(context.Context, uuid.UUID)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers = append(b.subscribers, handler)
}

func (b *InMemEventBus) PublishOrderCreated(ctx context.Context, orderID uuid.UUID) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	for _, handler := range b.subscribers {
		go handler(ctx, orderID)
	}
	return nil
}
