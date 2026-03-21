package http

type CreateOrderRequest struct {
	Price float64 `json:"price"`
}

type CreateOrderResponse struct {
	PaymentURL string `json:"payment_url"`
	OrderId    string `json:"order_id"`
}
