package application

import (
	"context"
	"log"
	"order-system/internal/contracts"
	"order-system/internal/contracts/repositories"
	"order-system/internal/domain/orders"
	"order-system/internal/domain/tasks"
	"os"
	"time"

	"github.com/google/uuid"
)

type OrderServiceImpl struct {
	repo     repositories.OrderRepository
	taskRepo repositories.TaskRepository
	bus      contracts.EventBus
	gateway  contracts.PaymentGateway
	logger   *log.Logger
}

func NewOrderService(
	repo repositories.OrderRepository,
	taskRepo repositories.TaskRepository,
	bus contracts.EventBus,
	gateway contracts.PaymentGateway,
) *OrderServiceImpl {
	return &OrderServiceImpl{
		repo:     repo,
		taskRepo: taskRepo,
		bus:      bus,
		gateway:  gateway,
		logger:   log.New(os.Stdout, "[OrderService] ", log.LstdFlags),
	}
}

func (s *OrderServiceImpl) CreateOrder(ctx context.Context, order *orders.Order) (string, error) {
	s.logger.Printf("Creating order: %s, price: %.2f", order.Id, order.Price)

	if err := s.repo.Save(ctx, order); err != nil {
		s.logger.Printf("failed to save order %s: %v", order.Id, err)
		return "", err
	}

	task := &tasks.ScheduledTask{
		Id:        uuid.New(),
		OrderId:   order.Id,
		ExecuteAt: time.Now().Add(15 * time.Minute),
		Status:    tasks.PENDING,
	}

	if err := s.taskRepo.Save(ctx, task); err != nil {
		s.logger.Printf("failed to save task for order %s: %v", order.Id, err)
		return "", err
	}

	url, err := s.gateway.GeneratePayLink(ctx, order.Id, order.Price)
	if err != nil {
		s.logger.Printf("failed to generate payment link for order %s: %v", order.Id, err)
		return "", err
	}

	_ = s.bus.PublishOrderCreated(ctx, order.Id)

	return url, nil
}

func (s *OrderServiceImpl) ConfirmOrder(ctx context.Context, orderId uuid.UUID) error {
	s.logger.Printf("confirming order: %s", orderId)

	err := s.repo.UpdateStatus(ctx, orderId, orders.PAID)
	if err != nil {
		s.logger.Printf("failed to confirm order %s: %v", orderId, err)
	}
	return err
}

func (s *OrderServiceImpl) CancelOrder(ctx context.Context, orderID uuid.UUID) error {
	s.logger.Printf("attempting to cancel order: %s", orderID)

	order, err := s.repo.GetById(ctx, orderID)
	if err != nil {
		s.logger.Printf("failed to get order %s for cancellation: %v", orderID, err)
		return err
	}

	if order.State == orders.PENDING {
		if err := s.repo.UpdateStatus(ctx, orderID, orders.CANCELED); err != nil {
			s.logger.Printf("failed to cancel order %s: %v", orderID, err)
			return err
		}
		s.logger.Printf("order %s has been CANCELED due to timeout", orderID)
	} else {
		s.logger.Printf("order %s cancellation skipped (current state: %s)", orderID, order.State)
	}

	return nil
}
