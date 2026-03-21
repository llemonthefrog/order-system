package in_memory

import (
	"context"
	"order-system/internal/domain/errors"
	"order-system/internal/domain/orders"
	"sync"

	"github.com/google/uuid"
)

type OrderRepository struct {
	mu        sync.RWMutex
	container map[uuid.UUID]*orders.Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		container: make(map[uuid.UUID]*orders.Order),
	}
}

func (repo *OrderRepository) Save(ctx context.Context, order *orders.Order) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.container[order.Id] = order
	return nil
}

func (repo *OrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status orders.OrderState) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	order, ok := repo.container[id]
	if !ok {
		return errors.EntityNotFoundError{Id: id}
	}

	order.State = status
	repo.container[id] = order
	return nil
}

func (repo *OrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*orders.Order, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	order, ok := repo.container[id]
	if !ok {
		return nil, errors.EntityNotFoundError{Id: id}
	}

	return order, nil
}
