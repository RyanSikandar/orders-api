package application

import (
	"fmt"
	"net/http"

	"github.com/RyanSikandar/orders-api/handler"
	"github.com/RyanSikandar/orders-api/repository/order"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	// Add logging middleware to track all incoming requests for debugging and monitoring
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World")
	})

	// Group all order-related routes under /orders prefix for better organization
	// This keeps route definitions modular and easier to maintain
	router.Route("/orders", a.loadOrderRoutes)

	a.router = router
}

func (a *App) loadOrderRoutes(r chi.Router) {
	orderHandler := &handler.Order{
		Repo: &order.RedisRepo{Client: a.rdb},
	}

	// Define RESTful routes following standard HTTP method conventions:
	// GET for retrieval, POST for creation, PUT for updates, DELETE for removal
	r.Get("/", orderHandler.List)
	r.Post("/", orderHandler.Create)
	r.Get("/{id}", orderHandler.GetByID)
	r.Put("/{id}", orderHandler.UpdateByID)
	r.Delete("/{id}", orderHandler.DeleteByID)
}
