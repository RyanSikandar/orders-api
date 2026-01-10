package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/RyanSikandar/orders-api/model"
	"github.com/RyanSikandar/orders-api/repository/order"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Order struct {
	Repo *order.RedisRepo
}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	// Defining Body of the Order as a struct
	var newOrder struct {
		CustomerID uuid.UUID        `json:"customer_id"`
		LineItems  []model.LineItem `json:"line_items"`
	}

	// Decode the JSON body into the newOrder struct
	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var now = time.Now().UTC()

	// Create a new Order instance
	order := model.Order{
		ID:         rand.Int(),
		CustomerID: newOrder.CustomerID,
		LineItems:  newOrder.LineItems,
		CreatedAt:  &now,
	}

	// Save the new order using the repository
	if err := o.Repo.Create(r.Context(), order); err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	responseBody := struct {
		OrderID int `json:"order_id"`
	}{
		OrderID: order.ID,
	}

	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")

	if cursorStr == "" {
		cursorStr = "0"
	}

	// Convert cursor to int64
	var cursor uint64
	_, err := fmt.Sscan(cursorStr, &cursor) // This parses the string to int64 or we can use the strconv package as well
	if err != nil {
		http.Error(w, "Invalid cursor value", http.StatusBadRequest)
		return
	}

	const size = 50

	res, err := o.Repo.List(r.Context(), order.OrderIterator{
		Size:   size,
		Offset: cursor,
	})

	if err != nil {
		http.Error(w, "Failed to list orders", http.StatusInternalServerError)
		return
	}

	var responseBody struct {
		Orders []model.Order `json:"orders"`
		Next   int           `json:"next"`
	}

	responseBody.Orders = res.Orders
	responseBody.Next = res.Cursor

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	// use chi to get the path params

	idParam := chi.URLParam(r, "id")

	// Convert idParam to int
	var id int
	_, err := fmt.Sscan(idParam, &id)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := o.Repo.GetByID(r.Context(), id)

	if err != nil {
		http.Error(w, "Failed to get order by ID", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	var updatedOrder struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updatedOrder); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Get the id from the URL params
	idParam := chi.URLParam(r, "id")

	// Convert idParam to int
	var id int
	_, err := fmt.Sscan(idParam, &id)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := o.Repo.GetByID(r.Context(), id)

	if err != nil {
		http.Error(w, "Failed to get order by ID", http.StatusInternalServerError)
		return
	}

	const completedStatus = "completed"
	const shippedStatus = "shipped"

	now := time.Now().UTC()

	switch updatedOrder.Status {
	case completedStatus:
		if order.ShippedAt != nil && order.CompletedAt == nil {
			order.CompletedAt = &now
			break
		}

		http.Error(w, "Order must be shipped before it can be completed", http.StatusBadRequest)
		return

	case shippedStatus:
		if order.ShippedAt != nil {
			http.Error(w, "Order is already shipped", http.StatusBadRequest)
			return
		}
		order.ShippedAt = &now
	default:
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	if err := o.Repo.UpdateByID(r.Context(), order); err != nil {
		http.Error(w, "Failed to update order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	// Convert idParam to int
	var id int
	_, err := fmt.Sscan(idParam, &id)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	if err := o.Repo.DeleteByID(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
