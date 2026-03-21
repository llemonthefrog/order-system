package postgres

import (
	"context"
	"database/sql"
	"errors"
	domainErrors "order-system/internal/domain/errors"
	"order-system/internal/domain/orders"

	"github.com/google/uuid"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (repo *OrderRepository) Save(ctx context.Context, order *orders.Order) error {
	query := `INSERT INTO orders (id, state) VALUES ($1, $2)`

	_, err := repo.db.ExecContext(ctx, query, order.Id, order.State)
	if err != nil {
		return err
	}

	return nil
}

func (repo *OrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status orders.OrderState) error {
	query := `UPDATE orders SET state = $1 WHERE id = $2`

	res, err := repo.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domainErrors.EntityNotFoundError{Id: id}
	}

	return nil
}

func (repo *OrderRepository) GetById(ctx context.Context, id uuid.UUID) (*orders.Order, error) {
	query := `SELECT id, state FROM orders WHERE id = $1`

	var order orders.Order
	err := repo.db.QueryRowContext(ctx, query, id).Scan(&order.Id, &order.State)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainErrors.EntityNotFoundError{Id: id}
		}
		return nil, err
	}

	return &order, nil
}
