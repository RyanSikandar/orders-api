package application

import (
	"fmt"
	"net/http"

	"github.com/RyanSikandar/orders-api/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()

	// Add logging middleware to track all incoming requests for debugging and monitoring
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World")
	})

	// Group all order-related routes under /orders prefix for better organization
	// This keeps route definitions modular and easier to maintain
	router.Route("/orders", loadOrderRoutes)

	return router
}

func loadOrderRoutes(r chi.Router) {
	orderHandler := &handler.Order{}

	// Define RESTful routes following standard HTTP method conventions:
	// GET for retrieval, POST for creation, PUT for updates, DELETE for removal
	r.Get("/", orderHandler.List)
	r.Post("/", orderHandler.Create)
	r.Get("/{id}", orderHandler.GetByID)
	r.Put("/{id}", orderHandler.UpdateByID)
	r.Delete("/{id}", orderHandler.DeleteByID)
}
