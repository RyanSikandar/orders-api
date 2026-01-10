package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/RyanSikandar/orders-api/model"
	"github.com/RyanSikandar/orders-api/repository/order"
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
	// Placeholder implementation for getting an order by ID
	fmt.Fprintln(w, "Order details by ID")
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation for updating an order by ID
	fmt.Fprintln(w, "Order updated by ID")
}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation for deleting an order by ID
	fmt.Fprintln(w, "Order deleted by ID")
}
