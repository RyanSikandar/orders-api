package handler

import (
	"fmt"
	"net/http"

	"github.com/RyanSikandar/orders-api/repository/order"
)

type Order struct{
	Repo *order.RedisRepo
}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation for creating an order
	fmt.Fprintln(w, "Order created successfully")
}

func (o *Order) List(w http.ResponseWriter, r *http.Request){
	// Placeholder implementation for listing orders
	fmt.Fprintln(w, "List of orders")
}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request){
	// Placeholder implementation for getting an order by ID
	fmt.Fprintln(w, "Order details by ID")
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request){
	// Placeholder implementation for updating an order by ID
	fmt.Fprintln(w, "Order updated by ID")
}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request){
	// Placeholder implementation for deleting an order by ID
	fmt.Fprintln(w, "Order deleted by ID")
}