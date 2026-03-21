package orders

import "github.com/google/uuid"

type OrderState string

const (
	PENDING   OrderState = "PENDING"
	COMPLETED OrderState = "COMPLETED"
	PAID      OrderState = "PAID"
	CANCELED  OrderState = "CANCELED"
)

type Order struct {
	Id    uuid.UUID
	State OrderState
	Price float64
}
