package http

import (
	"encoding/json"
	"net/http"
	"order-system/internal/contracts/services"
	"order-system/internal/domain/orders"

	"github.com/google/uuid"
)

type OrderHandler struct {
	service services.OrderService
}

func NewOrderHandler(svc services.OrderService) *OrderHandler {
	return &OrderHandler{service: svc}
}
func (h *OrderHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /orders/create", h.HandleCreateOrder)
	mux.HandleFunc("POST /orders/{id}/confirm", h.HandleConfirmOrder)
}

func (h *OrderHandler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	newOrder := &orders.Order{
		Id:    uuid.New(),
		Price: req.Price,
		State: orders.PENDING,
	}

	payURL, err := h.service.CreateOrder(r.Context(), newOrder)
	if err != nil {
		writeError(w, "failed to create order", http.StatusInternalServerError)
		return
	}

	res := CreateOrderResponse{
		PaymentURL: payURL,
		OrderId:    newOrder.Id.String(),
	}

	writeJSON(w, http.StatusCreated, res)
}

func (h *OrderHandler) HandleConfirmOrder(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")

	orderID, err := uuid.Parse(idParam)
	if err != nil {
		writeError(w, "invalid order Id format", http.StatusBadRequest)
		return
	}

	if err := h.service.ConfirmOrder(r.Context(), orderID); err != nil {
		writeError(w, "failed to confirm order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
